package prompts

// TextOption configures a TextPrompt.
type TextOption func(*TextPrompt)

// TextWithPlaceholder sets the placeholder text shown when the input is empty.

// TextWithDefault sets the default value.

// TextWithRequired makes the prompt required. Pass a string for a custom message.

// TextWithValidate sets the validation function.

// TextWithHint sets the hint text displayed below the prompt.

// TextWithTransform sets the transformation function applied to the value.

// TextPrompt handles single-line text input.
type TextPrompt struct {
	prompt      *Prompt
	TypedValue  TypedValue
	label       string
	placeholder string
	hint        string
}

func TextWithPlaceholder(s string) TextOption { return func(p *TextPrompt) { p.placeholder = s } }

func TextWithDefault(s string) TextOption {
	return func(p *TextPrompt) {
		p.TypedValue.SetValue(s)
		p.prompt.defaultValue = s
	}
}

func TextWithRequired(v any) TextOption { return func(p *TextPrompt) { p.prompt.required = v } }

func TextWithValidate(fn ValidateFunc) TextOption {
	return func(p *TextPrompt) { p.prompt.validate = fn }
}

func TextWithHint(s string) TextOption { return func(p *TextPrompt) { p.hint = s } }

func TextWithTransform(fn TransformFunc) TextOption {
	return func(p *TextPrompt) { p.prompt.transform = fn }
}

// Text displays a single-line text input prompt and returns the entered string.
func Text(label string, opts ...TextOption) (string, error) {
	p := &TextPrompt{
		prompt: newPrompt(),
		label:  label,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.prompt.valueFn = p.TypedValue.Value
	p.prompt.renderer = func(state State) string {
		return getTheme().TextRenderer(p, state)
	}

	p.TypedValue.TrackTypedValue(p.prompt, p.prompt.submit)

	return p.prompt.run()
}

// Value returns the current typed value.
func (p *TextPrompt) Value() string { return p.TypedValue.Value() }
