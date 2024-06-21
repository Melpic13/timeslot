package timeslot

import (
	"testing"
	"time"
)

func TestNewSlot(t *testing.T) {
	_, err := NewSlot(time.Now(), time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
