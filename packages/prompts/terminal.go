package prompts

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Terminal abstracts raw terminal I/O. The default implementation talks to a
// real TTY. For testing, install a FakeTerminal via [Fake].
type Terminal interface {
	// Read reads the next key input from the terminal.
	Read() (string, error)
	// SetTty applies terminal mode changes (e.g., "-icanon -broadcastclient").
	SetTty(mode string) error
	// RestoreTty restores the terminal to its initial mode.
	RestoreTty() error
	// Cols returns the terminal width in columns.
	Cols() int
	// Lines returns the terminal height in lines.
	Lines() int
	// Exit terminates the process.
	Exit()
	// SupportsTrueColor returns true if the terminal supports 24-bit color.
	SupportsTrueColor() bool
}

// Package-level terminal and writer holders, injectable for testing.

// realTerminal implements Terminal using the actual TTY.
type realTerminal struct {
	initialMode string
	cols        int
	lines       int
}

var (
	termMu      sync.Mutex
	terminalFn  = func() Terminal { return newRealTerminal() }
	writerFn    = func() Writer { return &ConsoleWriter{} }
	interactive = true
)

func getTerminal() Terminal {
	termMu.Lock()
	fn := terminalFn
	termMu.Unlock()

	return fn()
}

func getWriter() Writer {
	termMu.Lock()
	fn := writerFn
	termMu.Unlock()

	return fn()
}

func isInteractive() bool {
	termMu.Lock()

	defer termMu.Unlock()

	return interactive
}

var _ Terminal = (*realTerminal)(nil)

func newRealTerminal() *realTerminal {
	t := &realTerminal{}
	t.initDimensions()

	return t
}

func (t *realTerminal) Read() (string, error) {
	buf := make([]byte, 16)
	n, err := os.Stdin.Read(buf)

	if err != nil {
		return "", err
	}

	input := string(buf[:n])

	// If we got a bare ESC, wait briefly to see if more bytes follow
	// (disambiguates ESC key from escape sequence prefix).
	if input == "\x1b" {
		_ = setReadDeadline(os.Stdin, 50*time.Millisecond)
		n2, err2 := os.Stdin.Read(buf)
		_ = clearReadDeadline(os.Stdin)

		if err2 == nil && n2 > 0 {
			input += string(buf[:n2])
		}
	}

	return input, nil
}

func (t *realTerminal) SetTty(mode string) error {
	if t.initialMode == "" {
		out, err := exec.Command("stty", "-g").Output()

		if err == nil {
			t.initialMode = strings.TrimSpace(string(out))
		}
	}

	args := strings.Fields(mode)

	return exec.Command("stty", args...).Run()
}

func (t *realTerminal) RestoreTty() error {
	if t.initialMode == "" {
		return nil
	}

	return exec.Command("stty", t.initialMode).Run()
}

func (t *realTerminal) Cols() int {
	if t.cols == 0 {
		t.initDimensions()
	}

	return t.cols
}

func (t *realTerminal) Lines() int {
	if t.lines == 0 {
		t.initDimensions()
	}

	return t.lines
}

func (t *realTerminal) Exit() {
	os.Exit(1)
}

func (t *realTerminal) SupportsTrueColor() bool {
	colorTerm := os.Getenv("COLORTERM")

	return colorTerm == "truecolor" || colorTerm == "24bit"
}

func (t *realTerminal) initDimensions() {
	t.cols = 80
	t.lines = 24

	out, err := exec.Command("stty", "size").Output()

	if err != nil {
		return
	}

	parts := strings.Fields(strings.TrimSpace(string(out)))

	if len(parts) == 2 {
		if l, err := strconv.Atoi(parts[0]); err == nil && l > 0 {
			t.lines = l
		}

		if c, err := strconv.Atoi(parts[1]); err == nil && c > 0 {
			t.cols = c
		}
	}
}

// setReadDeadline sets a read deadline on the file if it supports it.
func setReadDeadline(f *os.File, d time.Duration) error {
	return f.SetReadDeadline(time.Now().Add(d))
}

// clearReadDeadline clears any read deadline on the file.
func clearReadDeadline(f *os.File) error {
	return f.SetReadDeadline(time.Time{})
}
