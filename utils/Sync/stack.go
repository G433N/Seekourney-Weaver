package Sync

import (
	"log"
	"sync"
)

// currently unused
//const _STACKMAXLEN_ = 1_000_000

type Stack[T any] struct {
	lock   sync.Mutex
	stack  []T
	lenSem Semaphore
}

func NewStack[T any](optionalSize ...int) Stack[T] {
	if len(optionalSize) == 1 {
		return Stack[T]{
			stack:  make([]T, 0, optionalSize[0]),
			lenSem: *NewSemaphore(0),
		}
	}
	if len(optionalSize) > 1 {
		log.Fatalln("NewStack", "only one size argument is allowed")
	}
	return Stack[T]{
		stack:  []T{},
		lenSem: *NewSemaphore(0),
	}
}

func (stack *Stack[T]) IsEmpty() bool {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	return len(stack.stack) == 0
}

func (stack *Stack[T]) Pop() T {
	stack.lenSem.Wait()
	stack.lock.Lock()
	URL := stack.stack[len(stack.stack)-1]
	stack.stack = stack.stack[:len(stack.stack)-1]
	stack.lock.Unlock()
	return URL
}

func (stack *Stack[T]) TryPop() (T, bool) {
	if !stack.lenSem.TryWait() {
		var zeroValue T
		return zeroValue, false
	}
	stack.lock.Lock()
	URL := stack.stack[len(stack.stack)-1]
	stack.stack = stack.stack[:len(stack.stack)-1]
	stack.lock.Unlock()
	return URL, true
}

func (stack *Stack[T]) Push(elem T) {
	stack.lenSem.Signal()
	stack.lock.Lock()
	stack.stack = append(stack.stack, elem)
	stack.lock.Unlock()
}

func (stack *Stack[T]) TryPush(elem T) bool {
	stack.lock.Lock()
	if !stack.lenSem.TrySignal() {
		return false
	}

	stack.stack = append(stack.stack, elem)
	stack.lock.Unlock()
	return true
}
