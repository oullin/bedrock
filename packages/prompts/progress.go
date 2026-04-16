package prompts

import (
	"fmt"
	"sync"
)

// ProgressOption configures a Progress call.
type ProgressOption func(*progressConfig)

type progressConfig struct {
	hint string
}

// ProgressWithHint sets the progress hint text.

// ProgressBar provides methods to control a progress bar during iteration.
type ProgressBar struct {
	mu       sync.Mutex
	label    string
	hint     string
	current  int
	total    int
	writer   Writer
	rendered bool
}

func ProgressWithHint(s string) ProgressOption { return func(c *progressConfig) { c.hint = s } }

// Advance increments the progress by n steps.
func (pb *ProgressBar) Advance(n int) {
	pb.mu.Lock()
	pb.current += n

	if pb.current > pb.total {
		pb.current = pb.total
	}

	pb.render()
	pb.mu.Unlock()
}

// Label updates the progress label.
func (pb *ProgressBar) Label(label string) {
	pb.mu.Lock()
	pb.label = label
	pb.render()
	pb.mu.Unlock()
}

// Hint updates the progress hint.
func (pb *ProgressBar) Hint(hint string) {
	pb.mu.Lock()
	pb.hint = hint
	pb.render()
	pb.mu.Unlock()
}

// Percentage returns the current completion as a decimal ratio (0.0–1.0).
func (pb *ProgressBar) Percentage() float64 {
	pb.mu.Lock()

	defer pb.mu.Unlock()

	return pb.percentage()
}

// percentage is the unlocked version, called while pb.mu is held.
func (pb *ProgressBar) percentage() float64 {
	if pb.total == 0 {
		return 0
	}

	return float64(pb.current) / float64(pb.total)
}

func (pb *ProgressBar) render() {
	if pb.rendered {
		lines := 2

		if pb.hint != "" {
			lines = 3
		}

		EraseLines(pb.writer, lines+1)
	}

	pb.writer.Write(getTheme().ProgressRenderer(pb.label, pb.percentage(), pb.hint))
	pb.rendered = true
}

func (pb *ProgressBar) finish() {
	pb.mu.Lock()

	defer pb.mu.Unlock()

	pb.current = pb.total
	pb.render()
	// Final erase and completion message.
	lines := 2

	if pb.hint != "" {
		lines = 3
	}

	EraseLines(pb.writer, lines+1)
	pb.writer.Write(fmt.Sprintf("  %s %s\n", Green("✓"), pb.label))
}

// Progress displays a progress bar, mapping over items with a callback.
func Progress[T any, R any](label string, steps []T, fn func(T, *ProgressBar) R, opts ...ProgressOption) ([]R, error) {
	cfg := &progressConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	w := getWriter()
	HideCursor(w)

	pb := &ProgressBar{
		label:  label,
		hint:   cfg.hint,
		total:  len(steps),
		writer: w,
	}

	// Render initial state.
	pb.mu.Lock()
	pb.render()
	pb.mu.Unlock()

	results := make([]R, len(steps))

	for i, step := range steps {
		results[i] = fn(step, pb)
		pb.Advance(1)
	}

	pb.finish()
	ShowCursor(w)

	return results, nil
}
