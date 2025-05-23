package concurrencyUtils

import (
	"log"
	"sync"
)

/*
Stack is a thread-safe generic stack implementation.
Can Pop, Push and has non-blocking versions too.
*/
type Stack[T any] struct {
	// lock is used to protect access to the stack
	lock sync.Mutex
	// stack is the slice that holds the elements
	stack []T
	// filledSem tracks the number of elements in the stack
	// used to block the Pop operation when the stack is empty
	filledSem Semaphore
	// emptySem tracks how many elements can be added before the stack is full
	// and is only used when a max size is given
	// used to block the Push operation when the stack is full
	emptySem Semaphore

	/*
		Pop
		retrieves the top element from the stack.
		If empty, blocks until its not
	*/
	Pop func() T

	/*
	   TryPop
	   retrieves the top element from the stack.
	   If empty, returns false without blocking.
	*/
	TryPop func() (T, bool)

	/*
		Push
		adds an element to the stack.
		If a max size is given and the stack is full,
		it blocks until space is available.
	*/
	Push func(T)

	/*
	   TryPush
	   attempts to add an element to the stack.
	   If a max size is given and stack is full,
	   it returns false without blocking.
	*/
	TryPush func(T) bool
}

/*
NewStack
creates a new Stack with the optional maximum size.
*/
func NewStack[T any](optionalMaxSize ...int) *Stack[T] {
	if len(optionalMaxSize) == 1 {
		newStack := Stack[T]{
			stack:     []T{},
			filledSem: *NewSemaphore(0),
			emptySem:  *NewSemaphore(optionalMaxSize[0]),
		}
		newStack.Pop = newStack.boundedPop
		newStack.TryPop = newStack.boundedTryPop
		newStack.Push = newStack.boundedPush
		newStack.TryPush = newStack.boundedTryPush

		return &newStack
	}
	if len(optionalMaxSize) > 1 {
		log.Fatalln("NewStack", "only one size argument is allowed")
	}

	newStack :=
		Stack[T]{
			stack:     []T{},
			filledSem: *NewSemaphore(0),
		}
	newStack.Pop = newStack.defaultPop
	newStack.TryPop = newStack.defaultTryPop
	newStack.Push = newStack.defaultPush
	newStack.TryPush = newStack.defaultTryPush

	return &newStack
}

/*
IsEmpty
returns true if the stack is empty.
*/
func (st *Stack[T]) IsEmpty() bool {
	st.lock.Lock()
	defer st.lock.Unlock()
	return len(st.stack) == 0
}

/*
defaultPop
retrieves the top element from the stack.
If empty, blocks until its not.
Is only used when a max size is not given.
*/
func (st *Stack[T]) defaultPop() T {
	st.filledSem.Wait()
	st.lock.Lock()
	elem := st.stack[len(st.stack)-1]
	st.stack = st.stack[:len(st.stack)-1]
	st.lock.Unlock()
	return elem
}

/*
defaultTryPop
tries to retrieve the top element from the stack.
If empty, returns false without blocking.
Is only used when a max size is not given.
*/
func (st *Stack[T]) defaultTryPop() (T, bool) {
	if !st.filledSem.TryWait() {
		var zeroValue T
		return zeroValue, false
	}
	st.lock.Lock()
	elem := st.stack[len(st.stack)-1]
	st.stack = st.stack[:len(st.stack)-1]
	st.lock.Unlock()
	return elem, true
}

/*
defaultPush
adds an element to the stack.
Is only used when a max size is not given.
*/
func (st *Stack[T]) defaultPush(elem T) {
	st.lock.Lock()
	st.stack = append(st.stack, elem)
	st.filledSem.Signal()
	st.lock.Unlock()
}

/*
defaultTryPush
tries to add an element to the stack.
Is only used when a max size is not given.
*/
func (stack *Stack[T]) defaultTryPush(elem T) bool {
	stack.lock.Lock()
	stack.stack = append(stack.stack, elem)
	stack.filledSem.Signal()
	stack.lock.Unlock()
	return true
}

/*
boundedPop
retrieves the top element from the stack.
If empty, blocks until its not.
Is only used when a max size is given.
*/
func (st *Stack[T]) boundedPop() T {

	st.filledSem.Wait()
	st.lock.Lock()
	elem := st.stack[len(st.stack)-1]
	st.stack = st.stack[:len(st.stack)-1]
	st.emptySem.Signal()
	st.lock.Unlock()
	return elem
}

/*
boundedTryPop
tries to retrieve the top element from the stack.
If empty, returns false without blocking.
Is only used when a max size is given.
*/
func (st *Stack[T]) boundedTryPop() (T, bool) {
	if !st.filledSem.TryWait() {
		var zeroValue T
		return zeroValue, false
	}
	st.lock.Lock()
	elem := st.stack[len(st.stack)-1]
	st.stack = st.stack[:len(st.stack)-1]
	st.emptySem.Signal()
	st.lock.Unlock()
	return elem, true
}

/*
boundedPush
adds an element to the stack.
If the stack is full it blocks until space is available.
Is only used when a max size is given.
*/
func (st *Stack[T]) boundedPush(elem T) {
	st.emptySem.Wait()
	st.lock.Lock()
	st.stack = append(st.stack, elem)
	st.filledSem.Signal()
	st.lock.Unlock()
}

/*
boundedTryPush
tries to add an element to the stack.
If the stack is full it returns false without blocking
Is only used when a max size is given.
*/
func (stack *Stack[T]) boundedTryPush(elem T) bool {
	if !stack.emptySem.TryWait() {
		return false
	}
	stack.lock.Lock()
	stack.stack = append(stack.stack, elem)
	stack.filledSem.Signal()
	stack.lock.Unlock()
	return true
}
