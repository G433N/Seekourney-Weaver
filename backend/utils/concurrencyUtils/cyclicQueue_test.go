package concurrencyUtils_test

import (
	"seekourney/utils/concurrencyUtils"
	"sync"
	"testing"
)

func TestCyclicQueueBasicTryPopPush(t *testing.T) {
	queue := concurrencyUtils.NewCyclicQueue[int](5)

	for i := range 5 {
		if !queue.TryPush(i) {
			t.Error("TryPush failed")
		}
	}

	if ok := queue.TryPush(10); ok {
		t.Error("TryPush should fail on full queue")
	}

	for i := range 5 {
		val, ok := queue.TryPop()
		if !ok || val != i {
			t.Error("TryPop failed")
		}
	}

	if _, ok := queue.TryPop(); ok {
		t.Error("TryPop should fail on empty queue")
	}

	for i := range 4 {
		if !queue.TryPush(i) {
			t.Error("TryPush failed")
		}
	}
	if !queue.TryPush(10) {
		t.Error("TryPush failed")
	}
	if ok := queue.TryPush(10); ok {
		t.Error("TryPush should fail on full queue")
	}

	for i := range 4 {
		val, ok := queue.TryPop()
		if !ok || val != i {
			t.Error("TryPop failed")
		}
	}
	if val, ok := queue.TryPop(); !ok || val != 10 {
		t.Error("TryPop failed")
	}

	if _, ok := queue.TryPop(); ok {
		t.Error("TryPop should fail on empty queue")
	}
}

func TestCyclicQueueBasicPushPop(t *testing.T) {
	queue := concurrencyUtils.NewCyclicQueue[int](5)

	for i := range 5 {
		queue.Push(i)
	}

	if ok := queue.TryPush(10); ok {
		t.Error("TryPush should fail on full queue")
	}

	for i := range 5 {
		val, ok := queue.TryPop()
		if !ok || val != i {
			t.Error("TryPop failed")
		}
	}

	if _, ok := queue.TryPop(); ok {
		t.Error("TryPop should fail on empty queue")
	}

	for i := range 4 {
		queue.Push(i)
	}
	if !queue.TryPush(10) {
		t.Error("TryPush failed")
	}
	if ok := queue.TryPush(10); ok {
		t.Error("TryPush should fail on full queue")
	}

	for i := range 4 {
		val, ok := queue.TryPop()
		if !ok || val != i {
			t.Error("TryPop failed")
		}
	}
	val, ok := queue.TryPop()
	if !ok || val != 10 {
		t.Error("TryPop failed")
	}

	if _, ok := queue.TryPop(); ok {
		t.Error("TryPop should fail on empty queue")
	}
}

func TestCyclicQueueConcurrent(t *testing.T) {
	queue := concurrencyUtils.NewCyclicQueue[int](5)

	var wg sync.WaitGroup

	// Push goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 10 {
			queue.Push(i)
		}
	}()

	// Pop goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 10 {
			val := queue.Pop()
			if val != i {
				t.Error("TryPop failed")
			}
		}
	}()

	wg.Wait()
}

func TestCyclicQueueConcurrency2(t *testing.T) {
	queueC := concurrencyUtils.NewCyclicQueue[int](10)

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := range 10 {
			queueC.Push(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := range 10 {
			queueC.Push(i)

		}
	}()
	go func() {
		defer wg.Done()
		for range 20 {
			queueC.Pop()
		}
	}()
	wg.Wait()

}

func TestCyclicQueueAdvancedConcurrency(t *testing.T) {
	queueA := concurrencyUtils.NewCyclicQueue[int](7)
	queueB := concurrencyUtils.NewCyclicQueue[int](4)
	queueC := concurrencyUtils.NewCyclicQueue[int](11)

	errSem := concurrencyUtils.NewSemaphore(0)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := range 10 {
			queueA.Push(i)
			queueC.Push(i)
		}
		queueA.Push(-1)

	}()
	go func() {
		defer wg.Done()
		for i := range 10 {
			queueB.Push(i)
			queueC.Push(i)

		}
		queueB.Push(-1)

	}()
	go func() {
		defer wg.Done()
		a := queueA.Pop()
		b := queueB.Pop()
		for range 20 {
			c := queueC.Pop()
			switch c {
			case a:
				a = queueA.Pop()
			case b:
				b = queueB.Pop()
			default:
				errSem.Signal()
				return
			}
		}
		if _, ok := queueA.TryPop(); ok || a != -1 {
			errSem.Signal()
		}
		if _, ok := queueB.TryPop(); ok || b != -1 {
			errSem.Signal()
		}
	}()
	wg.Wait()
	if errSem.TryWait() {
		t.Error("queueC didn't match queueA and queueB")
	}
}
