package prompts

// GridOption configures Grid display.
type GridOption func(*gridConfig)

type gridConfig struct {
	maxWidth int
}

// GridWithMaxWidth sets the maximum grid width.
func GridWithMaxWidth(w int) GridOption {
	return func(c *gridConfig) { c.maxWidth = w }
}

// Grid displays items in a responsive grid layout.
func Grid(items []string, opts ...GridOption) {
	cfg := &gridConfig{maxWidth: 80}

	for _, opt := range opts {
		opt(cfg)
	}

	w := getWriter()
	w.Write(getTheme().GridRenderer(items, cfg.maxWidth))
}
