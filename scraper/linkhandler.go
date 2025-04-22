package scraper

import "seekourney/utils/Sync"

type (
	linkInputWrap struct {
		prio bool
		URL  URLString
	}

	linkHandler struct {
		storageStack []URLString

		priorityQueue []URLString

		inputChan chan linkInputWrap

		outputChan chan URLString

		outputSem Sync.Semaphore
		quit      bool
	}
)

func (lH *linkHandler) inputHandler() {
	for {
		inputWrap, more := <-lH.inputChan
		_ = inputWrap
		if !more {
			break
		}

	}
}

func linkHandlerCreate() *linkHandler {

	lH := linkHandler{}

	return &lH

}

func (linkHandler *linkHandler) destroy() {

}

func (linkHandler *linkHandler) input(URL URLString) {

}
