package prompts

import "sync"

// Theme holds a set of renderer functions for each prompt type.
type Theme struct {
	TextRenderer         func(p *TextPrompt, state State) string
	PasswordRenderer     func(p *PasswordPrompt, state State) string
	ConfirmRenderer      func(p *ConfirmPrompt, state State) string
	NumberRenderer       func(p *NumberPrompt, state State) string
	PauseRenderer        func(p *PausePrompt, state State) string
	SelectRenderer       func(p *SelectPrompt, state State) string
	MultiSelectRenderer  func(p *MultiSelectPrompt, state State) string
	TextareaRenderer     func(p *TextareaPrompt, state State) string
	SuggestRenderer      func(p *SuggestPrompt, state State) string
	AutocompleteRenderer func(p *AutocompletePrompt, state State) string
	SearchRenderer       func(p *SearchPrompt, state State) string
	MultiSearchRenderer  func(p *MultiSearchPrompt, state State) string
	NoteRenderer         func(message string, noteType string) string
	TableRenderer        func(headers []string, rows [][]string) string
	GridRenderer         func(items []string, maxWidth int) string
	SpinnerRenderer      func(message string) string
	ProgressRenderer     func(label string, percentage float64, hint string) string
	DataTableRenderer    func(p *DataTablePrompt, state State) string
}

var (
	themeMu     sync.RWMutex
	themes      = map[string]*Theme{}
	activeTheme = "default"
)

// RegisterTheme registers a named theme.
func RegisterTheme(name string, theme *Theme) {
	themeMu.Lock()
	themes[name] = theme
	themeMu.Unlock()
}

// SetTheme sets the active theme by name.
func SetTheme(name string) {
	themeMu.Lock()
	activeTheme = name
	themeMu.Unlock()
}

// getTheme returns the active theme, falling back to "default".
func getTheme() *Theme {
	themeMu.RLock()

	defer themeMu.RUnlock()

	if t, ok := themes[activeTheme]; ok {
		return t
	}

	if t, ok := themes["default"]; ok {
		return t
	}

	return &Theme{}
}

func init() {
	RegisterTheme("default", defaultTheme())
}
