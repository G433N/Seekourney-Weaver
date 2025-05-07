package Sync_test

import (
	"seekourney/utils/Sync"
	"sync"
	"testing"
)

func TestArrayPlusAdvanced(t *testing.T) {
	arrPlus := Sync.NewArrayPlus[int](10)
	errSem := Sync.NewSemaphore()

	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := range 20 {
		go func() {
			defer wg.Done()
			number := i
			for range 20 {
				index := arrPlus.Push(number)
				if number != arrPlus.Peek(index) {
					errSem.Signal()
				}
				newNumber := arrPlus.Pop(index)
				if number != newNumber {
					errSem.Signal()
				}
				number = newNumber + 20
			}
		}()
	}
	wg.Wait()
	for errSem.TryWait() {
		t.Error("Error in ArrayPlus")
	}

}

func TestArrayPlusMaxSizeBlocking(t *testing.T) {
	arrPlus := Sync.NewArrayPlus[int](17)
	errSem := Sync.NewSemaphore()

	for range 17 {
		arrPlus.Push(0)
	}
	for range 20 {
		go func() {
			arrPlus.Push(0)
			errSem.Signal()
		}()
	}
	if errSem.TryWait() {
		t.Error("MaxSize blocking failed")
	}
}
