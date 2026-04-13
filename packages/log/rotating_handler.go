package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// RotatingHandler writes log records to daily-rotated files. Old files beyond
// the configured retention are pruned automatically.
type RotatingHandler struct {
	ProcessableHandler
	FormattableHandler
	mu          sync.Mutex
	basePath    string
	level       Level
	maxFiles    int
	currentDate string
	writer      *os.File
	perm        os.FileMode
}

var _ Handler = (*RotatingHandler)(nil)

// NewRotatingHandler creates a handler that rotates log files daily. maxFiles
// controls how many rotated files to keep (0 means unlimited).
func NewRotatingHandler(basePath string, maxFiles int, level Level) *RotatingHandler {
	return &RotatingHandler{
		basePath: basePath,
		level:    level,
		maxFiles: maxFiles,
		perm:     0644,
	}
}

// Handle writes the record to the current day's file, rotating if necessary.
func (h *RotatingHandler) Handle(record Record) error {
	if !h.IsHandling(record.Level) {
		return nil
	}

	h.mu.Lock()

	defer h.mu.Unlock()

	date := record.Time.Format("2006-01-02")

	if date != h.currentDate || h.writer == nil {
		if err := h.rotate(date); err != nil {
			return err
		}
	}

	record = h.ProcessRecord(record)

	formatted, err := h.GetFormatter().Format(record)

	if err != nil {
		return err
	}

	_, err = h.writer.Write(formatted)

	return err
}

// IsHandling reports whether this handler handles the given level.
func (h *RotatingHandler) IsHandling(level Level) bool {
	return level >= h.level
}

// Close closes the current file writer.
func (h *RotatingHandler) Close() error {
	h.mu.Lock()

	defer h.mu.Unlock()

	if h.writer != nil {
		err := h.writer.Close()
		h.writer = nil

		return err
	}

	return nil
}

func (h *RotatingHandler) rotate(date string) error {
	if h.writer != nil {
		h.writer.Close()
		h.writer = nil
	}

	path := h.filePath(date)

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, h.perm)

	if err != nil {
		return err
	}

	h.writer = f
	h.currentDate = date

	if h.maxFiles > 0 {
		h.pruneOld()
	}

	return nil
}

func (h *RotatingHandler) filePath(date string) string {
	ext := filepath.Ext(h.basePath)
	base := strings.TrimSuffix(h.basePath, ext)

	if ext == "" {
		ext = ".log"
	}

	return fmt.Sprintf("%s-%s%s", base, date, ext)
}

func (h *RotatingHandler) pruneOld() {
	ext := filepath.Ext(h.basePath)
	base := strings.TrimSuffix(h.basePath, ext)

	if ext == "" {
		ext = ".log"
	}

	pattern := fmt.Sprintf("%s-*%s", base, ext)

	matches, err := filepath.Glob(pattern)

	if err != nil || len(matches) <= h.maxFiles {
		return
	}

	sort.Strings(matches)

	toRemove := matches[:len(matches)-h.maxFiles]

	for _, path := range toRemove {
		now := time.Now()
		_ = os.Chtimes(path, now, now)
		os.Remove(path)
	}
}
