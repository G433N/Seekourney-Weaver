package concurrencyUtils_test

import (
	"seekourney/utils/concurrencyUtils"
	"testing"
)

func TestStack(t *testing.T) {
	stack := concurrencyUtils.NewStack[int]()
	for i := range 20 {
		if ok := stack.TryPush(i); !ok {
			t.Errorf("Push failed for value %d", i)
		}
	}

	for i := range 20 {
		if ok := stack.TryPush(i); !ok {
			t.Errorf("Push failed for value %d", i)
		}
	}

	for i := range 20 {
		val, ok := stack.TryPop()
		if !ok || val != 19-i {
			t.Errorf("Pop failed, expected %d, got %d", 19-i, val)
		}
	}

	for i := range 20 {
		if ok := stack.TryPush(i + 20); !ok {
			t.Errorf("Push failed for value %d", i)
		}
	}

	for i := range 40 {
		val, ok := stack.TryPop()
		if !ok || val != 39-i {
			t.Errorf("Pop failed, expected %d, got %d", 39-i, val)
		}
	}
}
func TestStackBlocking(t *testing.T) {
	stack1 := concurrencyUtils.NewStack[int]()
	errSem1 := concurrencyUtils.NewSemaphore()

	for range 20 {
		go func() {
			stack1.Pop()
			errSem1.Signal()
		}()
	}

	stack2 := concurrencyUtils.NewStack[int](1)
	stack2.Push(1)
	errSem2 := concurrencyUtils.NewSemaphore()
	for range 20 {
		go func() {
			stack2.Push(1)
			errSem2.Signal()
		}()
	}
	if errSem1.TryWait() {
		t.Error("Empty Pop blocking failed")
	}
	if errSem2.TryWait() {
		t.Error("Full Push blocking failed")
	}
}
