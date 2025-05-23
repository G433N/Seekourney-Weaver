package concurrencyUtils

import "sync"

/*
ArrayPlusIndex represents the index of an element in the array.

# Invariants:
  - Should only be gained as a handle after a Push operation.
  - Should itself and all of its copies be discarded after a Pop operation.
*/
type ArrayPlusIndex int

/*
ArrayPlus is a thread-safe wraper for a slice.
Meant to be used as kind of a variation on a map
where instead of a key, an index is used as a hadle
*/
type ArrayPlus[T any] struct {
	lock sync.Mutex
	// Stack to keep track of free indexes
	gaps Stack[ArrayPlusIndex]
	// The wrapped slice
	arr []T
	// The maximum length of the slice
	maxLen int
}

const _CHUNKSIZE_ = 100

/*
NewArrayPlus
creates a new ArrayPlus with the given size.
*/
func NewArrayPlus[T any](maxSize int) *ArrayPlus[T] {

	lenght := min(maxSize, _CHUNKSIZE_)

	arrPlus := ArrayPlus[T]{
		gaps:   *NewStack[ArrayPlusIndex](maxSize),
		arr:    make([]T, lenght),
		maxLen: maxSize,
	}

	for i := range lenght {
		arrPlus.gaps.Push(ArrayPlusIndex(i))
	}

	return &arrPlus
}

/*
Pop
retrieves the element at the given index,
which means also removing it from the array.
*/
func (arrPlus *ArrayPlus[T]) Pop(indexTyped ArrayPlusIndex) T {

	index := int(indexTyped)

	elem := arrPlus.arr[index]

	arrPlus.gaps.Push(indexTyped)

	return elem
}

/*
Peek
returns the element at the given index without removing it.
*/
func (arrPlus *ArrayPlus[T]) Peek(indexTyped ArrayPlusIndex) T {

	index := int(indexTyped)

	elem := arrPlus.arr[index]

	return elem
}

/*
GrowSlice
grows the slice by a fixed amount.
This is used when the array is full and a new element is added.
Assumes that the lock is already held.
*/
func (arrPlus *ArrayPlus[T]) growSlice(chunkLen int) ArrayPlusIndex {

	oldLen := len(arrPlus.arr)
	addLength := min(arrPlus.maxLen-oldLen, chunkLen)
	newChuck := make([]T, addLength)
	arrPlus.arr = append(arrPlus.arr, newChuck...)

	for i := range addLength - 1 {
		arrPlus.gaps.Push(ArrayPlusIndex(i + (oldLen + 1)))
	}
	return ArrayPlusIndex(oldLen)
}

/*
Push
adds an element to the array and returns the index of the element.
If the array is full, it will either grow the array if maxLen is not reached,
or it will block until space is available.
*/
func (arrPlus *ArrayPlus[T]) Push(elem T) ArrayPlusIndex {

	index, ok := arrPlus.gaps.TryPop()

	if ok {
		arrPlus.arr[index] = elem
		return index
	}

	arrPlus.lock.Lock()

	if len(arrPlus.arr) >= arrPlus.maxLen {

		arrPlus.lock.Unlock()
		// waits for a new gap to open
		index = arrPlus.gaps.Pop()

		arrPlus.arr[index] = elem

		return index
	}

	newIndex := arrPlus.growSlice(_CHUNKSIZE_)

	arrPlus.lock.Unlock()

	arrPlus.arr[newIndex] = elem

	return newIndex
}

/*
TryPush
attempts to add an element to the array and returns the index of the element.
If the array is full, it will either grow the array if maxLen is not reached,
or it will return false without blocking.
*/
func (arrPlus *ArrayPlus[T]) TryPush(elem T) (ArrayPlusIndex, bool) {

	index, ok := arrPlus.gaps.TryPop()

	if ok {
		arrPlus.arr[index] = elem
		return index, true
	}

	arrPlus.lock.Lock()
	if len(arrPlus.arr) >= arrPlus.maxLen {
		arrPlus.lock.Unlock()
		return 0, false
	}

	newIndex := arrPlus.growSlice(_CHUNKSIZE_)

	arrPlus.lock.Unlock()

	arrPlus.arr[newIndex] = elem

	return newIndex, true
}
