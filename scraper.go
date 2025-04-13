package main

/*
Function name
desc.

newline continued desc.

# Parameters:
  - param1 type

desc.

  - param2 type

...

# Returns:
  - type

desc.
*/

import (
	"fmt"
	"log"
	"regexp"
	"sync"

	"github.com/gocolly/colly/v2"
)

const _DEBUG_ = false

const _QUEUEMAXLEN_ = 5
const _LINKQUEUELEN_ = 200
const _IDKEY_ = `QueueIndex`

var shortendLinkRegex *regexp.Regexp
var nonAllowedRegex *regexp.Regexp

type WorkspaceID int
type URLString string

type context struct {

	// buffered queue channel of links scraped from previously visited sites
	linkQueue chan URLString
	// buffered queue channel of links inputed to the scraper
	priorityLinkQueue chan URLString
	// worspace for the collectors async requests
	// each request gets their own index and append the text they recieve to the slice
	workspaceBuffer [_QUEUEMAXLEN_][]string
	// channel of indexes in the 'workspace' array ready to be assigned to a request
	// (buffer of size 'QUEUEMAXLEN')
	emptyIndexes chan WorkspaceID
	//  channel of indexes in the 'workspace' array ready to be read
	// (buffer of size 'QUEUEMAXLEN')
	finishedIndexes chan WorkspaceID
}

type counter struct {
	// mutex used to sync changes to the two counters in context
	counterLock sync.Mutex
	// currently working on amount
	workingCounter int
	// currently finished amount should be in sync with len(finishedIndexes)
	finishedCounter int
}

type CollectorStruct struct {
	// struct holding all cotext to make the inteface with the collector as simple as possible
	context context
	//
	counter counter

	// the colly collector used for webb scraping and formatting
	collectorColly *colly.Collector
}

func debugPrint(a ...any) {
	if _DEBUG_ {
		fmt.Println(a...)
	}
}

func main() {
	collector := CollectorSetup()

	RequestVisitToSite(collector, "https://en.wikipedia.org/wiki/Cucumber")
	go CollectorRepopulate(collector)

	ReadAndPrint(collector)

	ReadAndPrint(collector)

	ReadAndPrint(collector)

}

func counterSync(collector *CollectorStruct, f func(counter *counter)) {
	counter := &collector.counter
	counter.counterLock.Lock()
	f(counter)
	counter.counterLock.Unlock()

}

/*
CollectorRepopulate
is used to request the scraper to scrape enough websites to fill the buffer.

It will block until it has enough links in the queue for all its requests.
Is safe to run in a seperate go rutine.

# Parameters:
  - collector *CollectorStruct

The struct containin the scraper and all used context.
*/
func CollectorRepopulate(collector *CollectorStruct) {
	counterSync(collector, func(counter *counter) {

		amountFilled := counter.finishedCounter + counter.workingCounter
		amountEmpty := _QUEUEMAXLEN_ - amountFilled
		for range amountEmpty {
			visitNextLink(collector, counter)
		}
	})
}

/*
CollectorRepopulateFixedNumber
is used to request the scraper to scrape a specified amount of websites.

The scraper is using a fixed sized buffer
which means that it isnt always possible to fit the amount of requests made.
Therefore it will return the amount of requests that didnt get prossesed.

It will block until it has enough links in the queue for all its requests.
Is safe to run in a seperate go rutine.

# Parameters:
  - collector *CollectorStruct

The struct containing the scraper and all used context.

  - n int

The amount of sites to scrape.

# Returns:
  - int

The amount of requests that couldn't be fullfilled.
*/
func CollectorRepopulateFixedNumber(collector *CollectorStruct, n int) int {
	amountDidntFit := 0
	counterSync(collector, func(counter *counter) {

		amountFilled := counter.finishedCounter + counter.workingCounter
		amountEmpty := _QUEUEMAXLEN_ - amountFilled
		if amountEmpty < n {
			amountDidntFit = n - amountEmpty
			n = amountEmpty
		}
		for range n {
			visitNextLink(collector, counter)
		}
	})

	return amountDidntFit
}

/*
ReadAndPrint
reads the first avaliable fully scraped site and prints the content.

# Parameters:
  - collector *CollectorStruct

The struct containing the scraper and all used context.
*/
func ReadAndPrint(collector *CollectorStruct) {
	stringSlice := ReadFinished(collector)
	for _, text := range stringSlice {
		fmt.Println(text)
	}
}

/*
claimNewIndex
claims and initialises a space in the worspace buffer.

# Parameters:
  - collector *CollectorStruct

The struct containing the scraper and all used context.

  - url string

The url that the worker is speaking with and is used to initialise the slice.

# Returns:
  - int

The index of the claimed space in the buffer
*/
func claimNewIndex(context *context, url string) WorkspaceID {
	ID := <-context.emptyIndexes
	context.workspaceBuffer[ID] = []string{url}
	return ID
}

/*
ReadFinished
retrieves a fully scraped page and returns it.
If there is no page to retrieve it will block until one gets avaliable.

# Parameters:
  - collector *CollectorStruct

The struct containing the scraper and all used context.

# Returns:
  - []string

The slice containing the text from the scraped page.
*/
func ReadFinished(collector *CollectorStruct) []string {
	context := &collector.context
	// removes 1 from finished
	counterSync(collector, func(counter *counter) {
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

# Parameters:
  - collector *CollectorStruct

The struct containing the scraper and all used context.

  - index int

The index or ID of the workspace to use.

  - text string

The text to appended.
*/
func writeToWorkspace(context *context, ID WorkspaceID, text string) {
	path := &context.workspaceBuffer[ID]
	*path = append(*path, text)
}

/*
Function name
desc.

newline continued desc.

# Parameters:
  - param1 type

desc.

  - param2 type

...

# Returns:
  - type

desc.
*/
func CollectorSetup() *CollectorStruct {
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

	shortendLinkRegex, _ = regexp.Compile(`^/wiki/`)
	nonAllowedRegex, _ = regexp.Compile(`/wiki/(File|Wikipedia|Special|User):`)

	c := colly.NewCollector(colly.AllowedDomains("en.wikipedia.org"))
	collector.collectorColly = c

	c.Async = true

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
		ID := claimNewIndex(context, url)
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
		writeToWorkspace(context, ID, e.Text)

	})

	// Find and visit all links
	c.OnHTML(`a[href*="/wiki/"]`, func(e *colly.HTMLElement) {
		addScrapedLinkToQueue(e, context)
	})

	// triggered once scraping is done
	c.OnScraped(func(r *colly.Response) {
		mapValue := r.Ctx.GetAny(_IDKEY_)
		ID, ok := mapValue.(WorkspaceID)
		if !ok {
			log.Fatal("couldn't find ID")
		}
		counterSync(collector, func(counter *counter) {
			counter.workingCounter--
			counter.finishedCounter++
		})
		context.finishedIndexes <- ID

	})
	return collector
}

/*
Function name
desc.

newline continued desc.

# Parameters:
  - param1 type

desc.

  - param2 type

...

# Returns:
  - type

desc.
*/
func addScrapedLinkToQueue(e *colly.HTMLElement, context *context) {

	// TODO: filter away already visited before adding to the queue

	link := e.Attr("href")
	if !shortendLinkRegex.MatchString(link) || nonAllowedRegex.MatchString(link) {
		debugPrint("non allowed link: ", link)
		return
	}
	link = "https://en.wikipedia.org" + link
	select {
	case context.linkQueue <- URLString(link):
	default:
		debugPrint("linkQueue full skipped link: ", link)
	}

}

/*
Function name
desc.

newline continued desc.

# Parameters:
  - param1 type

desc.

  - param2 type

...

# Returns:
  - type

desc.
*/
func RequestVisitToSite(collector *CollectorStruct, link string) {
	if len(collector.context.linkQueue) == 0 {
		collector.context.linkQueue <- URLString(link)
		return
	}

	collector.context.priorityLinkQueue <- URLString(link)
}

/*
Function name
desc.

should only be called while having the counter lock

# Parameters:
  - param1 type

desc.

  - param2 type

...

# Returns:
  - type

desc.
*/
func visitNextLink(collector *CollectorStruct, counter *counter) {
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
