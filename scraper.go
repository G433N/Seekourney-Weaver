package main

// Start test https://en.wikipedia.org/wiki/Cucumber

// End goal https://tracklock.gg/players/120846718

import (
	"fmt"
	"log"
	"sync"

	"github.com/gocolly/colly/v2"
)

const QUEUEMAXLEN = 5
const INDEXKEY = `QueueIndex`

// initialize a data structure to keep the scraped data
type Context struct {
	// currently working on amount
	workingCounter int
	// currently finished amount
	finishedCounter int
	// queue channel of links to visit
	// TODO: filter away already visited and illegal before adding to the queue
	linkQueue chan string
	// finished queue and working space
	finished *[QUEUEMAXLEN][]string
	// queue index where to put new ones
	queueEndIndex int
	// queue index where to read from
	queueStartIndex int
}

type CollectorStruct struct {
	FinishedQueueLock sync.Mutex
	LinkQueueLock     sync.Mutex
	CounterLock       sync.Mutex

	context        Context
	collectorColly *colly.Collector
}

func main2() {
	collector := collectorSetup()
	//temp := []string{}

	//... scraping logic

	fmt.Println("Start")
	// open the target URL
	err := c.Visit("https://en.wikipedia.org/wiki/Cucumber")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Waaa")

	c.Wait()
	close(scrapedLinks)
	for elem := range scrapedLinks {
		fmt.Println("hihi", elem)
	}
	fmt.Println("done")

}
func collectorRepopulateQueue(collector *CollectorStruct) {

}

func claimNewIndex(collector *CollectorStruct, url string) int {
	lock := &collector.FinishedQueueLock
	lock.Lock()
	defer lock.Unlock()
	index := collector.context.queueEndIndex
	collector.context.queueEndIndex = (index + 1) % QUEUEMAXLEN
	collector.context.finished[index] = []string{url}
	return index
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
	context.queueEndIndex = 0
	context.queueStartIndex = 0
	context.linkQueue = make(chan string, 100)

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
		FollowLink(e, collector)
	})

	// triggered once scraping is done (e.g., write the data to a CSV file)
	c.OnScraped(func(r *colly.Response) {
		//finished <- temp
		//temp = []string{}
		scrapedLinks <- r.Request.URL.Path
		fmt.Println(r.Request.URL, " scraped!")
	})
	return c
}

func FollowLink(e *colly.HTMLElement, collector *CollectorStruct) {
	counterLock := &collector.CounterLock
	counterLock.Lock()
	if *counter <= 0 {
		counterLock.Unlock()
		return
	}

	link := e.Attr("href")
	err := e.Request.Visit(link)
	if err == nil {
		*counter -= 1
		counterLock.Unlock()
		return
	}
	counterLock.Unlock()

	switch err.Error() {
	case `Forbidden domain`:
		fmt.Println(err, link)
	case `URL already visited`:
		fmt.Println(err, link)
	default:
		fmt.Println(err)
	}
}
