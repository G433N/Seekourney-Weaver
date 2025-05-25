package scraper

import (
	"seekourney/utils/concurrencyUtils"
	"sync"
)

func (filter *filter) toURLString(compactURL *URLCompact) URLString {
	var fullURL URLString

	if compactURL.webFileBool {
		elem := filter.webhosts.Peek(compactURL.host)
		fullURL = "https://" + URLString(elem)
	} else {
		fullURL = "file://"
	}
	fullURL += URLString(compactURL.inner)
	debugPrint("URL: ", fullURL)
	return fullURL
}

func (lH *linkHandler) outputHandler() {
	defer lH.handlersWorking.Done()
	for {
		lH.outputSem.Wait()
		if lH.quit {
			break
		}
		var URL *URLCompact
		var sucess bool
		lH.storedSem.Wait()
		URL, sucess = lH.priorityQueue.TryPop()
		if !sucess {
			URL = lH.storageStack.Pop()
		}

		lH.outputChan <- lH.filter.toURLString(URL)
	}
	close(lH.outputChan)
}

func (lH *linkHandler) inputHandler() {
	defer lH.handlersWorking.Done()
	for {
		inputWrap, more := <-lH.inputChan
		if !more {
			break
		}
		prio := inputWrap.prio
		URL := inputWrap.URL
		notFull := true
		if prio {
			notFull = lH.priorityQueue.TryPush(URL)
			if notFull {
				lH.storedSem.Signal()
				continue
			}
			debugPrint("PriorityQueue full, redirecting to stack: ", URL)
		}
		notFull = lH.storageStack.TryPush(URL)
		if notFull {
			lH.storedSem.Signal()
			continue
		}
		debugPrint("stack full skipping: ", URL)
	}

}

func linkHandlerCreate() *linkHandler {

	lH := linkHandler{
		storageStack: *concurrencyUtils.NewStack[*URLCompact](),
		priorityQueue: concurrencyUtils.NewCyclicQueue[*URLCompact](
			_PRIOQUEUEMAXLEN_,
		),
		inputChan:       make(chan linkInputWrap),
		outputChan:      make(chan URLString),
		outputSem:       *concurrencyUtils.NewSemaphore(0),
		storedSem:       *concurrencyUtils.NewSemaphore(0),
		quit:            false,
		handlersWorking: sync.WaitGroup{},
		filter: filter{
			webhosts:  *concurrencyUtils.NewArrayPlus[hostPath](1000),
			filterMap: map[hostPath]filterMapInner{},
		},
	}
	lH.filter.Whitelist(`file://`, false)
	lH.handlersWorking.Add(2)
	go lH.inputHandler()
	go lH.outputHandler()

	return &lH

}

func (lH *linkHandler) getLink() URLString {
	lH.outputSem.Signal()
	return <-lH.outputChan

}

func (lH *linkHandler) destroy() {
	lH.quit = true
	close(lH.inputChan)
	lH.handlersWorking.Wait()
}

func (lH *linkHandler) inputLocalFile(
	URL URLString,
	prio bool,
) {
	filter := &lH.filter
	newURL, ok := filter.URLCompactCreate("file://", innerPath(URL))
	if !ok {
		return
	}
	lH.inputChan <- linkInputWrap{
		prio: prio,
		URL:  newURL,
	}
}

func (lH *linkHandler) inputWeb(
	hostPath hostPath,
	innerPath innerPath,
	prio bool,
) {
	URL, ok := lH.filter.URLCompactCreate(hostPath, innerPath)
	if !ok {
		return
	}

	switch hostPath {
	case `en.wikipedia.org`:
		if WikipediaBadRegex.MatchString(string(innerPath)) {
			debugPrint("Not worth indexing: ", innerPath)
			return
		}
	default:
	}

	lH.inputChan <- linkInputWrap{
		prio: prio,
		URL:  URL,
	}
}

func (filter *filter) URLCompactCreate(
	hostPath hostPath,
	innerPath innerPath,
) (*URLCompact, bool) {

	val, ok := filter.filterMap[hostPath]
	if !ok {
		debugPrint("Host:", hostPath, "not whitelisted")
		return nil, false
	}
	alreadyDone := val.filterMap[innerPath]
	if alreadyDone {
		debugPrint(
			"URL:",
			"https://"+string(hostPath)+string(innerPath),
			"already visited",
		)
		return nil, false
	}
	val.filterMap[innerPath] = true
	return &URLCompact{
		webFileBool: val.webFileBool,
		host:        val.index,
		inner:       innerPath,
	}, true
}

func (filter *filter) Whitelist(hostPath hostPath, webFileBool bool) {
	index, ok := filter.webhosts.TryPush(hostPath)
	if !ok {
		debugPrint("webhost array full")
		return
	}
	filter.filterMap[hostPath] = filterMapInner{
		webFileBool: webFileBool,
		index:       index,
		filterMap:   map[innerPath]bool{},
	}
}
