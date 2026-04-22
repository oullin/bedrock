package prompts

import (
	"strconv"
)

// NumberOption configures a NumberPrompt.
type NumberOption func(*NumberPrompt)

// NumberWithPlaceholder sets placeholder text.

// NumberWithDefault sets the default value.

// NumberWithRequired makes the prompt required.

// NumberWithValidate sets validation.

// NumberWithHint sets hint text.

// NumberWithTransform sets the transform function.

// NumberWithMin sets the minimum allowed value.

// NumberWithMax sets the maximum allowed value.

// NumberWithStep sets the step increment for arrow keys.

// NumberPrompt handles numeric input with optional min/max/step.
type NumberPrompt struct {
	prompt      *Prompt
	TypedValue  TypedValue
	label       string
	placeholder string
	hint        string
	min         *int
	max         *int
	step        int
}

func NumberWithPlaceholder(s string) NumberOption {
	return func(p *NumberPrompt) { p.placeholder = s }
}

func NumberWithDefault(n int) NumberOption {
	return func(p *NumberPrompt) {
		p.TypedValue.SetValue(strconv.Itoa(n))
		p.prompt.defaultValue = strconv.Itoa(n)
	}
}

func NumberWithRequired(v any) NumberOption {
	return func(p *NumberPrompt) { p.prompt.required = v }
}

func NumberWithValidate(fn ValidateFunc) NumberOption {
	return func(p *NumberPrompt) { p.prompt.validate = fn }
}

func NumberWithHint(s string) NumberOption { return func(p *NumberPrompt) { p.hint = s } }

func NumberWithTransform(fn TransformFunc) NumberOption {
	return func(p *NumberPrompt) { p.prompt.transform = fn }
}

func NumberWithMin(n int) NumberOption { return func(p *NumberPrompt) { p.min = &n } }

func NumberWithMax(n int) NumberOption { return func(p *NumberPrompt) { p.max = &n } }

func NumberWithStep(n int) NumberOption { return func(p *NumberPrompt) { p.step = n } }

// Number displays a numeric input prompt and returns the entered integer.
func Number(label string, opts ...NumberOption) (int, error) {
	p := &NumberPrompt{
		prompt: newPrompt(),
		label:  label,
		step:   1,
	}

	for _, opt := range opts {
		opt(p)
	}

	validate := p.prompt.validate
	if p.min != nil || p.max != nil || validate != nil {
		p.prompt.validate = func(value string) string {
			if p.min != nil || p.max != nil {
				if n, err := strconv.Atoi(value); err == nil {
					if p.min != nil && n < *p.min {
						return "Minimum value is " + strconv.Itoa(*p.min) + "."
					}

					if p.max != nil && n > *p.max {
						return "Maximum value is " + strconv.Itoa(*p.max) + "."
					}
				}
			}

			if validate != nil {
				return validate(value)
			}

			return ""
		}
	}

	p.prompt.valueFn = p.TypedValue.Value
	p.prompt.renderer = func(state State) string {
		return getTheme().NumberRenderer(p, state)
	}

	// Handle arrow keys for increment/decrement.
	p.prompt.On("key", func(key string) {
		if p.prompt.state != StateActive && p.prompt.state != StateError {
			return
		}

		switch {
		case IsUpKey(key):
			p.increment()
		case IsDownKey(key):
			p.decrement()
		}
	})

	p.TypedValue.TrackTypedValue(p.prompt, p.prompt.submit)

	result, err := p.prompt.run()

	if err != nil {
		return 0, err
	}

	n, err := strconv.Atoi(result)

	if err != nil {
		return 0, err
	}

	return n, nil
}

func (p *NumberPrompt) increment() {
	val := p.currentInt()

	if p.Value() == "" && p.min != nil {
		val = *p.min
	} else {
		val += p.step
	}

	if p.max != nil && val > *p.max {
		val = *p.max
	}

	p.TypedValue.SetValue(strconv.Itoa(val))
}

func (p *NumberPrompt) decrement() {
	val := p.currentInt()

	if p.Value() == "" && p.min != nil {
		val = *p.min
	} else {
		val -= p.step
	}

	if p.min != nil && val < *p.min {
		val = *p.min
	}

	p.TypedValue.SetValue(strconv.Itoa(val))
}

func (p *NumberPrompt) currentInt() int {
	n, _ := strconv.Atoi(p.TypedValue.Value())

	return n
}

// Value returns the current typed value.
func (p *NumberPrompt) Value() string { return p.TypedValue.Value() }
