// Package watch wraps fsnotify with a debounce + ignore list so the brain
// scanner can rescan on `*.go` changes without firing on every keystroke
// of a multi-file save. This is the only internal new dependency the plan
// allows; everything else routes through an existing bedrock package.
package watch

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Options configures the watcher. Zero values produce sensible defaults.
type Options struct {
	Roots          []string
	Debounce       time.Duration
	IgnoreSegments []string
	Extensions     []string
}

// Run blocks until ctx is done, calling onChange every time a debounced
// burst of relevant events lands. onChange is invoked from a single
// goroutine — implementations may rescan without locking.
func Run(ctx context.Context, opts Options, onChange func()) error {
	if opts.Debounce == 0 {
		opts.Debounce = 200 * time.Millisecond
	}
	if len(opts.IgnoreSegments) == 0 {
		opts.IgnoreSegments = []string{
			".git", "vendor", "node_modules", "storage", ".turbo", "dist",
		}
	}
	if len(opts.Extensions) == 0 {
		opts.Extensions = []string{".go"}
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()
	for _, root := range opts.Roots {
		if err := addRecursive(w, root, opts.IgnoreSegments); err != nil {
			return err
		}
	}

	var timer *time.Timer
	fire := func() { onChange() }

	for {
		select {
		case <-ctx.Done():
			return nil
		case ev, ok := <-w.Events:
			if !ok {
				return nil
			}
			if !relevant(ev.Name, opts) {
				continue
			}
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(opts.Debounce, fire)
			// follow newly-created directories
			if ev.Op&fsnotify.Create != 0 {
				_ = addRecursive(w, ev.Name, opts.IgnoreSegments)
			}
		case <-w.Errors:
		}
	}
}

func relevant(path string, opts Options) bool {
	for _, seg := range opts.IgnoreSegments {
		if strings.Contains(path, string(filepath.Separator)+seg+string(filepath.Separator)) {
			return false
		}
		if strings.HasSuffix(path, string(filepath.Separator)+seg) {
			return false
		}
	}
	for _, ext := range opts.Extensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}

func addRecursive(w *fsnotify.Watcher, root string, ignore []string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		base := d.Name()
		for _, seg := range ignore {
			if base == seg {
				return fs.SkipDir
			}
		}
		return w.Add(path)
	})
}
