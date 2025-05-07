package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
)

func debugPrint(a ...any) {
	if _DEBUG_ {
		log.Println(a...)
	}
}

/*
counterSync
syncs changes to the counter using a mutex.

# Parameters:

  - f func(counter *counter)

The function to run while owning the mutex
*/
func (collector *CollectorStruct) counterSync(f func(counter *counter)) {
	counter := &collector.counter

	counter.counterLock.Lock()
	f(counter)
	counter.counterLock.Unlock()

}

/*
CollectorRepopulate
requests the scraper to scrape enough websites to fill the buffer.

It will block until it has enough links in the queue for all its requests.
Is safe to run in a seperate go rutine.
*/
func (collector *CollectorStruct) CollectorRepopulate() {
	collector.counterSync(func(counter *counter) {

		amountFilled := counter.finishedCounter + counter.workingCounter
		amountEmpty := _WORKSPACES_ - amountFilled
		for range amountEmpty {
			collector.visitNextLink(counter)
		}
	})
}

/*
CollectorRepopulateFixedNumber
requests the scraper to scrape a specified amount of websites.

The scraper is using a fixed sized buffer
which means that it isnt always possible to fit the amount of requests made.
Therefore it will return the amount of requests that didnt get prossesed.

It will block until it has enough links in the queue for all its requests.
Is safe to run in a seperate go rutine.

# Returns:

The amount of requests that couldn't be fullfilled.
*/
func (collector *CollectorStruct) CollectorRepopulateFixedNumber(
	amountToScrape int) int {
	amountDidntFit := 0
	collector.counterSync(func(counter *counter) {

		amountFilled := counter.finishedCounter + counter.workingCounter
		amountEmpty := _WORKSPACES_ - amountFilled
		if amountEmpty < amountToScrape {
			amountDidntFit = amountToScrape - amountEmpty
			amountToScrape = amountEmpty
		}
		for range amountToScrape {
			collector.visitNextLink(counter)
		}
	})

	return amountDidntFit
}

/*
ReadAndPrint
reads the first avaliable fully scraped site and prints the content.
*/
func (collector *CollectorStruct) ReadAndPrint() {
	readStrings := collector.ReadFinished()

	var fullString string

	fullString = strings.Join(readStrings, " ")

	iterator := strings.SplitSeq(fullString, "	")

	fullString = ""
	for x := range iterator {
		if x != "" {
			fullString += x + " "
		}
	}
	iterator = strings.SplitSeq(fullString, "\n")
	fullString = ""
	for x := range iterator {
		if x != "" {
			fullString += x + " "

		}
	}
	iterator = strings.SplitSeq(fullString, " ")
	fullString = ""
	for x := range iterator {
		if x != "" {
			fullString += x + " "
		}
	}
	fmt.Println("--------------------------------------------------")
	fmt.Print("\n\n\n")
	fmt.Println(fullString)
	fmt.Print("\n\n\n")
	fmt.Println("--------------------------------------------------")
}

/*
claimNewIndex
claims and initialises a space in the worspace buffer.

# Returns:

The ID of the claimed workspace.
*/
func (context *context) claimNewIndex(bufferInit URLString) WorkspaceID {
	ID := <-context.emptyIndexes
	context.workspaceBuffer[ID] = []string{string(bufferInit)}
	return ID
}

/*
ReadFinished
retrieves a fully scraped page and returns it.
If there is no page to retrieve it will block until one gets avaliable.

# Returns:

A slice containing the text from the scraped page.
*/
func (collector *CollectorStruct) ReadFinished() []string {
	context := &collector.context
	// removes 1 from finished
	collector.counterSync(func(counter *counter) {
		// can currently become negative by this which
		// isn't a case deeply explored but should work fine
		// TODO: test for cases where FinishedCounter becomes negative
		counter.finishedCounter--
	})

	// waits for a Workspace to finish
	ID := <-context.finishedIndexes

	// retrieves the content and empties out the workspace
	pos := &context.workspaceBuffer[ID]
	pageText := *pos
	*pos = nil

	// adds the index to the list of unused/empty workspaces
	context.emptyIndexes <- ID

	return pageText
}

/*
writeToWorkspace
appends text to the specified workspace.
*/
func (context *context) writeToWorkspace(ID WorkspaceID, text string) {
	path := &context.workspaceBuffer[ID]
	*path = append(*path, text)

}

/*
NewCollector
sets up a new collector.

# Parameters:

boolean for turning on and off asyncronous work

# Returns:

A new collector ready to be used.
*/
func NewCollector(async bool, localFiles bool) *CollectorStruct {
	lH := linkHandlerCreate()
	collector := &CollectorStruct{
		context: context{
			linkHandler:     lH,
			workspaceBuffer: [_WORKSPACES_][]string{},
			emptyIndexes:    make(chan WorkspaceID, _WORKSPACES_),
			finishedIndexes: make(chan WorkspaceID, _WORKSPACES_),
		},
		counter: counter{
			workingCounter:  0,
			finishedCounter: 0,
		},
		collectorColly: colly.NewCollector(),
		settings: settings{
			async: async,
		},
	}
	if localFiles {
		collector.settings.HtmlFileType = HtmlFileType("file://")
	} else {
		collector.settings.HtmlFileType = HtmlFileType("https://")
	}
	context := &collector.context
	for x := range _WORKSPACES_ {
		context.emptyIndexes <- WorkspaceID(x)
	}

	c := collector.collectorColly

	c.Async = async
	if localFiles {
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

		c.WithTransport(t)
	}

	// called before an HTTP request is triggered
	c.OnRequest(func(r *colly.Request) {
		debugPrint("Visiting: ", r.URL)
	})

	// triggered when the scraper encounters an error
	c.OnError(func(r *colly.Response, err error) {
		debugPrint(
			"Something went wrong: ", err,
			"\nWhen trying to scrape:", r.Request.URL,
		)
	})

	// fired when the server responds
	c.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.EscapedPath()
		debugPrint("Page visited: ", r.Request.URL)
		host := r.Request.URL.Host
		if shortendLinkRegex.MatchString(url) {
			url = r.Request.URL.Scheme + "://" + host + url
		}
		ID := context.claimNewIndex(URLString(url))
		r.Ctx.Put(_IDKEY_, ID)
	})

	// triggered when a CSS selector matches an element
	c.OnHTML("p, title, h1, h2, h3", func(e *colly.HTMLElement) {
		mapValue := e.Response.Ctx.GetAny(_IDKEY_)
		ID, ok := mapValue.(WorkspaceID)
		if !ok {
			log.Fatal("couldn't find ID")
		}
		context.writeToWorkspace(ID, e.Text)

	})

	// Find and visit all links
	c.OnHTML(`[href]`, func(e *colly.HTMLElement) {
		lH.linkFixer(e, localFiles)

	})

	// triggered once scraping is done
	c.OnScraped(func(r *colly.Response) {
		mapValue := r.Ctx.GetAny(_IDKEY_)
		ID, ok := mapValue.(WorkspaceID)
		if !ok {
			log.Fatal("couldn't find ID")
		}
		debugPrint("Scraped: ", r.Request.URL)

		f := func(counter *counter) {
			counter.workingCounter--
			counter.finishedCounter++
		}

		// needed bypass for synced scraping to work
		if async {
			collector.counterSync(f)
		} else {
			f(&collector.counter)
		}
		context.finishedIndexes <- ID

	})
	return collector
}

/*
linkHandler
is a helperfunction used for adding a scraped link to the queue.

# Parameters:

the matched HTMLElemnt.

# Returns:

bool whether the the link is valid
*/
func (lH *linkHandler) linkFixer(e *colly.HTMLElement, localFile bool) {
	link := e.Attr("href")
	shortenedLink := shortendLinkRegex.MatchString(link)
	fullLocalLink := fullLocalLinkRegex.MatchString(link)
	fullWebLink := fullWebLinkRegex.MatchString(link)

	var host hostPath
	switch {
	case localFile && shortenedLink:
		arr := strings.SplitAfter(e.Request.URL.EscapedPath(), `/`)
		arr[len(arr)-1] = link
		link = strings.Join(arr, "")
		lH.inputLocalFile(URLString(link), false)
	case localFile && fullLocalLink:
		link = strings.Replace(link, `file://`, ``, 1)
		lH.inputLocalFile(URLString(link), false)
	case !localFile && shortenedLink:
		host = hostPath(e.Request.URL.Host)
		lH.inputWeb(host, innerPath(link), false)
	case !localFile && fullWebLink:
		link = strings.Replace(link, `https://`, ``, 1)
		arr := strings.Split(link, `/`)
		host = hostPath(arr[0])
		link = strings.Join(arr[1:], `/`)
		lH.inputWeb(host, innerPath(link), false)
	default:
	}
}

/*
RequestVisitToSite
adds the requested site to the queue of links to visit.
*/
func (collector *CollectorStruct) RequestVisitToSite(link string) {
	fullLocalLink := fullLocalLinkRegex.MatchString(link)
	fullWebLink := fullWebLinkRegex.MatchString(link)
	var host hostPath
	lH := collector.context.linkHandler
	switch {

	case fullLocalLink:
		link = strings.Replace(link, `file://`, ``, 1)
		lH.inputLocalFile(URLString(link), true)

	case fullWebLink:
		link = strings.Replace(link, `https://`, ``, 1)
		arr := strings.Split(link, `/`)
		host = hostPath(arr[0])
		link = `/` + strings.Join(arr[1:], `/`)
		lH.inputWeb(host, innerPath(link), true)
		debugPrint("host:", host, "link:", link, "")
	default:
	}
}

/*
visitNextLink
dispatches a new worker to scrape the next link in the queue.

Should only be called in the scope of [counterSync].
*/
func (collector *CollectorStruct) visitNextLink(counter *counter) {
	for {
		link := collector.context.linkHandler.getLink()
		err := collector.collectorColly.Visit(string(link))

		if err == nil {
			counter.workingCounter++
			break
		}
		switch err.Error() {
		case `Forbidden domain`:
			debugPrint(err, link)
		case `URL already visited`:
			debugPrint(err, link)
		default:
			debugPrint(err)
		}
	}
}

func (collector *CollectorStruct) Whitelist(host string) {
	collector.context.linkHandler.filter.Whitelist(hostPath(host), true)
}
