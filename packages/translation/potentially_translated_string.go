package translation

// PotentiallyTranslatedString is a lazy translation wrapper that holds an
// original string and defers translation until String() is called or
// Translate() is invoked explicitly.  It mirrors the upstream // Illuminate\Translation\PotentiallyTranslatedString.
type PotentiallyTranslatedString struct {
	original   string
	translator *Translator
	translated *string
}

// NewPotentiallyTranslatedString returns a new wrapper around original backed
// by the given Translator.
func NewPotentiallyTranslatedString(original string, translator *Translator) *PotentiallyTranslatedString {
	return &PotentiallyTranslatedString{
		original:   original,
		translator: translator,
	}
}

// Translate executes the translation and stores the result.  Subsequent calls
// to String() return this result.  The method is chainable.
func (p *PotentiallyTranslatedString) Translate(replace map[string]any, locale *string) *PotentiallyTranslatedString {
	result := p.translator.Get(p.original, replace, locale)

	var s string

	if r, ok := result.(string); ok {
		s = r
	} else {
		s = p.original
	}

	p.translated = &s

	return p
}

// TranslateChoice executes a pluralised translation and stores the result.
func (p *PotentiallyTranslatedString) TranslateChoice(number any, replace map[string]any, locale *string) *PotentiallyTranslatedString {
	s := p.translator.Choice(p.original, number, replace, locale)
	p.translated = &s

	return p
}

// Original returns the untranslated source string.
func (p *PotentiallyTranslatedString) Original() string { return p.original }

// String returns the translated string if available, otherwise the original.
// Implements fmt.Stringer.
func (p *PotentiallyTranslatedString) String() string {
	if p.translated != nil {
		return *p.translated
	}

	return p.original
}
