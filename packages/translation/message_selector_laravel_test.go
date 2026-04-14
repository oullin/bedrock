package translation_test

// Ports of Illuminate\Tests\Translation\TranslationMessageSelectorTest.

import (
	"fmt"
	"testing"

	"github.com/bedrock/packages/translation"
)

func TestMessageSelectorChoose(t *testing.T) {
	t.Parallel()

	s := translation.NewMessageSelector()

	cases := []struct {
		line   string
		n      float64
		locale string
		want   string
	}{
		// Simple pipe-delimited (no conditions).
		{"foo", 1, "en", "foo"},
		{"foo|bar", 1, "en", "foo"},
		{"foo|bar", 10, "en", "bar"},

		// Exact number conditions {N}.
		{"{0} zero|{1} one|[2,*] many", 0, "en", "zero"},
		{"{0} zero|{1} one|[2,*] many", 1, "en", "one"},
		{"{0} zero|{1} one|[2,*] many", 2, "en", "many"},
		{"{0} zero|{1} one|[2,*] many", 10, "en", "many"},

		// Range conditions [M,N].
		{"[2,*] many", 5, "en", "many"},
		{"[*,-1] negative|{0} zero|[1,*] positive", -5, "en", "negative"},
		{"[*,-1] negative|{0} zero|[1,*] positive", 0, "en", "zero"},
		{"[*,-1] negative|{0} zero|[1,*] positive", 5, "en", "positive"},

		// Mixed conditions.
		{"{0} no apples|{1} apple|[2,*] apples", 0, "en", "no apples"},
		{"{0} no apples|{1} apple|[2,*] apples", 1, "en", "apple"},
		{"{0} no apples|{1} apple|[2,*] apples", 5, "en", "apples"},

		// Floats use plural form via truncation.
		{"one|many", 0.5, "en", "many"},
		{"one|many", 1.5, "en", "many"},
		{"one|many", 2.3, "en", "many"},

		// Zero treated as plural in English.
		{"singular|plural", 0, "en", "plural"},

		// Multiline content preserved.
		{"line1\nline2|other", 2, "en", "other"},

		// Empty string result.
		{"{1} |{0} empty", 1, "en", ""},
		{"{1} |{0} empty", 0, "en", "empty"},

		// Whitespace tolerance.
		{"  foo  |  bar  ", 1, "en", "foo"},

		// Bracket range — lower bound inclusive.
		{"[1,*] positive", 1, "en", "positive"},

		// Upper bound only.
		{"[*,4] few|[5,*] many", 3, "en", "few"},
		{"[*,4] few|[5,*] many", 5, "en", "many"},
	}

	for _, c := range cases {
		c := c
		t.Run(fmt.Sprintf("%s_%v_%s", c.line, c.n, c.locale), func(t *testing.T) {
			t.Parallel()
			got := s.Choose(c.line, c.n, c.locale)

			if got != c.want {
				t.Errorf("Choose(%q, %v, %q) = %q; want %q", c.line, c.n, c.locale, got, c.want)
			}
		})
	}
}

func TestMessageSelectorChoosePluralizesFloats(t *testing.T) {
	t.Parallel()

	s := translation.NewMessageSelector()

	// 0.5 is not 1, so English returns plural.
	got := s.Choose("singular|plural", 0.5, "en")

	if got != "plural" {
		t.Errorf("Choose(singular|plural, 0.5, en) = %q; want %q", got, "plural")
	}
}

func TestMessageSelectorChooseWithFloatDoesNotPanic(t *testing.T) {
	t.Parallel()

	s := translation.NewMessageSelector()

	// Polish pluralisation uses modulo arithmetic; float must not panic.
	// Mirrors Laravel's testChooseWithFloatDoesNotTriggerDeprecation.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Choose panicked with float for Polish locale: %v", r)
		}
	}()

	_ = s.Choose("jeden|:count elementy|:count elementów", 1.5, "pl")
}
