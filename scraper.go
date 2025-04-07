package main

// Start test https://en.wikipedia.org/wiki/Cucumber

// End goal https://tracklock.gg/players/120846718

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
	workingCounter int
	// currently finished amount should be in sync with len(finishedIndexes)
	finishedCounter int
	// buffered queue channel of links scraped from previously visited sites
	linkQueue chan string
	// worspace for the collectors async requests
	// each request gets their own index and append the text they recieve to the slice
	Workspace [QUEUEMAXLEN][]string
	// channel of indexes in the 'workspace' array ready to be assigned to a request
	// (buffer of size 'QUEUEMAXLEN')
	EmptyIndexes chan int
	//  channel of indexes in the 'workspace' array ready to be read
	// (buffer of size 'QUEUEMAXLEN')
	finishedIndexes chan int
}

type CollectorStruct struct {
	// mutex used to sync changes to the two counters in context
	CounterLock sync.Mutex

	// struct holding all cotext to make the inteface with the collector as simple as possible
	context Context

	// the colly collector used for webb scraping and formatting
	collectorColly *colly.Collector
}

func debugPrint(a ...any) {
	if DEBUG {
		fmt.Println(a...)
	}
}

func main() {
	collector := collectorSetup()

	collector.context.linkQueue <- "https://en.wikipedia.org/wiki/Cucumber"
	collectorRepopulateQueue(collector)

	readAndPrint(collector)

	readAndPrint(collector)

	readAndPrint(collector)

}
func collectorRepopulateQueue(collector *CollectorStruct) {
	context := &collector.context
	countLock := &collector.CounterLock
	countLock.Lock()
	amountFilled := context.finishedCounter + context.workingCounter
	amountEmpty := QUEUEMAXLEN - amountFilled
	for range amountEmpty {
		VisitNextLink(collector)
	}
	countLock.Unlock()

}

func readAndPrint(collector *CollectorStruct) {
	stringSlice := readFinished(collector)
	for _, text := range stringSlice {
		fmt.Println(text)
	}
}

func claimNewIndex(collector *CollectorStruct, url string) int {
	index := <-collector.context.EmptyIndexes
	collector.context.Workspace[index] = []string{url}
	return index
}

func readFinished(collector *CollectorStruct) []string {
	countLock := &collector.CounterLock
	countLock.Lock()
	collector.context.finishedCounter--
	countLock.Unlock()

	index := <-collector.context.finishedIndexes
	collector.context.EmptyIndexes <- index

	pos := &collector.context.Workspace[index]
	result := *pos
	*pos = nil
	return result
}

func writeToWorkspace(collector *CollectorStruct, index int, text string) {
	path := &collector.context.Workspace[index]
	*path = append(*path, text)
}

func collectorSetup() *CollectorStruct {
	collector := &CollectorStruct{}
	collector.context = Context{}
	context := &collector.context
	context.finishedCounter = 0
	context.workingCounter = 0
	context.EmptyIndexes = make(chan int, QUEUEMAXLEN)
	for x := range QUEUEMAXLEN {
		context.EmptyIndexes <- x
	}
	context.finishedIndexes = make(chan int, QUEUEMAXLEN)
	context.linkQueue = make(chan string, LINKQUEUELEN)
	context.Workspace = [QUEUEMAXLEN][]string{}

	ShortendLinkRegex, _ = regexp.Compile(`^/wiki/`)
	NonAllowedRegex, _ = regexp.Compile(`/wiki/(File|Wikipedia|Special|User):`)

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
		AddLinkToQueue(e, collector)
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
		context.workingCounter--
		context.finishedCounter++
		lock.Unlock()
		context.finishedIndexes <- index

	})
	return collector
}

func AddLinkToQueue(e *colly.HTMLElement, collector *CollectorStruct) {

	// TODO: filter away already visited before adding to the queue

	context := collector.context
	link := e.Attr("href")
	if !ShortendLinkRegex.MatchString(link) || NonAllowedRegex.MatchString(link) {
		debugPrint("non allowed link: ", link)
		return
	}
	link = "https://en.wikipedia.org" + link
	select {
	case context.linkQueue <- link:
	default:
		debugPrint("linkQueue full skipped link: ", link)
	}

}

// should only be called while having the counter lock
func VisitNextLink(collector *CollectorStruct) {
	for {
		link := <-collector.context.linkQueue
		err := collector.collectorColly.Visit(link)
		if err == nil {

			collector.context.workingCounter++
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
