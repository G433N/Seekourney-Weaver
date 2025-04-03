package main

// Start test https://en.wikipedia.org/wiki/Cucumber

// End goal https://tracklock.gg/players/120846718

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
)

// initialize a data structure to keep the scraped data
type Context struct {
	counter      int
	scrapedLinks chan string
	finished     chan []string
	linkQueue    chan string
	maxAmount    int
}

type CollectorStruct struct {
	ContextLock sync.Mutex
	context     Context
	collector   colly.Collector
}

func main() {
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

func collectorSetup() CollectorStruct {
	var counterLock sync.Mutex
	var counter = 5
	//linkQueue := make(chan string, 100)
	//finished := make(chan []string, 100)
	scrapedLinks := make(chan string, 100)

	c := colly.NewCollector(colly.AllowedDomains("en.wikipedia.org"))
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
	})

	// triggered when a CSS selector matches an element
	c.OnHTML("p, div.mw-heading", func(e *colly.HTMLElement) {
		// printing all URLs associated with the <p> tag on the page
		//fmt.Println(e.Text)

	})

	// Find and visit all links
	c.OnHTML(`a[href*="/wiki/"]`, func(e *colly.HTMLElement) {
		FollowLink(e, &counter, &counterLock)
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

func FollowLink(e *colly.HTMLElement, counter *int, counterLock *sync.Mutex) {
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
