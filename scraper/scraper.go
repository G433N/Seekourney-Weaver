package scraper

import (
	"fmt"
	"log"
	"regexp"
	"sync"

	"github.com/gocolly/colly/v2"
)

const (

	// Used for enabling debug prints
	_DEBUG_ = false
	// Total amount of worspaces in the collector
	_QUEUEMAXLEN_ = 5
	// Buffer size of the link and priority link channels
	_LINKQUEUELEN_ = 20
	// Key used to access workspace ID in a map
	_IDKEY_ = `WorkspaceID`
)

var (
	//Matches when the link is in the same main domain as the current webbsite
	shortendLinkRegex = regexp.MustCompile(`^/[a-zA-Z]`)
	//Parts of wikipedia not worth indexing
	WikipediaBadRegex = regexp.MustCompile(
		`/wiki/(File|Wikipedia|Special|User)|/static/|/w/`)
)

type (
	WorkspaceID int
	URLString   string

	context struct {

		// buffered queue channel of links scraped from previously visited sites
		linkQueue chan URLString

		// buffered queue channel of links inputed to the scraper
		priorityLinkQueue chan URLString

		// worspace for the collectors async requests
		// each request gets their own index
		// and then appends the text they recieve to the slice
		workspaceBuffer [_QUEUEMAXLEN_][]string

		// channel of indexes in the `workspace` array
		// ready to be assigned to a request
		// (buffer of size 'QUEUEMAXLEN')
		emptyIndexes chan WorkspaceID

		//  channel of indexes in the `workspace` array ready to be read
		// (buffer of size 'QUEUEMAXLEN')
		finishedIndexes chan WorkspaceID
	}

	counter struct {

		// mutex used to sync changes to the two counters in context
		counterLock sync.Mutex

		// currently working on amount
		workingCounter int

		// currently finished amount should be in sync with len(finishedIndexes)
		finishedCounter int
	}

	CollectorStruct struct {

		// struct holding all context
		// to make the inteface with the collector as simple as possible
		context context

		// used to keep track of number of workspaces in use
		counter counter

		// the colly collector used for webb scraping and formatting
		collectorColly *colly.Collector
	}
)

func debugPrint(a ...any) {
	if _DEBUG_ {
		fmt.Println(a...)
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
		amountEmpty := _QUEUEMAXLEN_ - amountFilled
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
		amountEmpty := _QUEUEMAXLEN_ - amountFilled
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
	stringSlice := collector.ReadFinished()
	fmt.Print("\n\n\n")
	for _, text := range stringSlice {
		fmt.Println(text)
	}
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
func NewCollector(async bool) *CollectorStruct {
	collector := &CollectorStruct{}
	collector.context = context{}
	context := &collector.context
	collector.counter.finishedCounter = 0
	collector.counter.workingCounter = 0
	context.emptyIndexes = make(chan WorkspaceID, _QUEUEMAXLEN_)
	for x := range _QUEUEMAXLEN_ {
		context.emptyIndexes <- WorkspaceID(x)
	}
	context.finishedIndexes = make(chan WorkspaceID, _QUEUEMAXLEN_)
	context.linkQueue = make(chan URLString, _LINKQUEUELEN_)
	context.priorityLinkQueue = make(chan URLString, _LINKQUEUELEN_)
	context.workspaceBuffer = [_QUEUEMAXLEN_][]string{}

	c := colly.NewCollector()
	collector.collectorColly = c

	c.Async = async

	// called before an HTTP request is triggered
	c.OnRequest(func(r *colly.Request) {
		debugPrint("Visiting: ", r.URL)
	})

	// triggered when the scraper encounters an error
	c.OnError(func(_ *colly.Response, err error) {
		debugPrint("Something went wrong: ", err)
	})

	// fired when the server responds
	c.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.EscapedPath()
		debugPrint("Page visited: ", r.Request.URL)
		if shortendLinkRegex.MatchString(url) {
			url = "https://en.wikipedia.org" + url
		}
		ID := context.claimNewIndex(URLString(url))
		r.Ctx.Put(_IDKEY_, ID)
	})

	// triggered when a CSS selector matches an element
	c.OnHTML("p, div.mw-heading", func(e *colly.HTMLElement) {
		// printing all URLs associated with the <p> tag on the page
		mapValue := e.Response.Ctx.GetAny(_IDKEY_)
		ID, ok := mapValue.(WorkspaceID)
		if !ok {
			log.Fatal("couldn't find ID")
		}
		context.writeToWorkspace(ID, e.Text)

	})

	// Find and visit all links
	c.OnHTML(`[href]`, func(e *colly.HTMLElement) {
		link, valid := linkHandler(e)
		if !valid {
			return
		}
		select {
		case context.linkQueue <- URLString(link):
		default:
			debugPrint("linkQueue full skipped link: ", link)
		}
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
func linkHandler(e *colly.HTMLElement) (string, bool) {

	// TODO: filter away already visited

	link := e.Attr("href")
	host := e.Request.URL.Host

	switch host {
	case `en.wikipedia.org`:
		if WikipediaBadRegex.MatchString(link) {
			debugPrint("Not worth indexing: ", link)
			return "", false
		}
	default:
	}
	if !shortendLinkRegex.MatchString(link) {
		debugPrint("Not allowed for crawler to leave host\nhost:",
			host, "\nlink:", link)
		return "", false
	}
	return "https://" + host + link, true
}

/*
RequestVisitToSite
adds the requested site to the queue of links to visit.
*/
func (collector *CollectorStruct) RequestVisitToSite(webbpageLink string) {
	if len(collector.context.linkQueue) == 0 {
		collector.context.linkQueue <- URLString(webbpageLink)
		return
	}

	collector.context.priorityLinkQueue <- URLString(webbpageLink)
}

/*
visitNextLink
dispatches a new worker to scrape the next link in the queue.

Should only be called in the scope of [counterSync].
*/
func (collector *CollectorStruct) visitNextLink(counter *counter) {
	for {
		var link URLString
		select {
		case link = <-collector.context.priorityLinkQueue:
		default:
			link = <-collector.context.linkQueue
		}
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
