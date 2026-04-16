package prompts

// PauseOption configures a PausePrompt.
type PauseOption func(*PausePrompt)

// PauseWithMessage sets the pause message.

// PausePrompt waits for the user to press Enter.
type PausePrompt struct {
	prompt  *Prompt
	message string
}

func PauseWithMessage(s string) PauseOption { return func(p *PausePrompt) { p.message = s } }

// Pause displays a "Press enter to continue..." prompt and waits.
func Pause(opts ...PauseOption) (bool, error) {
	p := &PausePrompt{
		prompt:  newPrompt(),
		message: "Press enter to continue...",
	}

	for _, opt := range opts {
		opt(p)
	}

	p.prompt.valueFn = func() string { return "true" }
	p.prompt.renderer = func(state State) string {
		return getTheme().PauseRenderer(p, state)
	}

	p.prompt.On("key", func(key string) {
		if key == KeyEnter {
			p.prompt.submit()
		}
	})

	_, err := p.prompt.run()

	if err != nil {
		return false, err
	}

	return true, nil
}
