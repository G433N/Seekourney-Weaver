package scraper

import (
	"seekourney/utils/concurrencyUtils"
	"sync"
)

func (PQ *PriorityQueue) push(URL *URLCompact) bool {
	PQ.lock.Lock()
	defer PQ.lock.Unlock()
	if PQ.len == _PRIOQUEUEMAXLEN_ {
		return false
	}
	PQ.Queue[PQ.write] = URL
	PQ.write = (PQ.write + 1) % _PRIOQUEUEMAXLEN_
	PQ.len++
	return true
}

func (PQ *PriorityQueue) pop() (*URLCompact, bool) {
	PQ.lock.Lock()
	defer PQ.lock.Unlock()
	if PQ.len == 0 {
		return nil, false
	}
	URL := PQ.Queue[PQ.read]
	PQ.read = (PQ.read + 1) % _PRIOQUEUEMAXLEN_
	PQ.len--
	return URL, true
}

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
		URL, sucess = lH.priorityQueue.pop()
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
			notFull = lH.priorityQueue.push(URL)
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
		priorityQueue: PriorityQueue{
			Queue: [_PRIOQUEUEMAXLEN_]*URLCompact{},
			read:  0,
			write: 0,
			len:   0,
		},
		inputChan:       make(chan linkInputWrap),
		outputChan:      make(chan URLString),
		outputSem:       *concurrencyUtils.NewSemaphore(),
		storedSem:       *concurrencyUtils.NewSemaphore(),
		quit:            false,
		handlersWorking: sync.WaitGroup{},
		filter: filter{
			webhosts:  *concurrencyUtils.NewArrayPlus[hostPath](10_000),
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
