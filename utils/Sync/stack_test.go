package Sync_test

import (
	"seekourney/utils/Sync"
	"testing"
)

func TestStack(t *testing.T) {
	stack := Sync.NewStack[int]()
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
