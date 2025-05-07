package Sync

import "sync"

/*
NewType used to represent the index of an element in the array.

# Invariants:
  - Should only be gained as a handle after a Push operation.
  - Should itself and all of its copies be discarded after a Pop operation.
*/
type ArrayPlusIndex int

type ArrayPlus[T any] struct {
	lock   sync.Mutex
	gaps   Stack[ArrayPlusIndex]
	arr    []T
	maxLen int
}

const _CHUNKSIZE_ = 100

/*
NewArrayPlus
creates a new ArrayPlus with the given size.
*/
func NewArrayPlus[T any](maxSize int) *ArrayPlus[T] {
	lenght := min(maxSize/3, _CHUNKSIZE_)
	newArr := ArrayPlus[T]{
		gaps:   NewStack[ArrayPlusIndex](maxSize),
		arr:    make([]T, lenght),
		maxLen: maxSize,
	}
	for i := range lenght {
		newArr.gaps.Push(ArrayPlusIndex(i))
	}
	return &newArr
}

/*
Pop
retrieves the element at the given index.
*/
func (stack *ArrayPlus[T]) Pop(indexTyped ArrayPlusIndex) T {
	index := int(indexTyped)

	elem := stack.arr[index]

	stack.gaps.Push(indexTyped)

	return elem
}

/*
// with current invariants, this function is not needed
func (stack *ArrayPlus[T]) TryPop(indexTyped ArrayPlusIndex) (T, bool) {

	index := int(indexTyped)

	elem := stack.arr[index]

	stack.gaps.Push(indexTyped)

	return elem, true
}
*/

/*
Peek
returns the element at the given index without removing it.
*/
func (stack *ArrayPlus[T]) Peek(indexTyped ArrayPlusIndex) T {
	index := int(indexTyped)

	elem := stack.arr[index]

	return elem
}

/*
Push
adds an element to the array and returns the index of the element.
If the array is full, it will either grow the array if maxLen is not reached,
or it will block until space is available.
*/
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

/*
TryPush
attempts to add an element to the array and returns the index of the element.
If the array is full, it will either grow the array if maxLen is not reached,
or it will return false without blocking.
*/
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
