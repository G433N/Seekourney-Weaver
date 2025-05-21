package concurrencyUtils

import (
	"log"
	"math"
	"sync"
)

/*
A very basic semaphore implementation
*/
type Semaphore struct {
	//internal value to keep track of amount signals
	value int
	// mutex used for the internal value
	syncLock sync.Mutex

	/*
		Signal
		Increments the semaphores value.
		If it was 0 it will unblock any blocked Wait thread.
		If bounded, it will block if the value is currently the maximum value.
	*/
	Signal func()

	/*
		Wait
		Decrement the semaphore’s value.
		It will block if the value is currently 0.
		If bounded, it will unblock any blocked Signal thread.
	*/
	Wait func()

	/*
		TryWait
		Tries to decrement the semaphore’s value and return sucess/true.
		If the value is 0 it will return false.
	*/
	TryWait func() bool

	/*
		TrySignal
		Tries to increment the semaphore’s value and return sucess/true.
		If the value is currently the maximum value it will return false.
	*/
	TrySignal func() bool

	// the maximum value of the semaphore
	maxValue int

	// mutex used for the signaling waitgroup
	// only used for the bounded semaphore
	signalingLock  sync.Mutex
	signalingGroup sync.WaitGroup

	// mutex used for the waiting waitgroup
	waitingLock  sync.Mutex
	waitingGroup sync.WaitGroup
}

/*
Wait
Decrement the semaphore’s value; block if the value is currently 0.
*/
func (semaphore *Semaphore) defaultWait() {

	// so that multiple threads don't use the waitgroup at the same time
	semaphore.waitingLock.Lock()
	defer semaphore.waitingLock.Unlock()

	semaphore.waitingGroup.Wait()

	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()

	if semaphore.value == 1 {
		semaphore.waitingGroup.Add(1)
	}

	semaphore.value--

}

/*
Wait
Decrement the semaphore’s value; block if the value is currently 0.
*/
func (semaphore *Semaphore) boundedWait() {

	// so that multiple threads don't use the waitgroup at the same time
	semaphore.waitingLock.Lock()
	defer semaphore.waitingLock.Unlock()

	//
	semaphore.waitingGroup.Wait()

	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()

	if semaphore.value == 1 {
		semaphore.waitingGroup.Add(1)
	}

	if semaphore.value == semaphore.maxValue {
		semaphore.signalingGroup.Done()
	}

	semaphore.value--

}

/*
TryWait
Tries to decrement the semaphore’s value and return sucess/true.
If the value is 0 it will return false.
*/
func (semaphore *Semaphore) boundedTryWait() bool {
	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()
	if semaphore.value == 0 {
		return false
	}
	if semaphore.value == 1 {
		semaphore.waitingGroup.Add(1)
	}

	if semaphore.value == semaphore.maxValue {
		semaphore.signalingGroup.Done()
	}

	semaphore.value--
	return true
}

/*
TryWait
Tries to decrement the semaphore’s value and return sucess/true.
If the value is 0 it will return false.
*/
func (semaphore *Semaphore) defaultTryWait() bool {
	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()
	if semaphore.value == 0 {
		return false
	}
	if semaphore.value == 1 {
		semaphore.waitingGroup.Add(1)
	}
	semaphore.value--
	return true
}

/*
Signal
Increments the semaphores value.
If it was 0 it will unblock any waiting threads.
*/
func (semaphore *Semaphore) defaultSignal() {
	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()
	if semaphore.value == 0 {
		semaphore.waitingGroup.Done()
	}
	semaphore.value++
}

/*
Signal
Increments the semaphores value.
If it was 0 it will unblock any waiting threads.
*/
func (semaphore *Semaphore) boundedSignal() {

	// so that multiple threads don't use the waitgroup at the same time
	semaphore.signalingLock.Lock()
	defer semaphore.signalingLock.Unlock()

	semaphore.signalingGroup.Wait()

	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()
	if semaphore.value == 0 {
		semaphore.waitingGroup.Done()
	}

	semaphore.value++

	if semaphore.value == semaphore.maxValue {
		semaphore.signalingGroup.Add(1)
	}

}

func (semaphore *Semaphore) defaultTrySignal() bool {
	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()
	if semaphore.value == 0 {
		semaphore.waitingGroup.Done()
	}

	semaphore.value++

	return true
}

func (semaphore *Semaphore) boundedTrySignal() bool {
	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()
	if semaphore.value == 0 {
		semaphore.waitingGroup.Done()
	}

	if semaphore.value == semaphore.maxValue {
		return false
	}

	semaphore.value++

	if semaphore.value == semaphore.maxValue {
		semaphore.signalingGroup.Add(1)
	}

	return true
}

/*
NewSemaphore
creates a new semaphore
with the option to set the initial value and the maximum value.

# OptionalParameters:
  - initialValue: the initial value of the semaphore. Default is 0.
  - maxValue: the maximum value of the semaphore. Default is math.MaxInt.

# Returns:
  - A pointer to the new semaphore.
*/
func NewSemaphore(arg ...int) *Semaphore {
	var initialValue int
	maxValue := math.MaxInt
	semaphore := Semaphore{}
	semaphore.Signal = semaphore.defaultSignal
	semaphore.Wait = semaphore.defaultWait
	semaphore.TryWait = semaphore.defaultTryWait
	semaphore.TrySignal = semaphore.defaultTrySignal
	switch len(arg) {
	case 2:
		maxValue = arg[1]
		if maxValue < 1 {
			log.Fatalln("maxValue must be greater than 0")
		}
		semaphore.Signal = semaphore.boundedSignal
		semaphore.Wait = semaphore.boundedWait
		semaphore.TryWait = semaphore.boundedTryWait
		semaphore.TrySignal = semaphore.boundedTrySignal
		fallthrough
	case 1:
		initialValue = arg[0]
		if initialValue < 0 {
			log.Fatalln("initialValue must be greater than 0")
		}
		if initialValue > maxValue {
			log.Fatalln("initialValue must be less than maxValue")
		}
		fallthrough
	case 0:

	default:
		log.Fatalln("Invalid number of arguments")
		return nil
	}

	semaphore.maxValue = maxValue
	semaphore.value = initialValue

	if initialValue == 0 {
		semaphore.waitingGroup.Add(1)
	}
	if initialValue == maxValue {
		semaphore.signalingGroup.Add(1)
	}
	return &semaphore
}
