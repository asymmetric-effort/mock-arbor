package arbor

import (
	"testing"
	"time"
)

func TestStartTimeStable(t *testing.T) {
	t1 := StartTime()
	time.Sleep(10 * time.Millisecond)
	t2 := StartTime()
	if !t1.Equal(t2) {
		t.Fatalf("StartTime not stable across calls: %v vs %v", t1, t2)
	}
	if t1.After(time.Now()) {
		t.Fatalf("StartTime is in the future: %v", t1)
	}
}
