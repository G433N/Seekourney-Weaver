package concurrencyUtils_test

import (
	"seekourney/utils/concurrencyUtils"
	"sync"
	"testing"
	"time"
)

type chanSem chan bool

func (s chanSem) Signal() {
	s <- true
}
func (s chanSem) Wait() {
	<-s
}
func (s chanSem) TryWait() bool {
	select {
	case <-s:
		return true
	default:
		return false
	}
}
func (s chanSem) TrySignal() bool {
	select {
	case s <- true:
		return true
	default:
		return false
	}
}

type semaphoreInterface interface {
	Signal()
	Wait()
	TryWait() bool
	TrySignal() bool
}

func benchSemaphore[C semaphoreInterface](sem C,
	iterations int,
	threads int) {
	wg := sync.WaitGroup{}
	wg.Add(2 * threads)
	for range threads {
		go func() {
			defer wg.Done()
			for range iterations {
				sem.Signal()
			}
		}()
		go func() {
			defer wg.Done()
			for range iterations {
				sem.Wait()
			}
		}()

	}
	wg.Wait()
}

func BenchmarkSemaphoreChanelBasic10(b *testing.B) {
	for b.Loop() {
		benchSemaphore(chanSem(make(chan bool, 1000)), 10, 10)
	}

}

func BenchmarkSemaphoreNewBasic10(b *testing.B) {
	for b.Loop() {
		benchSemaphore(concurrencyUtils.NewSemaphore(), 10, 10)
	}
}

func BenchmarkSemaphoreChanelBasic10000_100(b *testing.B) {
	for b.Loop() {
		benchSemaphore(chanSem(make(chan bool, 10000000)), 10000, 100)
	}

}

func BenchmarkSemaphoreNewBasic10000_100(b *testing.B) {
	for b.Loop() {
		benchSemaphore(concurrencyUtils.NewSemaphore(), 10000, 100)
	}
}

func BenchmarkSemaphoreNewBoundedBasic10000_100(b *testing.B) {
	for b.Loop() {
		benchSemaphore(concurrencyUtils.NewSemaphore(0, 10000000), 10000, 100)
	}
}

func BenchmarkSemaphoreChanelSmallBuffer1000_4000(b *testing.B) {
	for b.Loop() {
		benchSemaphore(chanSem(make(chan bool, 20)), 1000, 4000)
	}

}

func BenchmarkSemaphoreNewSmallbuffer1000_4000(b *testing.B) {
	for b.Loop() {
		benchSemaphore(concurrencyUtils.NewSemaphore(0, 20), 1000, 4000)
	}
}

func BenchmarkSemaphoreNewNobuffer1000_4000(b *testing.B) {
	for b.Loop() {
		benchSemaphore(concurrencyUtils.NewSemaphore(), 1000, 4000)
	}
}

func benchSemaphore2[C semaphoreInterface](
	sem C,
	iterations int,
	threads int) {
	wg := sync.WaitGroup{}
	wg.Add(threads)
	for range threads {
		go func() {
			defer wg.Done()
			for range iterations {
				sem.Signal()
				time.Sleep(10 * time.Microsecond)
				sem.Wait()
			}
		}()

	}
	wg.Wait()
}

func BenchmarkSemaphoreChanelSmallBuffer1000_400(b *testing.B) {
	for b.Loop() {
		benchSemaphore2(chanSem(make(chan bool, 20)), 1000, 400)
	}

}

func BenchmarkSemaphoreNewSmallbuffer1000_400(b *testing.B) {
	for b.Loop() {
		benchSemaphore2(concurrencyUtils.NewSemaphore(0, 20), 1000, 400)
	}
}

func BenchmarkSemaphoreNewNobuffer1000_400(b *testing.B) {
	for b.Loop() {
		benchSemaphore2(concurrencyUtils.NewSemaphore(), 1000, 400)
	}
}
