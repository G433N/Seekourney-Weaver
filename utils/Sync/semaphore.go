package Sync

import (
	"sync"
)

type Semaphore struct {
	signals  int
	syncLock sync.Mutex
	waitLock sync.Mutex

	waitGroup sync.WaitGroup
}

func (semaphore *Semaphore) Wait() {
	semaphore.waitLock.Lock()
	defer semaphore.waitLock.Unlock()

	semaphore.waitGroup.Wait()

	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()

	if semaphore.signals == 1 {
		semaphore.waitGroup.Add(1)
	}

	semaphore.signals--

}

func (semaphore *Semaphore) TryWait() bool {
	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()
	if semaphore.signals == 0 {
		return false
	}
	if semaphore.signals == 1 {
		semaphore.waitGroup.Add(1)
	}
	semaphore.signals--
	return true
}

func (semaphore *Semaphore) Signal() {
	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()
	if semaphore.signals == 0 {
		semaphore.waitGroup.Done()
	}
	semaphore.signals++
}

func NewSemaphore(initialValue int) *Semaphore {
	semaphore := Semaphore{}
	if initialValue == 0 {
		semaphore.waitGroup.Add(1)
	}
	semaphore.signals = initialValue
	return &semaphore
}
