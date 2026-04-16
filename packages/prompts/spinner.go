package prompts

import (
	"sync"
	"time"
)

// SpinOption configures a Spin call.
type SpinOption func(*spinConfig)

type spinConfig struct {
	message  string
	interval time.Duration
}

// SpinWithMessage sets the spinner message.
func SpinWithMessage(s string) SpinOption { return func(c *spinConfig) { c.message = s } }

// SpinWithInterval sets the animation interval.
func SpinWithInterval(d time.Duration) SpinOption { return func(c *spinConfig) { c.interval = d } }

var spinnerFrames = []string{"◒", "◐", "◓", "◑"}

// Spin shows an animated spinner while fn executes. Returns the callback result.
func Spin[T any](fn func() (T, error), opts ...SpinOption) (T, error) {
	cfg := &spinConfig{
		message:  "Loading...",
		interval: 100 * time.Millisecond,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	w := getWriter()
	HideCursor(w)

	var (
		mu      sync.Mutex
		stopped bool
		frame   int
	)

	// Animation goroutine.
	done := make(chan struct{})
	go func() {
		defer close(done)

		ticker := time.NewTicker(cfg.interval)

		defer ticker.Stop()

		// Render initial frame.
		mu.Lock()
		w.Write(getTheme().SpinnerRenderer(cfg.message))
		mu.Unlock()

		for {
			select {
			case <-ticker.C:
				mu.Lock()

				if stopped {
					mu.Unlock()

					return
				}

				frame = (frame + 1) % len(spinnerFrames)
				EraseLines(w, 2)
				w.Write(getTheme().SpinnerRenderer(cfg.message))
				mu.Unlock()
			}
		}
	}()

	// Execute callback.
	result, err := fn()

	// Stop animation.
	mu.Lock()
	stopped = true
	mu.Unlock()
	<-done

	// Clean up spinner output.
	EraseLines(w, 2)
	ShowCursor(w)

	return result, err
}
