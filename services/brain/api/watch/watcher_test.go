package watch

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestRelevant(t *testing.T) {
	defaults := Options{
		IgnoreSegments: []string{".git", "vendor", "node_modules", "storage"},
		Extensions:     []string{".go"},
	}

	sep := string(filepath.Separator)
	cases := []struct {
		path string
		opts Options
		want bool
	}{
		{filepath.Join("x", "foo.go"), defaults, true},
		{filepath.Join("x", "foo.txt"), defaults, false},
		{sep + filepath.Join("x", ".git", "foo.go"), defaults, false},
		{sep + filepath.Join("x", "vendor", "y.go"), defaults, false},
		{sep + filepath.Join("x", "node_modules"), defaults, false},
		{filepath.Join("x", "y.vue"), Options{Extensions: []string{".vue"}}, true},
		{filepath.Join("x", "y.go"), Options{Extensions: []string{".vue"}}, false},
	}

	for _, c := range cases {
		if got := relevant(c.path, c.opts); got != c.want {
			t.Errorf("relevant(%q, %v) = %v, want %v", c.path, c.opts.Extensions, got, c.want)
		}
	}
}

func TestAddRecursiveSkipsIgnoredDirs(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, root,
		"pkg/sub",
		"vendor/lib",
		"node_modules/foo",
		".git/objects",
	)

	w, err := fsnotify.NewWatcher()

	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}

	t.Cleanup(func() { _ = w.Close() })

	if err := addRecursive(w, root, []string{".git", "vendor", "node_modules"}); err != nil {
		t.Fatalf("addRecursive: %v", err)
	}

	watched := w.WatchList()
	want := map[string]bool{
		root:                              false,
		filepath.Join(root, "pkg"):        false,
		filepath.Join(root, "pkg", "sub"): false,
	}

	for _, p := range watched {
		if _, ok := want[p]; ok {
			want[p] = true
		}

		for _, banned := range []string{"vendor", "node_modules", ".git"} {
			if filepath.Base(p) == banned || filepath.Base(filepath.Dir(p)) == banned {
				t.Errorf("watching ignored path %q", p)
			}
		}
	}

	for path, found := range want {
		if !found {
			t.Errorf("missing watch entry %q (watched: %v)", path, watched)
		}
	}
}

func TestRunDebouncesAndFires(t *testing.T) {
	root := t.TempDir()

	var calls int32
	var once sync.Once
	fired := make(chan struct{}, 1)

	onChange := func() {
		atomic.AddInt32(&calls, 1)
		once.Do(func() { close(fired) })
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{
			Roots:    []string{root},
			Debounce: 30 * time.Millisecond,
		}, onChange)
	}()

	// Give the watcher a moment to attach.
	time.Sleep(50 * time.Millisecond)

	if err := os.WriteFile(filepath.Join(root, "trigger.go"), []byte("package x"), 0o644); err != nil {
		t.Fatalf("write trigger: %v", err)
	}

	select {
	case <-fired:
	case <-time.After(2 * time.Second):
		t.Fatal("onChange did not fire within 2s")
	}

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run returned error: %v", err)
		}
	case <-time.After(time.Second):
		t.Error("Run did not return after ctx.Cancel within 1s")
	}

	if atomic.LoadInt32(&calls) < 1 {
		t.Errorf("calls = %d, want >= 1", calls)
	}
}

func mkdirs(t *testing.T, root string, paths ...string) {
	t.Helper()

	for _, p := range paths {
		if err := os.MkdirAll(filepath.Join(root, p), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", p, err)
		}
	}
}
