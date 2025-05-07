package Sync

import "sync"

type ArrayPlus[T any] struct {
	lock sync.Mutex
	gaps Stack[int]
	arr  []T
}

func NewArrayPlus[T any](size int) *ArrayPlus[T] {
	newArr := ArrayPlus[T]{
		gaps: NewStack[int](),
		arr:  make([]T, size),
	}
	for i := range size {
		newArr.gaps.Push(i)
	}
	return &newArr
}

func (stack *ArrayPlus[T]) Pop(index int) T {
	stack.lock.Lock()

	if index < 0 || index >= len(stack.arr) {
		var zeroValue T

		return zeroValue, false
	}
	elem := stack.arr[index]

	stack.lock.Unlock()
	stack.gaps.Push(index)

	return elem, true
}
func (stack *ArrayPlus[T]) TryPop(index int) (T, bool) {
	stack.lock.Lock()

	if index < 0 || index >= len(stack.arr) {
		var zeroValue T

		return zeroValue, false
	}
	elem := stack.arr[index]

	stack.lock.Unlock()
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
