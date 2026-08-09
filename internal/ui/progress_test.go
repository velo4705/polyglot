package ui

import (
	"sync"
	"testing"
)

func TestNewSpinner(t *testing.T) {
	s := NewSpinner("loading")
	if s.message != "loading" {
		t.Errorf("message = %q, want %q", s.message, "loading")
	}
	if len(s.frames) == 0 {
		t.Error("frames should not be empty")
	}
}

func TestSpinnerStartStopNoColors(t *testing.T) {
	ColorsEnabled = false
	s := NewSpinner("test")
	s.Start()
	s.Stop()
}

func TestSpinnerDoubleStop(t *testing.T) {
	ColorsEnabled = false
	s := NewSpinner("test")
	s.Start()
	s.Stop()
	// Second stop should not panic or block
	s.Stop()
}

func TestSpinnerConcurrentStartStop(t *testing.T) {
	ColorsEnabled = false
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := NewSpinner("concurrent")
			s.Start()
			s.Stop()
		}()
	}
	wg.Wait()
}

func TestProgressBar(t *testing.T) {
	ColorsEnabled = false
	pb := NewProgressBar(10, "progress")
	if pb.total != 10 {
		t.Errorf("total = %d, want 10", pb.total)
	}
	if pb.current != 0 {
		t.Errorf("current = %d, want 0", pb.current)
	}
}

func TestProgressBarIncrement(t *testing.T) {
	ColorsEnabled = false
	pb := NewProgressBar(5, "test")
	pb.Increment()
	if pb.current != 1 {
		t.Errorf("after Increment, current = %d, want 1", pb.current)
	}
}

func TestProgressBarUpdate(t *testing.T) {
	ColorsEnabled = false
	pb := NewProgressBar(10, "test")
	pb.Update(7)
	if pb.current != 7 {
		t.Errorf("after Update(7), current = %d, want 7", pb.current)
	}
}

func TestProgressBarComplete(t *testing.T) {
	ColorsEnabled = false
	pb := NewProgressBar(5, "test")
	pb.Complete()
	if pb.current != 5 {
		t.Errorf("after Complete, current = %d, want 5", pb.current)
	}
}
