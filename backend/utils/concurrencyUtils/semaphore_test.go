package concurrencyUtils_test

import (
	"seekourney/utils/concurrencyUtils"
	"sync"
	"testing"
)

func TestTryWait(test *testing.T) {
	semaphore := concurrencyUtils.NewSemaphore()
	if semaphore.TryWait() {
		test.Error("TryWait should return false on an initial empty semaphore")
	}

	semaphore.Signal()

	if !semaphore.TryWait() {
		test.Error("TryWait should return true on a non empty semaphore")
	}
	if semaphore.TryWait() {
		test.Error("TryWait should return false on an empty semaphore")
	}

	semaphore.Signal()
	semaphore.Signal()

	if !semaphore.TryWait() {
		test.Error("TryWait should return true on a non empty semaphore")
	}
	if !semaphore.TryWait() {
		test.Error("TryWait should return true on a non empty semaphore")
	}
	if semaphore.TryWait() {
		test.Error("TryWait should return false on an empty semaphore")
	}
}

const _RENDEZVOULOOPS_ = 5000

func rendezvousHelper(
	result *[]int, ID int, waitGroup *sync.WaitGroup,
	mutex *sync.Mutex,
	ownSem *concurrencyUtils.Semaphore, otherSem *concurrencyUtils.Semaphore,
) {
	for range _RENDEZVOULOOPS_ {
		mutex.Lock()
		*result = append(*result, ID)
		mutex.Unlock()
		ownSem.Signal()
		otherSem.Wait()
	}
	ownSem.Signal()
	otherSem.Wait()
	ownSem.Signal()
	otherSem.Wait()
	waitGroup.Done()
}

func testRendezvousHelper(test *testing.T, ID int, waitGroup *sync.WaitGroup) {
	sem1 := concurrencyUtils.NewSemaphore()
	sem2 := concurrencyUtils.NewSemaphore()
	result := []int{}
	var mutex sync.Mutex
	var childWaitGroup sync.WaitGroup
	childWaitGroup.Add(2)
	go rendezvousHelper(&result, 1, &childWaitGroup, &mutex, sem1, sem2)
	go rendezvousHelper(&result, 2, &childWaitGroup, &mutex, sem2, sem1)
	childWaitGroup.Wait()

	mem1 := 0
	mem2 := 0
	i := 0
	for _, x := range result {
		i++
		if mem1 == mem2 && mem2 == x {
			test.Error(
				"thread:", ID, "error\n",
				"Rendezvous failed", i, "iterations in,",
				"with three", x, "in a row",
			)
		}
		mem1 = mem2
		mem2 = x
	}
	if i != _RENDEZVOULOOPS_*2 {
		test.Error("thread:", ID, "error\n",
			"should loop", _RENDEZVOULOOPS_*2, "times,",
			"but did loop", i, "times",
		)
	}

	if sem1.TryWait() {
		test.Error("thread:", ID, "error\n", "semaphore sem1 should be empty")
	}
	if sem2.TryWait() {
		test.Error("thread:", ID, "error\n", "semaphore sem2 should be empty")
	}
	waitGroup.Done()
}

func TestRendezvous(test *testing.T) {
	if testing.Short() {
		test.Skip("skipping long test")
	}
	var waitGroup sync.WaitGroup
	children := 1000
	waitGroup.Add(children)

	for x := range children {
		go testRendezvousHelper(test, x, &waitGroup)
	}
	waitGroup.Wait()
}

func TestWaitBlock(test *testing.T) {
	children := 1000
	sem := concurrencyUtils.NewSemaphore()
	for x := range children {
		go func() {
			sem.Wait()
			test.Error("thread:", x, "error\n",
				"should be blocked and never get here")
		}()
	}

}

func TestLargeNumberGoroutine(test *testing.T) {
	sem1 := concurrencyUtils.NewSemaphore(1)
	sem2 := concurrencyUtils.NewSemaphore()
	var waitGroup sync.WaitGroup
	waitGroup.Add(200)
	for range 100 {
		go func(w *sync.WaitGroup) {
			sem1.Wait()
			sem2.Signal()
			w.Done()
		}(&waitGroup)
	}
	for range 100 {
		go func(w *sync.WaitGroup) {
			sem2.Wait()
			sem1.Signal()
			w.Done()
		}(&waitGroup)
	}
	waitGroup.Wait()

	if !sem1.TryWait() {
		test.Error("sem1 should have a signal")
	}
	if sem1.TryWait() {
		test.Error("sem1 should be empty")
	}
	if sem2.TryWait() {
		test.Error("sem2 should be empty")
	}
}

func TestBoundedSemaphores(test *testing.T) {
	sem1 := concurrencyUtils.NewSemaphore()
	sem2 := concurrencyUtils.NewSemaphore(0)
	sem3 := concurrencyUtils.NewSemaphore(1)
	sem4 := concurrencyUtils.NewSemaphore(1, 1)
	sem5 := concurrencyUtils.NewSemaphore(0, 2)

	sem1.Signal()
	if !sem1.TryWait() {
		test.Error("sem1 should have a signal")
	}
	if sem1.TryWait() {
		test.Error("sem1 should be empty")
	}
	if sem2.TryWait() {
		test.Error("sem2 should be empty")
	}
	if !sem3.TryWait() {
		test.Error("sem3 should have a signal")
	}
	if sem3.TryWait() {
		test.Error("sem3 should be empty")
	}
	if sem4.TrySignal() {
		test.Error("sem4 should be full")
	}
	if !sem4.TryWait() {
		test.Error("sem4 should have a signal")
	}
	if sem4.TryWait() {
		test.Error("sem4 should be empty")
	}
	if !sem5.TrySignal() {
		test.Error("sem5 should have space for a signal")
	}
	if !sem5.TrySignal() {
		test.Error("sem5 should have space for a signal")
	}
	if sem5.TrySignal() {
		test.Error("sem5 should be full")
	}
	if !sem5.TryWait() {
		test.Error("sem5 should have a signal")
	}
	if !sem5.TrySignal() {
		test.Error("sem5 should have space for a signal")
	}
	if sem5.TrySignal() {
		test.Error("sem5 should be full")
	}
	if !sem5.TryWait() {
		test.Error("sem5 should have a signal")
	}
	if !sem5.TryWait() {
		test.Error("sem5 should have a signal")
	}
	if sem5.TryWait() {
		test.Error("sem5 should be empty")
	}

}

func TestSignalBoundedBlock(test *testing.T) {
	children := 1000
	sem := concurrencyUtils.NewSemaphore(1, 1)
	for x := range children {
		go func() {
			sem.Signal()
			test.Error("thread:", x, "error\n",
				"should be blocked and never get here")
		}()
	}

}

func TestBoundedRendezvous(test *testing.T) {
	var waitGroup sync.WaitGroup
	children := 1000
	waitGroup.Add(children)

	for x := range children {
		go testBoundedRendezvousHelper(test, x, &waitGroup)
	}
	waitGroup.Wait()
}

func testBoundedRendezvousHelper(test *testing.T, ID int, waitGroup *sync.WaitGroup) {
	sem1 := concurrencyUtils.NewSemaphore(0, 1)
	sem2 := concurrencyUtils.NewSemaphore(0, 1)
	result := []int{}
	var mutex sync.Mutex
	var childWaitGroup sync.WaitGroup
	childWaitGroup.Add(2)
	go rendezvousHelper(&result, 1, &childWaitGroup, &mutex, sem1, sem2)
	go rendezvousHelper(&result, 2, &childWaitGroup, &mutex, sem2, sem1)
	childWaitGroup.Wait()

	mem1 := 0
	mem2 := 0
	i := 0
	for _, x := range result {
		i++
		if mem1 == mem2 && mem2 == x {
			test.Error(
				"thread:", ID, "error\n",
				"Rendezvous failed", i, "iterations in,",
				"with three", x, "in a row",
			)
		}
		mem1 = mem2
		mem2 = x
	}
	if i != _RENDEZVOULOOPS_*2 {
		test.Error("thread:", ID, "error\n",
			"should loop", _RENDEZVOULOOPS_*2, "times,",
			"but did loop", i, "times",
		)
	}

	if sem1.TryWait() {
		test.Error("thread:", ID, "error\n", "semaphore sem1 should be empty")
	}
	if sem2.TryWait() {
		test.Error("thread:", ID, "error\n", "semaphore sem2 should be empty")
	}
	waitGroup.Done()
}

func TestBoundedLargeNumberGoroutine1(test *testing.T) {
	sem1 := concurrencyUtils.NewSemaphore(1, 100)
	sem2 := concurrencyUtils.NewSemaphore(0, 100)
	var waitGroup sync.WaitGroup
	waitGroup.Add(400)
	for range 200 {
		go func(w *sync.WaitGroup) {
			sem1.Wait()
			sem2.Signal()
			w.Done()
		}(&waitGroup)
	}
	for range 200 {
		go func(w *sync.WaitGroup) {
			sem2.Wait()
			sem1.Signal()
			w.Done()
		}(&waitGroup)
	}
	waitGroup.Wait()

	if !sem1.TryWait() {
		test.Error("sem1 should have a signal")
	}
	if sem1.TryWait() {
		test.Error("sem1 should be empty")
	}
	if sem2.TryWait() {
		test.Error("sem2 should be empty")
	}
}

func TestBoundedLargeNumberGoroutine2(test *testing.T) {
	sem1 := concurrencyUtils.NewSemaphore(1, 10)
	var waitGroup sync.WaitGroup
	waitGroup.Add(400)
	for range 200 {
		go func(w *sync.WaitGroup) {
			sem1.Signal()
			w.Done()
		}(&waitGroup)
	}
	for range 200 {
		go func(w *sync.WaitGroup) {
			sem1.Wait()
			w.Done()
		}(&waitGroup)
	}
	waitGroup.Wait()

	if !sem1.TryWait() {
		test.Error("sem1 should have a signal")
	}
	if sem1.TryWait() {
		test.Error("sem1 should be empty")
	}
}
