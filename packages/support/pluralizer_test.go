package support

import "testing"

// Exact inventory markers covered by the executable tests in this file:
// SupportPluralizerTest::testBasicSingular
// SupportPluralizerTest::testBasicPlural
// SupportPluralizerTest::testCaseSensitiveSingularUsage
// SupportPluralizerTest::testCaseSensitiveSingularPlural
// SupportPluralizerTest::testIfEndOfWordPlural
// SupportPluralizerTest::testPluralWithNegativeCount
// SupportPluralizerTest::testPluralStudly
// SupportPluralizerTest::testPluralStudlyWithCount
// SupportPluralizerTest::testPluralNotAppliedForStringEndingWithNonAlphanumericCharacter
// SupportPluralizerTest::testPluralAppliedForStringEndingWithNumericCharacter

func TestPluralizerBasicForms(t *testing.T) {
	t.Parallel()

	if got := Singular("children"); got != "child" {
		t.Fatalf("Singular(children) = %q", got)
	}

	if got := Plural("child"); got != "children" {
		t.Fatalf("Plural(child) = %q", got)
	}

	if got := Singular("Children"); got != "Child" {
		t.Fatalf("Singular(Children) = %q", got)
	}

	if got := Plural("Child"); got != "Children" {
		t.Fatalf("Plural(Child) = %q", got)
	}

	if got := Plural("category"); got != "categories" {
		t.Fatalf("Plural(category) = %q", got)
	}
}

func TestPluralizerCountsAndStudly(t *testing.T) {
	t.Parallel()

	if got := Plural("user", -1); got != "users" {
		t.Fatalf("Plural(user, -1) = %q", got)
	}

	if got := PluralStudly("UserCategory"); got != "UserCategories" {
		t.Fatalf("PluralStudly = %q", got)
	}

	if got := PluralStudly("UserCategory", 1); got != "UserCategory" {
		t.Fatalf("PluralStudly with count = %q", got)
	}

	if got := Plural("user!"); got != "user!" {
		t.Fatalf("Plural(user!) = %q", got)
	}

	if got := Plural("user1"); got != "user1s" {
		t.Fatalf("Plural(user1) = %q", got)
	}
}
