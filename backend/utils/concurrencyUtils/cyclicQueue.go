package Sync

import (
	"sync"
)

/*
CyclicQueue is a thread-safe generic cyclic queue implementation.
Can Pop, Push and has non-blocking versions too.
*/
type CyclicQueue[T any] struct {

	// lock is used to protect access to the queue
	lock sync.Mutex

	// queue is the underlying slice that holds the elements
	queue []T

	// read is the index of the next element to be read
	read int

	// write is the index where the next element can be written
	write int

	//  emptySem tracks the number of empty slots in the queue
	emptySem *Semaphore

	// filledSem tracks the number of filled slots in the queue
	filledSem *Semaphore
}

/*
NewCyclicQueue
creates a new CyclicQueue with the given size.
*/
func NewCyclicQueue[T any](size int) CyclicQueue[T] {
	return CyclicQueue[T]{
		lock:      sync.Mutex{},
		queue:     make([]T, size),
		read:      0,
		write:     0,
		emptySem:  NewSemaphore(size),
		filledSem: NewSemaphore(0),
	}
}

/*
Push
adds an element to the queue.
If the queue is full, it blocks until space is available.
*/
func (cq *CyclicQueue[T]) Push(elem T) {
	cq.emptySem.Wait()
	cq.lock.Lock()
	defer cq.lock.Unlock()
	cq.queue[cq.write] = elem
	cq.write = (cq.write + 1) % len(cq.queue)
	cq.filledSem.Signal()
}

/*
TryPush
attempts to add an element to the queue.
If the queue is full, it returns false without blocking.
*/
func (cq *CyclicQueue[T]) TryPush(elem T) bool {
	if !cq.emptySem.TryWait() {
		return false
	}
	cq.lock.Lock()
	defer cq.lock.Unlock()
	cq.queue[cq.write] = elem
	cq.write = (cq.write + 1) % len(cq.queue)
	cq.filledSem.Signal()
	return true
}

/*
Pop
retrieves an element from the queue.
If the queue is empty, it blocks until an element is available.
*/
func (cq *CyclicQueue[T]) Pop() T {
	cq.filledSem.Wait()
	cq.lock.Lock()
	defer cq.lock.Unlock()
	elem := cq.queue[cq.read]
	cq.read = (cq.read + 1) % len(cq.queue)
	cq.emptySem.Signal()

	return elem
}

/*
TryPop
attempts to retrive an element from the queue.
If the queue is empty, it returns false without blocking.
*/
func (cq *CyclicQueue[T]) TryPop() (T, bool) {
	if !cq.filledSem.TryWait() {
		var zeroValue T
		return zeroValue, false
	}
	cq.lock.Lock()
	defer cq.lock.Unlock()
	elem := cq.queue[cq.read]
	cq.read = (cq.read + 1) % len(cq.queue)
	cq.emptySem.Signal()
	return elem, true
}
