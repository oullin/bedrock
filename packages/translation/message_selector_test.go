package translation_test

import (
	"testing"

	"github.com/bedrock/packages/translation"
)

func TestMessageSelectorChooseSingleSegment(t *testing.T) {
	t.Parallel()

	s := translation.NewMessageSelector()
	got := s.Choose("hello", 1, "en")

	if got != "hello" {
		t.Errorf("got %q", got)
	}
}

func TestMessageSelectorChooseEmptyLine(t *testing.T) {
	t.Parallel()

	s := translation.NewMessageSelector()
	got := s.Choose("", 1, "en")

	if got != "" {
		t.Errorf("got %q", got)
	}
}

func TestMessageSelectorChoosePipeOnly(t *testing.T) {
	t.Parallel()

	s := translation.NewMessageSelector()
	got := s.Choose("|", 1, "en")

	if got != "" {
		t.Errorf("got %q", got)
	}
}

func TestMessageSelectorGetPluralIndexEnglish(t *testing.T) {
	t.Parallel()

	s := translation.NewMessageSelector()

	if s.Choose("one|many", 1, "en") != "one" {
		t.Error("English n=1 should select first form")
	}

	if s.Choose("one|many", 0, "en") != "many" {
		t.Error("English n=0 should select second form")
	}

	if s.Choose("one|many", 2, "en") != "many" {
		t.Error("English n=2 should select second form")
	}
}

func TestMessageSelectorGetPluralIndexRussian(t *testing.T) {
	t.Parallel()

	s := translation.NewMessageSelector()

	// Russian: 1→0, 2→1, 5→2, 11→2, 21→0.
	cases := []struct {
		n    float64
		want string
	}{
		{1, "один"},
		{2, "два"},
		{5, "пять"},
		{11, "пять"},
		{21, "один"},
	}

	for _, c := range cases {
		got := s.Choose("один|два|пять", c.n, "ru")

		if got != c.want {
			t.Errorf("Russian n=%v: got %q; want %q", c.n, got, c.want)
		}
	}
}

func TestMessageSelectorGetPluralIndexJapanese(t *testing.T) {
	t.Parallel()

	s := translation.NewMessageSelector()

	// Japanese always returns form 0.
	for _, n := range []float64{0, 1, 2, 100} {
		got := s.Choose("個", n, "ja")

		if got != "個" {
			t.Errorf("Japanese n=%v: got %q", n, got)
		}
	}
}
