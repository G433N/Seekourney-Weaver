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
func (queue *CyclicQueue[T]) Push(URL T) {
	queue.emptySem.Wait()
	queue.lock.Lock()
	defer queue.lock.Unlock()
	queue.queue[queue.write] = URL
	queue.write = (queue.write + 1) % len(queue.queue)
	queue.filledSem.Signal()
}

/*
TryPush
attempts to add an element to the queue.
If the queue is full, it returns false without blocking.
*/
func (queue *CyclicQueue[T]) TryPush(URL T) bool {
	if !queue.emptySem.TryWait() {
		return false
	}
	queue.lock.Lock()
	defer queue.lock.Unlock()
	queue.queue[queue.write] = URL
	queue.write = (queue.write + 1) % len(queue.queue)
	queue.filledSem.Signal()
	return true
}

/*
Pop
retrieves an element from the queue.
If the queue is empty, it blocks until an element is available.
*/
func (queue *CyclicQueue[T]) Pop() T {
	queue.filledSem.Wait()
	queue.lock.Lock()
	defer queue.lock.Unlock()
	URL := queue.queue[queue.read]
	queue.read = (queue.read + 1) % len(queue.queue)
	queue.emptySem.Signal()

	return URL
}

/*
TryPop
attempts to retrive an element from the queue.
If the queue is empty, it returns false without blocking.
*/
func (queue *CyclicQueue[T]) TryPop() (T, bool) {
	if !queue.filledSem.TryWait() {
		var zeroValue T
		return zeroValue, false
	}
	queue.lock.Lock()
	defer queue.lock.Unlock()
	URL := queue.queue[queue.read]
	queue.read = (queue.read + 1) % len(queue.queue)
	queue.emptySem.Signal()
	return URL, true
}
