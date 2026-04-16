package prompts

// PasswordOption configures a PasswordPrompt.
type PasswordOption func(*PasswordPrompt)

// PasswordWithPlaceholder sets placeholder text.

// PasswordWithRequired makes the prompt required.

// PasswordWithValidate sets validation.

// PasswordWithHint sets the hint text.

// PasswordWithTransform sets the transformation function.

// PasswordPrompt handles masked password input.
type PasswordPrompt struct {
	prompt      *Prompt
	TypedValue  TypedValue
	label       string
	placeholder string
	hint        string
}

func PasswordWithPlaceholder(s string) PasswordOption {
	return func(p *PasswordPrompt) { p.placeholder = s }
}

func PasswordWithRequired(v any) PasswordOption {
	return func(p *PasswordPrompt) { p.prompt.required = v }
}

func PasswordWithValidate(fn ValidateFunc) PasswordOption {
	return func(p *PasswordPrompt) { p.prompt.validate = fn }
}

func PasswordWithHint(s string) PasswordOption {
	return func(p *PasswordPrompt) { p.hint = s }
}

func PasswordWithTransform(fn TransformFunc) PasswordOption {
	return func(p *PasswordPrompt) { p.prompt.transform = fn }
}

// Password displays a masked password input prompt.
func Password(label string, opts ...PasswordOption) (string, error) {
	p := &PasswordPrompt{
		prompt: newPrompt(),
		label:  label,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.prompt.valueFn = p.TypedValue.Value
	p.prompt.renderer = func(state State) string {
		return getTheme().PasswordRenderer(p, state)
	}

	p.TypedValue.TrackTypedValue(p.prompt, p.prompt.submit)

	return p.prompt.run()
}

// Value returns the current typed value.
func (p *PasswordPrompt) Value() string { return p.TypedValue.Value() }
