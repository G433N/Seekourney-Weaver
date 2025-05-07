package Sync

import "sync"

type ArrayPlusIndex int

type ArrayPlus[T any] struct {
	lock   sync.Mutex
	gaps   Stack[ArrayPlusIndex]
	arr    []T
	maxLen int
}

const _CHUNKSIZE_ = 100

func NewArrayPlus[T any](maxSize int) *ArrayPlus[T] {
	lenght := min(maxSize, _CHUNKSIZE_)
	newArr := ArrayPlus[T]{
		gaps:   NewStack[ArrayPlusIndex](lenght),
		arr:    make([]T, lenght),
		maxLen: maxSize,
	}
	for i := range lenght {
		newArr.gaps.Push(ArrayPlusIndex(i))
	}
	return &newArr
}

func (stack *ArrayPlus[T]) Pop(indexTyped ArrayPlusIndex) T {
	index := int(indexTyped)

	elem := stack.arr[index]

	stack.gaps.Push(indexTyped)

	return elem
}

// with current invariants, this function is not needed
//func (stack *ArrayPlus[T]) TryPop(indexTyped ArrayPlusIndex) (T, bool) {
//
//	index := int(indexTyped)
//
//	elem := stack.arr[index]
//
//	stack.gaps.Push(indexTyped)
//
//	return elem, true
//}

func (stack *ArrayPlus[T]) Peek(indexTyped ArrayPlusIndex) T {
	index := int(indexTyped)

	elem := stack.arr[index]

	return elem
}

func (stack *ArrayPlus[T]) Push(elem T) ArrayPlusIndex {
	index, ok := stack.gaps.TryPop()
	if ok {
		stack.arr[index] = elem
		return index
	}
	stack.lock.Lock()
	if len(stack.arr) >= stack.maxLen {
		stack.lock.Unlock()
		index = stack.gaps.Pop()
		stack.arr[index] = elem
		return index
	}
	oldLen := len(stack.arr)
	addLength := min(stack.maxLen-len(stack.arr), _CHUNKSIZE_)
	newChuck := make([]T, addLength)
	stack.arr = append(stack.arr, newChuck...)
	for i := range addLength - 1 {
		stack.gaps.Push(ArrayPlusIndex(i + (oldLen + 1)))
	}
	stack.arr[oldLen] = elem
	stack.lock.Unlock()
	return ArrayPlusIndex(oldLen)
}

func (stack *ArrayPlus[T]) TryPush(elem T) (ArrayPlusIndex, bool) {
	index, ok := stack.gaps.TryPop()
	if ok {
		stack.arr[index] = elem
		return index, true
	}
	stack.lock.Lock()
	if len(stack.arr) >= stack.maxLen {
		stack.lock.Unlock()
		return 0, false
	}
	oldLen := len(stack.arr)
	addLength := min(stack.maxLen-len(stack.arr), _CHUNKSIZE_)
	newChuck := make([]T, addLength)
	stack.arr = append(stack.arr, newChuck...)
	for i := range addLength - 1 {
		stack.gaps.Push(ArrayPlusIndex(i + (oldLen + 1)))
	}
	stack.arr[oldLen] = elem
	stack.lock.Unlock()
	return ArrayPlusIndex(oldLen), true
}
