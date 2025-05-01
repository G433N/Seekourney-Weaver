package Sync

import "sync"

const _ARRAYPLUSMAXLEN_ = 1_000_000

type ArrayPlus[T any] struct {
	lock sync.Mutex
	gaps Stack[int]
	arr  []T
}

func NewArrayPlus[T any]() ArrayPlus[T] {
	return ArrayPlus[T]{
		gaps: NewStack[int](),
		arr:  []T{},
	}
}

func (stack *ArrayPlus[T]) Pop(index int) (T, bool) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	if index < 0 || index >= len(stack.arr) {
		var zeroValue T

		return zeroValue, false
	}
	elem := stack.arr[index]
	stack.gaps.Push(index)

	return elem, true
}

func (stack *ArrayPlus[T]) Peek(index int) (T, bool) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	if index < 0 || index >= len(stack.arr) {
		var zeroValue T
		return zeroValue, false
	}
	elem := stack.arr[index]

	return elem, true
}

func (stack *ArrayPlus[T]) Push(elem T) (int, bool) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	index, ok := stack.gaps.Pop()
	if ok {
		stack.arr[index] = elem
		return index, true
	}
	if len(stack.arr) >= _ARRAYPLUSMAXLEN_ {
		return 0, false
	}
	stack.arr = append(stack.arr, elem)
	return len(stack.arr) - 1, true
}
