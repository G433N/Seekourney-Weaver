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

const QUEUEMAXLEN = 5
const LINKQUEUELEN = 200
const INDEXKEY = `QueueIndex`

var AllowedDomainsRegex *regexp.Regexp
var NonAllowedRegex *regexp.Regexp

// initialize a data structure to keep the scraped data
type Context struct {
	// currently working on amount
	workingCounter int
	// currently finished amount
	finishedCounter int
	// queue channel of links to visit
	// TODO: filter away already visited before adding to the queue
	linkQueue chan string
	// finished queue and working space
	finished [QUEUEMAXLEN][]string
	// queue index where to put new ones
	EmptyIndexes chan int
	// queue index where to read from
	finishedIndexes chan int
}

type CollectorStruct struct {
	FinishedQueueLock sync.Mutex
	LinkQueueLock     sync.Mutex
	CounterLock       sync.Mutex

	context        Context
	collectorColly *colly.Collector
}

func main() {
	collector := collectorSetup()
	//temp := []string{}

	//... scraping logic

	fmt.Println("Start")
	collector.context.linkQueue <- "https://en.wikipedia.org/wiki/Cucumber"
	VisitNextLink(collector)

	VisitNextLink(collector)

	readAndPrint(collector)

	fmt.Print("\n\nhiohio\n\n\n")
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
	lock := &collector.FinishedQueueLock
	lock.Lock()
	defer lock.Unlock()
	index := <-collector.context.EmptyIndexes
	collector.context.finished[index] = []string{url}
	return index
}

func readFinished(collector *CollectorStruct) []string {
	countLock := &collector.CounterLock
	countLock.Lock()
	collector.context.finishedCounter--
	countLock.Unlock()
	FQLock := &collector.FinishedQueueLock

	index := <-collector.context.finishedIndexes
	collector.context.EmptyIndexes <- index

	FQLock.Lock()
	defer FQLock.Unlock()
	pos := &collector.context.finished[index]
	result := *pos
	*pos = nil
	return result
}

func writeToWorkspace(collector *CollectorStruct, index int, text string) {
	path := &collector.context.finished[index]
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
	context.finished = [QUEUEMAXLEN][]string{}

	AllowedDomainsRegex, _ = regexp.Compile(`^/wiki/`)
	NonAllowedRegex, _ = regexp.Compile(`/wiki/(File|Wikipedia|Special|User):`)

	c := colly.NewCollector(colly.AllowedDomains("en.wikipedia.org"))
	collector.collectorColly = c

	c.Async = true

	// called before an HTTP request is triggered
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	// triggered when the scraper encounters an error
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	// fired when the server responds
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
		index := claimNewIndex(collector, r.Request.URL.EscapedPath())
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
	context := collector.context
	link := e.Attr("href")
	if !AllowedDomainsRegex.MatchString(link) || NonAllowedRegex.MatchString(link) {
		fmt.Println("non allowed link: ", link)
		return
	}
	link = "https://en.wikipedia.org" + link
	select {
	case context.linkQueue <- link:
	default:
		fmt.Println("linkQueue full skipped link: ", link)
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
			fmt.Println(err, link)
		case `URL already visited`:
			fmt.Println(err, link)
		default:
			fmt.Println(err)
		}
	}
}
