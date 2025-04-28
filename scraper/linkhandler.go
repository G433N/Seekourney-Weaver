package scraper

import (
	"log"
	"seekourney/utils/Sync"
	"sync"
)

const _PRIOQUEUEMAXLEN_ = 100
const _STACKMAXLEN_ = 1_000_000

type (
	hostPath   string
	innerPath  string
	URLCompact struct {
		webFileBool bool
		host        int
		inner       innerPath
	}

	filter struct {
		filehosts []hostPath
		webhosts  []hostPath
		filterMap map[hostPath]filterMapInner
	}
	filterMapInner struct {
		webFileBool bool
		index       int
		filterMap   map[innerPath]bool
	}

	linkInputWrap struct {
		prio bool
		URL  *URLCompact
	}

	PriorityQueue struct {
		lock sync.Mutex

		Queue [_PRIOQUEUEMAXLEN_]*URLCompact

		read int

		write int

		len int
	}
	storageStack struct {
		lock sync.Mutex

		Stack []*URLCompact
	}

	linkHandler struct {
		filter       filter
		storageStack storageStack

		priorityQueue PriorityQueue

		inputChan chan linkInputWrap

		outputChan chan URLString

		outputSem Sync.Semaphore

		storedSem       Sync.Semaphore
		quit            bool
		handlersWorking sync.WaitGroup
	}
)

func (stack *storageStack) pop() (*URLCompact, bool) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	if len(stack.Stack) == 0 {
		return nil, false
	}
	URL := stack.Stack[len(stack.Stack)-1]
	stack.Stack = stack.Stack[:len(stack.Stack)-1]

	return URL, true
}

func (stack *storageStack) push(URL *URLCompact) bool {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	if len(stack.Stack) >= _STACKMAXLEN_ {
		return false
	}
	stack.Stack = append(stack.Stack, URL)
	return true
}

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
		fullURL = "https://" + URLString(filter.webhosts[compactURL.host])
	} else {
		fullURL = "file://" + URLString(filter.filehosts[compactURL.host])
	}
	fullURL += URLString(compactURL.inner)
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
			URL, sucess = lH.storageStack.pop()
		}
		if !sucess {
			log.Fatal(
				"Couldn't find any URL",
				"even though the semaphore had a signal",
			)
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
		notFull = lH.storageStack.push(URL)
		if notFull {
			lH.storedSem.Signal()
			continue
		}
		debugPrint("stack full skipping: ", URL)
	}

}

func linkHandlerCreate() *linkHandler {

	lH := linkHandler{
		storageStack: storageStack{
			Stack: []*URLCompact{},
		},
		priorityQueue: PriorityQueue{
			Queue: [_PRIOQUEUEMAXLEN_]*URLCompact{},
			read:  0,
			write: 0,
			len:   0,
		},
		inputChan:       make(chan linkInputWrap),
		outputChan:      make(chan URLString),
		outputSem:       *Sync.NewSemaphore(0),
		storedSem:       *Sync.NewSemaphore(0),
		quit:            false,
		handlersWorking: sync.WaitGroup{},
	}
	lH.handlersWorking.Add(2)
	go lH.inputHandler()
	go lH.outputHandler()

	return &lH

}

func (lH *linkHandler) destroy() {
	lH.quit = true
	close(lH.inputChan)
	lH.handlersWorking.Wait()
}

func (lH *linkHandler) input(
	hostPath hostPath,
	innerPath innerPath,
	prio bool,
) {

}
