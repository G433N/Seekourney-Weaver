package Sync

import (
	"sync"
)

type CyclicQueue[T any] struct {
	lock sync.Mutex

	queue []T

	read int

	write int

	currentLen int

	emptySem *Semaphore

	filledSem *Semaphore
}

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

func (queue *CyclicQueue[T]) Push(URL T) {
	queue.emptySem.Wait()
	queue.lock.Lock()
	defer queue.lock.Unlock()
	queue.queue[queue.write] = URL
	queue.write = (queue.write + 1) % len(queue.queue)
	queue.filledSem.Signal()
	return
}

func (queue *CyclicQueue[T]) TryPush(URL T) bool {
	if !queue.emptySem.TryWait() {
		return false
	}
	queue.lock.Lock()
	defer queue.lock.Unlock()
	queue.queue[queue.write] = URL
	queue.write = (queue.write + 1) % len(queue.queue)
	queue.currentLen++
	queue.filledSem.Signal()
	return true
}

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

func (queue *CyclicQueue[T]) Pop() T {
	queue.filledSem.Wait()
	queue.lock.Lock()
	defer queue.lock.Unlock()
	URL := queue.queue[queue.read]
	queue.read = (queue.read + 1) % len(queue.queue)
	queue.emptySem.Signal()

	return URL
}
