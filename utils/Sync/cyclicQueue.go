package Sync

import "sync"

type CyclicQueue[T any] struct {
	lock sync.Mutex

	queue []T

	read int

	write int

	currentLen int
}

func NewCyclicQueue[T any](size int) CyclicQueue[T] {
	return CyclicQueue[T]{
		lock:       sync.Mutex{},
		queue:      make([]T, size),
		read:       0,
		write:      0,
		currentLen: 0,
	}
}

func (queue *CyclicQueue[T]) Push(URL T) bool {
	queue.lock.Lock()
	defer queue.lock.Unlock()
	if queue.currentLen == len(queue.queue) {
		return false
	}
	queue.queue[queue.write] = URL
	queue.write = (queue.write + 1) % len(queue.queue)
	queue.currentLen++
	return true
}

func (queue *CyclicQueue[T]) Pop() (T, bool) {
	queue.lock.Lock()
	defer queue.lock.Unlock()
	if queue.currentLen == 0 {
		var zeroValue T
		return zeroValue, false
	}
	URL := queue.queue[queue.read]
	queue.read = (queue.read + 1) % len(queue.queue)
	queue.currentLen--
	return URL, true
}
