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
		signal
		Increments the semaphores value.
		If it was 0 it will unblock any blocked Wait thread.
		If bounded, it will block if the value is currently the maximum value.
	*/
	signal func()

	/*
		Wait
		Decrement the semaphore’s value.
		It will block if the value is currently 0.
		If bounded, it will unblock any blocked Signal thread.
	*/
	wait func()

	/*
		TryWait
		Tries to decrement the semaphore’s value and return sucess/true.
		If the value is 0 it will return false.
	*/
	tryWait func() bool

	/*
		TrySignal
		Tries to increment the semaphore’s value and return sucess/true.
		If the value is currently the maximum value it will return false.
	*/
	trySignal func() bool

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
	semaphore.signal = semaphore.defaultSignal
	semaphore.wait = semaphore.defaultWait
	semaphore.tryWait = semaphore.defaultTryWait
	semaphore.trySignal = semaphore.defaultTrySignal
	switch len(arg) {
	case 2:
		maxValue = arg[1]
		if maxValue < 1 {
			log.Fatalln("maxValue must be greater than 0")
		}
		semaphore.signal = semaphore.boundedSignal
		semaphore.wait = semaphore.boundedWait
		semaphore.tryWait = semaphore.boundedTryWait
		semaphore.trySignal = semaphore.boundedTrySignal
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

/*
signal
Increments the semaphores value.
If it was 0 it will unblock any blocked Wait thread.
If bounded, it will block if the value is currently the maximum value.
*/
func (semaphore *Semaphore) Signal() {
	semaphore.signal()
}

/*
Wait
Decrement the semaphore’s value.
It will block if the value is currently 0.
If bounded, it will unblock any blocked Signal thread.
*/
func (semaphore *Semaphore) Wait() {
	semaphore.wait()
}

/*
TryWait
Tries to decrement the semaphore’s value and return sucess/true.
If the value is 0 it will return false.
*/
func (semaphore *Semaphore) TryWait() bool {
	return semaphore.tryWait()
}

/*
TrySignal
Tries to increment the semaphore’s value and return sucess/true.
If the value is currently the maximum value it will return false.
*/
func (semaphore *Semaphore) TrySignal() bool {
	return semaphore.trySignal()
}

/*
defaultWait
Decrement the semaphore’s value
and will block if the value is currently 0.
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
defaultTryWait
Tries to decrement the semaphore’s value and return true.
If the value is 0 it will return false.
*/
func (semaphore *Semaphore) defaultTryWait() bool {

	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()

	if semaphore.value == 0 {
		return false
	}
	// a value of one means lowering it once more will make waits block
	if semaphore.value == 1 {

		if semaphore.waitingLock.TryLock() {
			// if there is no normal wait waiting; continue as normal
			defer semaphore.waitingLock.Unlock()

			semaphore.waitingGroup.Add(1)
		} else {
			// a normal wait was already being processed and
			// will be given priority
			return false
		}
	}
	semaphore.value--

	return true
}

/*
defaultSignal
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
defaultTrySignal
Tries to increment the semaphore’s value and return true.
Currently always returns true.
*/
func (semaphore *Semaphore) defaultTrySignal() bool {

	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()

	if semaphore.value == 0 {
		semaphore.waitingGroup.Done()
	}

	semaphore.value++

	return true
}

/*
boundedWait
Decrement the semaphore’s value
and will block if the value is currently 0.
adjusted for a bounded semaphore
*/
func (semaphore *Semaphore) boundedWait() {

	// so that multiple threads don't use the waitgroup at the same time
	semaphore.waitingLock.Lock()
	defer semaphore.waitingLock.Unlock()

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
boundedTryWait
Tries to decrement the semaphore’s value and return true.
If the value is 0 it will return false.
adjusted for a bounded semaphore
*/
func (semaphore *Semaphore) boundedTryWait() bool {

	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()

	if semaphore.value == 0 {
		return false
	}
	// a value of one means lowering it once more will make waits block
	if semaphore.value == 1 {

		if semaphore.waitingLock.TryLock() {
			// if there is no normal wait waiting; continue as normal
			defer semaphore.waitingLock.Unlock()

			semaphore.waitingGroup.Add(1)
		} else {
			// a normal wait was already being processed and
			// will be given priority
			return false
		}
	}

	if semaphore.value == semaphore.maxValue {
		semaphore.signalingGroup.Done()
	}

	semaphore.value--
	return true
}

/*
boundedSignal
Increments the semaphores value.
If it was 0 it will unblock any waiting threads.
If the value is currently the maximum value it will block.
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

/*
boundedTrySignal
Tries to increment the semaphore’s value and return true.
If the value is currently the maximum value it will return false.
*/
func (semaphore *Semaphore) boundedTrySignal() bool {

	semaphore.syncLock.Lock()
	defer semaphore.syncLock.Unlock()

	if semaphore.value == 0 {
		semaphore.waitingGroup.Done()
	}

	if semaphore.value == semaphore.maxValue {
		return false
	}
	// value of maxValue-1 means signaling it once more will make signals block
	if semaphore.value == semaphore.maxValue-1 {

		if semaphore.signalingLock.TryLock() {
			// if there is no normal signal waiting; continue as normal
			defer semaphore.signalingLock.Unlock()

			semaphore.signalingGroup.Add(1)
		} else {
			// a normal signal was already being processed and
			// will be given priority
			return false
		}
	}
	semaphore.value++

	return true
}
