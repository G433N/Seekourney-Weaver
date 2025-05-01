package Sync

import "sync"

const _STACKMAXLEN_ = 1_000_000

type Stack[T any] struct {
	lock  sync.Mutex
	stack []T
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{
		stack: []T{},
	}
}

func (stack *Stack[T]) IsEmpty() bool {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	return len(stack.stack) == 0
}

func (stack *Stack[T]) Pop() (T, bool) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	if len(stack.stack) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	URL := stack.stack[len(stack.stack)-1]
	stack.stack = stack.stack[:len(stack.stack)-1]

	return URL, true
}

func (stack *Stack[T]) Push(elem T) bool {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	if len(stack.stack) >= _STACKMAXLEN_ {
		return false
	}
	stack.stack = append(stack.stack, elem)
	return true
}
