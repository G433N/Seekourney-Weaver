package concurrencyUtils_test

import (
	"seekourney/utils/concurrencyUtils"
	"sync"
	"testing"
)

func TestArrayPlusAdvanced(t *testing.T) {
	arrPlus := concurrencyUtils.NewArrayPlus[int](10)
	errSem := concurrencyUtils.NewSemaphore()

	wg := sync.WaitGroup{}
	wg.Add(200)
	for i := range 200 {
		go func() {
			defer wg.Done()
			number := i
			for range 200 {
				index := arrPlus.Push(number)
				if number != arrPlus.Peek(index) {
					errSem.Signal()
				}
				newNumber := arrPlus.Pop(index)
				if number != newNumber {
					errSem.Signal()
				}
				number = newNumber + 200
			}
		}()
	}
	wg.Wait()
	for errSem.TryWait() {
		t.Error("Error in ArrayPlus")
	}

}

func TestArrayPlusMaxSizeBlocking(t *testing.T) {
	arrPlus := concurrencyUtils.NewArrayPlus[int](17)
	errSem := concurrencyUtils.NewSemaphore()

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
