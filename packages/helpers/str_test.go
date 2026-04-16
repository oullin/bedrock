package helpers

import (
	"strings"
	"testing"
)

// Port of Laravel\Helpers\Tests\HelpersTest::test_autoloading
func TestStrContainsDelegation(t *testing.T) {
	t.Parallel()

	if !StrContains("Laravel", "Lara") {
		t.Error("StrContains should find 'Lara' in 'Laravel'")
	}
}

func TestCamelCase(t *testing.T) {
	t.Parallel()

	if got := CamelCase("hello_world"); got != "helloWorld" {
		t.Errorf("CamelCase = %q, want %q", got, "helloWorld")
	}
}

func TestEndsWith(t *testing.T) {
	t.Parallel()

	if !EndsWith("jason", "on") {
		t.Error("EndsWith(on) should be true")
	}
	if EndsWith("jason", "no") {
		t.Error("EndsWith(no) should be false")
	}
}

func TestKebabCase(t *testing.T) {
	t.Parallel()

	if got := KebabCase("helloWorld"); got != "hello-world" {
		t.Errorf("KebabCase = %q, want %q", got, "hello-world")
	}
}

func TestSnakeCase(t *testing.T) {
	t.Parallel()

	if got := SnakeCase("helloWorld"); got != "hello_world" {
		t.Errorf("SnakeCase = %q, want %q", got, "hello_world")
	}
}

func TestStartsWith(t *testing.T) {
	t.Parallel()

	if !StartsWith("jason", "jas") {
		t.Error("StartsWith(jas) should be true")
	}
	if StartsWith("jason", "day") {
		t.Error("StartsWith(day) should be false")
	}
}

func TestStrAfterDelegation(t *testing.T) {
	t.Parallel()

	if got := StrAfter("hannah", "han"); got != "nah" {
		t.Errorf("StrAfter = %q, want %q", got, "nah")
	}
}

func TestStrBeforeDelegation(t *testing.T) {
	t.Parallel()

	if got := StrBefore("hannah", "nah"); got != "han" {
		t.Errorf("StrBefore = %q, want %q", got, "han")
	}
}

func TestStrFinishDelegation(t *testing.T) {
	t.Parallel()

	if got := StrFinish("hello", "/"); got != "hello/" {
		t.Errorf("StrFinish = %q, want %q", got, "hello/")
	}
	if got := StrFinish("hello/", "/"); got != "hello/" {
		t.Errorf("StrFinish already finished = %q, want %q", got, "hello/")
	}
}

func TestStrIsDelegation(t *testing.T) {
	t.Parallel()

	if !StrIs("foo*", "foobar") {
		t.Error("StrIs(foo*, foobar) should be true")
	}
	if StrIs("baz*", "foobar") {
		t.Error("StrIs(baz*, foobar) should be false")
	}
}

func TestStrLimitDelegation(t *testing.T) {
	t.Parallel()

	if got := StrLimit("The PHP framework for web artisans.", 7); got != "The PHP..." {
		t.Errorf("StrLimit = %q, want %q", got, "The PHP...")
	}
}

func TestStrPluralDelegation(t *testing.T) {
	t.Parallel()

	if got := StrPlural("car"); got != "cars" {
		t.Errorf("StrPlural = %q, want %q", got, "cars")
	}
	if got := StrPlural("car", 1); got != "car" {
		t.Errorf("StrPlural(1) = %q, want %q", got, "car")
	}
}

func TestStrRandomDelegation(t *testing.T) {
	t.Parallel()

	got := StrRandom(16)
	if len(got) == 0 {
		t.Error("StrRandom should return non-empty string")
	}
}

func TestStrReplaceArrayDelegation(t *testing.T) {
	t.Parallel()

	got := StrReplaceArray("?", []string{"8:30", "9:00"}, "? and ?")
	if got != "8:30 and 9:00" {
		t.Errorf("StrReplaceArray = %q, want %q", got, "8:30 and 9:00")
	}
}

func TestStrReplaceFirstDelegation(t *testing.T) {
	t.Parallel()

	got := StrReplaceFirst("the", "a", "the quick brown fox jumps over the lazy dog")
	if !strings.HasPrefix(got, "a quick") {
		t.Errorf("StrReplaceFirst = %q", got)
	}
}

func TestStrReplaceLastDelegation(t *testing.T) {
	t.Parallel()

	got := StrReplaceLast("the", "a", "the quick brown fox jumps over the lazy dog")
	if !strings.Contains(got, "a lazy") {
		t.Errorf("StrReplaceLast = %q", got)
	}
}

func TestStrSingularDelegation(t *testing.T) {
	t.Parallel()

	if got := StrSingular("cars"); got != "car" {
		t.Errorf("StrSingular = %q, want %q", got, "car")
	}
}

func TestStrSlugDelegation(t *testing.T) {
	t.Parallel()

	if got := StrSlug("Hello World"); got != "hello-world" {
		t.Errorf("StrSlug = %q, want %q", got, "hello-world")
	}
}

func TestStrStartDelegation(t *testing.T) {
	t.Parallel()

	if got := StrStart("path", "/"); got != "/path" {
		t.Errorf("StrStart = %q, want %q", got, "/path")
	}
	if got := StrStart("/path", "/"); got != "/path" {
		t.Errorf("StrStart already started = %q, want %q", got, "/path")
	}
}

func TestStudlyCase(t *testing.T) {
	t.Parallel()

	if got := StudlyCase("hello_world"); got != "HelloWorld" {
		t.Errorf("StudlyCase = %q, want %q", got, "HelloWorld")
	}
}

func TestTitleCase(t *testing.T) {
	t.Parallel()

	if got := TitleCase("hello world"); got != "Hello World" {
		t.Errorf("TitleCase = %q, want %q", got, "Hello World")
	}
}
