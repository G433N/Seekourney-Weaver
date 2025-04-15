package Sync

import (
	"sync"
)

type Semaphore struct {
	signals   int
	syncLock  sync.Mutex
	waitGroup sync.WaitGroup
}

func (semaphore *Semaphore) Wait() {
	semaphore.waitGroup.Wait()
	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()

	if semaphore.signals == 1 {
		semaphore.waitGroup.Add(1)
	}
	semaphore.signals--

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
