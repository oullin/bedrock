package prompts

import "strconv"

// ConfirmOption configures a ConfirmPrompt.
type ConfirmOption func(*ConfirmPrompt)

// ConfirmWithDefault sets the default confirmation state.

// ConfirmWithYes sets the "yes" label.

// ConfirmWithNo sets the "no" label.

// ConfirmWithHint sets the hint text.

// ConfirmWithRequired makes the prompt required.

// ConfirmWithValidate sets validation.

// ConfirmWithTransform sets the transform function.

// ConfirmPrompt handles yes/no confirmation.
type ConfirmPrompt struct {
	prompt    *Prompt
	label     string
	confirmed bool
	yesLabel  string
	noLabel   string
	hint      string
}

func ConfirmWithDefault(v bool) ConfirmOption {
	return func(p *ConfirmPrompt) { p.confirmed = v }
}

func ConfirmWithYes(s string) ConfirmOption { return func(p *ConfirmPrompt) { p.yesLabel = s } }

func ConfirmWithNo(s string) ConfirmOption { return func(p *ConfirmPrompt) { p.noLabel = s } }

func ConfirmWithHint(s string) ConfirmOption { return func(p *ConfirmPrompt) { p.hint = s } }

func ConfirmWithRequired(v any) ConfirmOption {
	return func(p *ConfirmPrompt) { p.prompt.required = v }
}

func ConfirmWithValidate(fn ValidateFunc) ConfirmOption {
	return func(p *ConfirmPrompt) { p.prompt.validate = fn }
}

func ConfirmWithTransform(fn TransformFunc) ConfirmOption {
	return func(p *ConfirmPrompt) { p.prompt.transform = fn }
}

// Confirm displays a yes/no confirmation prompt and returns the result.
func Confirm(label string, opts ...ConfirmOption) (bool, error) {
	p := &ConfirmPrompt{
		prompt:    newPrompt(),
		label:     label,
		confirmed: true,
		yesLabel:  "Yes",
		noLabel:   "No",
	}

	for _, opt := range opts {
		opt(p)
	}

	p.prompt.defaultValue = strconv.FormatBool(p.confirmed)
	p.prompt.valueFn = func() string {
		return strconv.FormatBool(p.confirmed)
	}

	p.prompt.renderer = func(state State) string {
		return getTheme().ConfirmRenderer(p, state)
	}

	p.prompt.On("key", func(key string) {
		if p.prompt.state != StateActive && p.prompt.state != StateError {
			return
		}

		switch {
		case key == KeyEnter:
			p.prompt.submit()
		case IsLeftKey(key) || key == "y" || key == "Y":
			p.confirmed = true
		case IsRightKey(key) || key == "n" || key == "N":
			p.confirmed = false
		case key == KeyTab:
			p.confirmed = !p.confirmed
		}
	})

	result, err := p.prompt.run()

	if err != nil {
		return false, err
	}

	return result == "true", nil
}
