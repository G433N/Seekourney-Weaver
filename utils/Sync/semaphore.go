package Sync

import (
	"sync"
)

/*
A very basic semaphore implementation
*/
type Semaphore struct {
	//internal value to keep track of amount signals
	value int
	// mutex used for the internal value
	syncLock sync.Mutex
	// mutex used for the waiting group
	waitLock sync.Mutex

	waitGroup sync.WaitGroup
}

/*
Wait
Decrement the semaphore’s value; block if the value is currently 0.
*/
func (semaphore *Semaphore) Wait() {
	// so that multiple threads don't use the waitgroup at the same time
	semaphore.waitLock.Lock()
	defer semaphore.waitLock.Unlock()

	semaphore.waitGroup.Wait()

	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()

	if semaphore.value == 1 {
		semaphore.waitGroup.Add(1)
	}

	semaphore.value--

}

/*
TryWait
Tries to decrement the semaphore’s value and return sucess/true.
If the value is 0 it will return false.

# Returns:

The amount of requests that couldn't be fullfilled.
*/
func (semaphore *Semaphore) TryWait() bool {
	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()
	if semaphore.value == 0 {
		return false
	}
	if semaphore.value == 1 {
		semaphore.waitGroup.Add(1)
	}
	semaphore.value--
	return true
}

/*
Signal
Increments the semaphores value.
If it was 0 it will unblock any waiting threads.
*/
func (semaphore *Semaphore) Signal() {
	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()
	if semaphore.value == 0 {
		semaphore.waitGroup.Done()
	}
	semaphore.value++
}

/*
NewSemaphore
creates a new semaphore and sets the initial value
*/
func NewSemaphore(initialValue int) *Semaphore {
	semaphore := Semaphore{}
	if initialValue == 0 {
		semaphore.waitGroup.Add(1)
	}
	semaphore.value = initialValue
	return &semaphore
}
