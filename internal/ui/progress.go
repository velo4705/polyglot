package ui

import (
	"fmt"
	"strings"
	"time"
)

// Spinner represents a loading spinner
type Spinner struct {
	message string
	frames  []string
	index   int
	active  bool
	done    chan bool
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		done:    make(chan bool, 1),
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	if !ColorsEnabled {
		fmt.Printf("%s...\n", s.message)
		return
	}

	s.active = true
	go func() {
		for s.active {
			select {
			case <-s.done:
				return
			default:
				frame := s.frames[s.index%len(s.frames)]
				fmt.Printf("\r%s %s", Colorize(Cyan, frame), s.message)
				s.index++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.active = false
	s.done <- true
	fmt.Print("\r" + strings.Repeat(" ", len(s.message)+10) + "\r")
}

// ProgressBar represents a progress bar
type ProgressBar struct {
	total   int
	current int
	width   int
	message string
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int, message string) *ProgressBar {
	return &ProgressBar{
		total:   total,
		current: 0,
		width:   40,
		message: message,
	}
}

// Update updates the progress bar
func (pb *ProgressBar) Update(current int) {
	pb.current = current
	pb.Render()
}

// Increment increments the progress bar
func (pb *ProgressBar) Increment() {
	pb.current++
	pb.Render()
}

// Render renders the progress bar
func (pb *ProgressBar) Render() {
	if !ColorsEnabled {
		fmt.Printf("%s: %d/%d\n", pb.message, pb.current, pb.total)
		return
	}

	percent := float64(pb.current) / float64(pb.total)
	filled := int(percent * float64(pb.width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", pb.width-filled)
	percentStr := fmt.Sprintf("%.0f%%", percent*100)

	fmt.Printf("\r%s [%s] %s %d/%d",
		pb.message,
		Colorize(Green, bar),
		percentStr,
		pb.current,
		pb.total,
	)

	if pb.current >= pb.total {
		fmt.Println()
	}
}

// Complete marks the progress bar as complete
func (pb *ProgressBar) Complete() {
	pb.current = pb.total
	pb.Render()
}
