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

const DEBUG = false

const QUEUEMAXLEN = 5
const LINKQUEUELEN = 200
const INDEXKEY = `QueueIndex`

var ShortendLinkRegex *regexp.Regexp
var NonAllowedRegex *regexp.Regexp

// initialize a data structure to keep the scraped data
type Context struct {
	// currently working on amount
	WorkingCounter int
	// currently finished amount should be in sync with len(finishedIndexes)
	FinishedCounter int
	// buffered queue channel of links scraped from previously visited sites
	LinkQueue chan string
	// buffered queue channel of links inputed to the scraper
	PriorityLinkQueue chan string
	// worspace for the collectors async requests
	// each request gets their own index and append the text they recieve to the slice
	WorkspaceBuffer [QUEUEMAXLEN][]string
	// channel of indexes in the 'workspace' array ready to be assigned to a request
	// (buffer of size 'QUEUEMAXLEN')
	EmptyIndexes chan int
	//  channel of indexes in the 'workspace' array ready to be read
	// (buffer of size 'QUEUEMAXLEN')
	FinishedIndexes chan int
}

type CollectorStruct struct {
	// mutex used to sync changes to the two counters in context
	CounterLock sync.Mutex

	// struct holding all cotext to make the inteface with the collector as simple as possible
	Context Context

	// the colly collector used for webb scraping and formatting
	CollectorColly *colly.Collector
}

func debugPrint(a ...any) {
	if DEBUG {
		fmt.Println(a...)
	}
}

func main() {
	collector := collectorSetup()

	RequestVisitToSite(collector, "https://en.wikipedia.org/wiki/Cucumber")
	go collectorRepopulate(collector)

	readAndPrint(collector)

	readAndPrint(collector)

	readAndPrint(collector)

}

/*
collectorRepopulate
is used to request the scraper to scrape enough websites to fill the buffer.

It will block until it has enough links in the queue for all its requests.
Is safe to run in a seperate go rutine.

# Parameters:
  - collector *CollectorStruct

The struct containin the scraper and all used context.
*/
func collectorRepopulate(collector *CollectorStruct) {
	context := &collector.Context
	countLock := &collector.CounterLock
	countLock.Lock()
	amountFilled := context.FinishedCounter + context.WorkingCounter
	amountEmpty := QUEUEMAXLEN - amountFilled
	for range amountEmpty {
		VisitNextLink(collector)
	}
	countLock.Unlock()

}

/*
collectorRepopulateFixedNumber
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
func collectorRepopulateFixedNumber(collector *CollectorStruct, n int) int {
	amountDidntFit := 0
	context := &collector.Context
	countLock := &collector.CounterLock
	countLock.Lock()
	amountFilled := context.FinishedCounter + context.WorkingCounter
	amountEmpty := QUEUEMAXLEN - amountFilled
	if amountEmpty < n {
		amountDidntFit = n - amountEmpty
		n = amountEmpty
	}
	for range n {
		VisitNextLink(collector)
	}
	countLock.Unlock()
	return amountDidntFit
}

/*
readAndPrint
reads the first avaliable fully scraped site and prints the content.

# Parameters:
  - collector *CollectorStruct

The struct containing the scraper and all used context.
*/
func readAndPrint(collector *CollectorStruct) {
	stringSlice := readFinished(collector)
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
func claimNewIndex(collector *CollectorStruct, url string) int {
	index := <-collector.Context.EmptyIndexes
	collector.Context.WorkspaceBuffer[index] = []string{url}
	return index
}

/*
readFinished
retrieves a fully scraped page and returns it.

# Parameters:
  - collector *CollectorStruct

The struct containing the scraper and all used context.

# Returns:
  - []string

The slice containing the text from the scraped page.
*/
func readFinished(collector *CollectorStruct) []string {
	// removes 1 from finished
	countLock := &collector.CounterLock
	countLock.Lock()
	// can currently become negative by this which
	// isn't a case deeply explored but should work fine
	// TODO: test for cases where FinishedCounter becomes negative
	collector.Context.FinishedCounter--
	countLock.Unlock()

	// waits for a Workspace to finish
	index := <-collector.Context.FinishedIndexes

	// retrieves the content and empties out the workspace
	pos := &collector.Context.WorkspaceBuffer[index]
	PageText := *pos
	*pos = nil

	// adds the index to the list of unused/empty workspaces
	collector.Context.EmptyIndexes <- index

	return PageText
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
func writeToWorkspace(collector *CollectorStruct, index int, text string) {
	path := &collector.Context.WorkspaceBuffer[index]
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
func collectorSetup() *CollectorStruct {
	collector := &CollectorStruct{}
	collector.Context = Context{}
	context := &collector.Context
	context.FinishedCounter = 0
	context.WorkingCounter = 0
	context.EmptyIndexes = make(chan int, QUEUEMAXLEN)
	for x := range QUEUEMAXLEN {
		context.EmptyIndexes <- x
	}
	context.FinishedIndexes = make(chan int, QUEUEMAXLEN)
	context.LinkQueue = make(chan string, LINKQUEUELEN)
	context.PriorityLinkQueue = make(chan string, LINKQUEUELEN)
	context.WorkspaceBuffer = [QUEUEMAXLEN][]string{}

	ShortendLinkRegex, _ = regexp.Compile(`^/wiki/`)
	NonAllowedRegex, _ = regexp.Compile(`/wiki/(File|Wikipedia|Special|User):`)

	c := colly.NewCollector(colly.AllowedDomains("en.wikipedia.org"))
	collector.CollectorColly = c

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
		if ShortendLinkRegex.MatchString(url) {
			url = "https://en.wikipedia.org" + url
		}
		index := claimNewIndex(collector, url)
		r.Ctx.Put(INDEXKEY, index)
	})

	// triggered when a CSS selector matches an element
	c.OnHTML("p, div.mw-heading", func(e *colly.HTMLElement) {
		// printing all URLs associated with the <p> tag on the page
		mapValue := e.Response.Ctx.GetAny(INDEXKEY)
		index, ok := mapValue.(int)
		if !ok {
			log.Fatal("couldn't find index")
		}
		writeToWorkspace(collector, index, e.Text)

	})

	// Find and visit all links
	c.OnHTML(`a[href*="/wiki/"]`, func(e *colly.HTMLElement) {
		AddScrapedLinkToQueue(e, collector)
	})

	// triggered once scraping is done
	c.OnScraped(func(r *colly.Response) {
		mapValue := r.Ctx.GetAny(INDEXKEY)
		index, ok := mapValue.(int)
		if !ok {
			log.Fatal("couldn't find index")
		}
		lock := &collector.CounterLock
		lock.Lock()
		context.WorkingCounter--
		context.FinishedCounter++
		lock.Unlock()
		context.FinishedIndexes <- index

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
func AddScrapedLinkToQueue(e *colly.HTMLElement, collector *CollectorStruct) {

	// TODO: filter away already visited before adding to the queue

	context := collector.Context
	link := e.Attr("href")
	if !ShortendLinkRegex.MatchString(link) || NonAllowedRegex.MatchString(link) {
		debugPrint("non allowed link: ", link)
		return
	}
	link = "https://en.wikipedia.org" + link
	select {
	case context.LinkQueue <- link:
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
	if len(collector.Context.LinkQueue) == 0 {
		collector.Context.LinkQueue <- link
		return
	}
	collector.Context.PriorityLinkQueue <- link
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
func VisitNextLink(collector *CollectorStruct) {
	for {
		var link string
		select {
		case link = <-collector.Context.PriorityLinkQueue:
		default:
			link = <-collector.Context.LinkQueue
		}
		err := collector.CollectorColly.Visit(link)
		if err == nil {

			collector.Context.WorkingCounter++
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
