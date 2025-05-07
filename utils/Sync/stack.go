package Sync

import (
	"log"
	"sync"
)

// currently unused
//const _STACKMAXLEN_ = 1_000_000

/*
Stack is a thread-safe generic stack implementation.
Can Pop, Push and has non-blocking versions too.
*/
type Stack[T any] struct {
	// lock is used to protect access to the stack
	lock sync.Mutex
	// stack is the slice that holds the elements
	stack []T
	// lenSem is a semaphore that tracks the number of elements in the stack
	// used to block the Pop operation when the stack is empty
	lenSem Semaphore
}

/*
NewStack
creates a new Stack with the optional maximum size.
*/
func NewStack[T any](optionalMaxSize ...int) Stack[T] {
	if len(optionalMaxSize) == 1 {
		return Stack[T]{
			stack:  []T{},
			lenSem: *NewSemaphore(0, optionalMaxSize[0]),
		}
	}
	if len(optionalMaxSize) > 1 {
		log.Fatalln("NewStack", "only one size argument is allowed")
	}
	return Stack[T]{
		stack:  []T{},
		lenSem: *NewSemaphore(0),
	}
}

/*
IsEmpty
returns true if the stack is empty.
*/
func (stack *Stack[T]) IsEmpty() bool {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	return len(stack.stack) == 0
}

/*
Pop
retrieves the top element from the stack.
If empty, blocks until its not
*/
func (stack *Stack[T]) Pop() T {
	stack.lenSem.Wait()
	stack.lock.Lock()
	URL := stack.stack[len(stack.stack)-1]
	stack.stack = stack.stack[:len(stack.stack)-1]
	stack.lock.Unlock()
	return URL
}

/*
TryPop
retrieves the top element from the stack.
If empty, returns false without blocking.
*/
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

/*
Push
adds an element to the stack.
If the stack is full, it blocks until space is available.
*/
func (stack *Stack[T]) Push(elem T) {
	stack.lenSem.Signal()
	stack.lock.Lock()
	stack.stack = append(stack.stack, elem)
	stack.lock.Unlock()
}

/*
TryPush
attempts to add an element to the stack.
If the stack is full, it returns false without blocking.
*/
func (stack *Stack[T]) TryPush(elem T) bool {
	stack.lock.Lock()
	if !stack.lenSem.TrySignal() {
		return false
	}

	stack.stack = append(stack.stack, elem)
	stack.lock.Unlock()
	return true
}
