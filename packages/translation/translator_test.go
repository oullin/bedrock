package translation_test

import (
	"testing"

	"github.com/bedrock/packages/translation"
)

func TestTranslatorGetLocale(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "fr")

	if tr.GetLocale() != "fr" {
		t.Errorf("GetLocale() = %q", tr.GetLocale())
	}
}

func TestTranslatorSetLocale(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "en")
	tr.SetLocale("de")

	if tr.GetLocale() != "de" {
		t.Errorf("SetLocale did not change locale: %q", tr.GetLocale())
	}
}

func TestTranslatorSetLocaleRejectsPathSeparators(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "en")
	tr.SetLocale("en/evil")

	if tr.GetLocale() != "en" {
		t.Errorf("SetLocale accepted path separator; locale = %q", tr.GetLocale())
	}

	tr.SetLocale(`en\evil`)

	if tr.GetLocale() != "en" {
		t.Errorf("SetLocale accepted backslash; locale = %q", tr.GetLocale())
	}
}

func TestTranslatorGetFallback(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "en")
	tr.SetFallback("es")

	if tr.GetFallback() != "es" {
		t.Errorf("GetFallback() = %q", tr.GetFallback())
	}
}

func TestTranslatorAddLines(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "en")
	tr.AddLines(map[string]any{"hello": "world"}, "en")

	got := tr.Get("hello", nil, nil)

	if got != "world" {
		t.Errorf("got %v", got)
	}
}

func TestTranslatorParseKeyWithNamespace(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "en")
	ns, group, item := tr.ParseKey("vendor::messages.greeting")

	if ns != "vendor" || group != "messages" || item != "greeting" {
		t.Errorf("ParseKey = (%q, %q, %q)", ns, group, item)
	}
}

func TestTranslatorParseKeyWithoutNamespace(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "en")
	ns, group, item := tr.ParseKey("messages.greeting")

	if ns != "*" || group != "messages" || item != "greeting" {
		t.Errorf("ParseKey = (%q, %q, %q)", ns, group, item)
	}
}

func TestTranslatorParseKeyFlatJson(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "en")
	ns, group, item := tr.ParseKey("Hello World")

	if ns != "*" || group != "*" || item != "Hello World" {
		t.Errorf("ParseKey = (%q, %q, %q)", ns, group, item)
	}
}

func TestTranslatorMakeReplacementsAtomicSubstitution(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "en")

	got := tr.MakeReplacements("Hello :foo!", map[string]any{
		"foo": "baz:bar",
		"bar": "WRONG",
	})

	if got != "Hello baz:bar!" {
		t.Errorf("got %q; want %q", got, "Hello baz:bar!")
	}
}

func TestTranslatorMakeReplacementsCapitalization(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "en")

	got := tr.MakeReplacements(":name :Name :NAME", map[string]any{"name": "alice"})

	if got != "alice Alice ALICE" {
		t.Errorf("got %q; want %q", got, "alice Alice ALICE")
	}
}

func TestTranslatorChoiceWithSlice(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "en")
	tr.AddLines(map[string]any{"items": "{1} one|[2,*] many"}, "en")

	got := tr.Choice("items", []int{1, 2, 3}, nil, nil)

	if got != "many" {
		t.Errorf("got %q", got)
	}
}

func TestTranslatorPotentiallyTranslatedString(t *testing.T) {
	t.Parallel()

	tr := translation.NewTranslator(translation.NewArrayLoader(), "en")
	tr.AddLines(map[string]any{"welcome": "Welcome :name"}, "en")

	pts := translation.NewPotentiallyTranslatedString("welcome", tr)
	pts.Translate(map[string]any{"name": "Alice"}, nil)

	if pts.String() != "Welcome Alice" {
		t.Errorf("String() = %q", pts.String())
	}

	if pts.Original() != "welcome" {
		t.Errorf("Original() = %q", pts.Original())
	}
}
