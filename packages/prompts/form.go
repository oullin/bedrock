package prompts

import "errors"

// FormBuilder creates multi-step form workflows with conditional steps
// and reversion support.
type FormBuilder struct {
	steps []formStep
}

type formStep struct {
	name                string
	fn                  func(responses map[string]any) (any, error)
	condition           any // bool or func(map[string]any) bool
	ignoreWhenReverting bool
}

// Form returns a new FormBuilder for multi-step form workflows.
func Form() *FormBuilder {
	return &FormBuilder{}
}

// Add registers a custom form step.
func (f *FormBuilder) Add(name string, step func(responses map[string]any) (any, error)) *FormBuilder {
	f.steps = append(f.steps, formStep{
		name: name,
		fn:   step,
	})

	return f
}

// AddIf conditionally registers a form step.
// condition accepts bool or func(map[string]any) bool.
func (f *FormBuilder) AddIf(condition any, name string, step func(responses map[string]any) (any, error)) *FormBuilder {
	f.steps = append(f.steps, formStep{
		name:      name,
		fn:        step,
		condition: condition,
	})

	return f
}

// AddIgnoringWhenReverting adds a step that is skipped during reversion.
func (f *FormBuilder) AddIgnoringWhenReverting(name string, step func(responses map[string]any) (any, error)) *FormBuilder {
	f.steps = append(f.steps, formStep{
		name:                name,
		fn:                  step,
		ignoreWhenReverting: true,
	})

	return f
}

// Text adds a text input step.
func (f *FormBuilder) Text(name, label string, opts ...TextOption) *FormBuilder {
	return f.Add(name, func(_ map[string]any) (any, error) {
		return Text(label, opts...)
	})
}

// Password adds a password input step.
func (f *FormBuilder) Password(name, label string, opts ...PasswordOption) *FormBuilder {
	return f.Add(name, func(_ map[string]any) (any, error) {
		return Password(label, opts...)
	})
}

// Textarea adds a textarea input step.
func (f *FormBuilder) Textarea(name, label string, opts ...TextareaOption) *FormBuilder {
	return f.Add(name, func(_ map[string]any) (any, error) {
		return Textarea(label, opts...)
	})
}

// Number adds a number input step.
func (f *FormBuilder) Number(name, label string, opts ...NumberOption) *FormBuilder {
	return f.Add(name, func(_ map[string]any) (any, error) {
		return Number(label, opts...)
	})
}

// Confirm adds a confirm step.
func (f *FormBuilder) Confirm(name, label string, opts ...ConfirmOption) *FormBuilder {
	return f.Add(name, func(_ map[string]any) (any, error) {
		return Confirm(label, opts...)
	})
}

// Select adds a select step.
func (f *FormBuilder) Select(name, label string, options any, opts ...SelectOption) *FormBuilder {
	return f.Add(name, func(_ map[string]any) (any, error) {
		return Select(label, options, opts...)
	})
}

// MultiSelect adds a multi-select step.
func (f *FormBuilder) MultiSelect(name, label string, options any, opts ...MultiSelectOption) *FormBuilder {
	return f.Add(name, func(_ map[string]any) (any, error) {
		return MultiSelect(label, options, opts...)
	})
}

// Suggest adds a suggest step.
func (f *FormBuilder) Suggest(name, label string, options any, opts ...SuggestOption) *FormBuilder {
	return f.Add(name, func(_ map[string]any) (any, error) {
		return Suggest(label, options, opts...)
	})
}

// Search adds a search step.
func (f *FormBuilder) Search(name, label string, options func(string) map[string]string, opts ...SearchOption) *FormBuilder {
	return f.Add(name, func(_ map[string]any) (any, error) {
		return Search(label, options, opts...)
	})
}

// MultiSearch adds a multi-search step.
func (f *FormBuilder) MultiSearch(name, label string, options func(string) map[string]string, opts ...MultiSearchOption) *FormBuilder {
	return f.Add(name, func(_ map[string]any) (any, error) {
		return MultiSearch(label, options, opts...)
	})
}

// Pause adds a pause step.
func (f *FormBuilder) Pause(opts ...PauseOption) *FormBuilder {
	f.steps = append(f.steps, formStep{
		fn: func(_ map[string]any) (any, error) {
			return Pause(opts...)
		},
	})

	return f
}

// NoteStep adds a note display step.
func (f *FormBuilder) NoteStep(message string) *FormBuilder {
	f.steps = append(f.steps, formStep{
		ignoreWhenReverting: true,
		fn: func(_ map[string]any) (any, error) {
			Note(message)

			return nil, nil
		},
	})

	return f
}

// IntroStep adds an intro display step.
func (f *FormBuilder) IntroStep(message string) *FormBuilder {
	f.steps = append(f.steps, formStep{
		ignoreWhenReverting: true,
		fn: func(_ map[string]any) (any, error) {
			Intro(message)

			return nil, nil
		},
	})

	return f
}

// OutroStep adds an outro display step.
func (f *FormBuilder) OutroStep(message string) *FormBuilder {
	f.steps = append(f.steps, formStep{
		ignoreWhenReverting: true,
		fn: func(_ map[string]any) (any, error) {
			Outro(message)

			return nil, nil
		},
	})

	return f
}

// Submit executes all form steps in order and returns collected responses.
// When a step is cancelled (Ctrl+C), the form reverts to the previous
// interactive step unless it's the first step, in which case the
// cancellation propagates.
func (f *FormBuilder) Submit() (map[string]any, error) {
	responses := make(map[string]any)

	i := 0

	for i < len(f.steps) {
		step := f.steps[i]

		// Check condition.
		if step.condition != nil {
			skip := false

			switch c := step.condition.(type) {
			case bool:
				skip = !c
			case func(map[string]any) bool:
				skip = !c(responses)
			}

			if skip {
				i++

				continue
			}
		}

		result, err := step.fn(responses)

		if err != nil {
			if errors.Is(err, ErrCancelled) {
				// Try to revert to the previous interactive step.
				prev := f.findPreviousInteractiveStep(i)

				if prev < 0 {
					// First step — propagate cancellation.
					return responses, err
				}
				// Clear the response for the step we're reverting to.
				if f.steps[prev].name != "" {
					delete(responses, f.steps[prev].name)
				}

				i = prev

				continue
			}

			return responses, err
		}

		if step.name != "" {
			responses[step.name] = result
		}

		i++
	}

	return responses, nil
}

// findPreviousInteractiveStep finds the index of the previous step that
// should be reverted to (skipping steps marked ignoreWhenReverting).
func (f *FormBuilder) findPreviousInteractiveStep(current int) int {
	for j := current - 1; j >= 0; j-- {
		if !f.steps[j].ignoreWhenReverting {
			return j
		}
	}

	return -1
}
