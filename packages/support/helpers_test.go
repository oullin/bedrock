package support

import (
	"os"
	"testing"
)

// Ref: @bedrock/code-0377
func TestBlank(t *testing.T) {
	t.Parallel()

	cases := []struct {
		value  any
		expect bool
	}{
		{nil, true},
		{"", true},
		{"  ", true},
		{"\t\n", true},
		{"hello", false},
		{" hi ", false},
		{[]any{}, true},
		{[]int{}, true},
		{[]int{1}, false},
		{map[string]any{}, true},
		{map[string]any{"a": 1}, false},
	}

	for _, tc := range cases {
		got := Blank(tc.value)

		if got != tc.expect {
			t.Errorf("Blank(%v) = %v, want %v", tc.value, got, tc.expect)
		}
	}
}

// Ref: @bedrock/code-0377
func TestFilled(t *testing.T) {
	t.Parallel()

	if Filled(nil) {
		t.Error("Filled(nil) should be false")
	}

	if Filled("") {
		t.Error("Filled(\"\") should be false")
	}

	if !Filled("hello") {
		t.Error("Filled(\"hello\") should be true")
	}

	if !Filled([]int{1}) {
		t.Error("Filled([]int{1}) should be true")
	}

	if Filled([]int{}) {
		t.Error("Filled([]int{}) should be false")
	}
}

// Ref: @bedrock/code-0377
func TestTap(t *testing.T) {
	t.Parallel()

	called := false
	result := Tap("hello", func(v string) {
		called = true

		if v != "hello" {
			t.Errorf("expected 'hello', got %q", v)
		}
	})

	if !called {
		t.Error("callback was not called")
	}

	if result != "hello" {
		t.Errorf("Tap should return original value, got %q", result)
	}
}

// Ref: @bedrock/code-0377
func TestWith(t *testing.T) {
	t.Parallel()

	// With without callback returns value
	r1 := With("hello")

	if r1 != "hello" {
		t.Errorf("expected 'hello', got %q", r1)
	}

	// With callback transforms value
	r2 := With("hello", func(s string) string {
		return s + " world"
	})

	if r2 != "hello world" {
		t.Errorf("expected 'hello world', got %q", r2)
	}
}

// Ref: @bedrock/code-0377
func TestValue(t *testing.T) {
	t.Parallel()

	if got := Value[string]("literal"); got != "literal" {
		t.Errorf("Value(literal) = %q", got)
	}

	if got := Value[string](func() string { return "callback" }); got != "callback" {
		t.Errorf("Value(callback) = %q", got)
	}

	if got := Value[string](func(value any) string { return value.(string) + " value" }, "passed"); got != "passed value" {
		t.Errorf("Value(callback with arg) = %q", got)
	}
}

// Ports of:
// Ref: @bedrock/code-0377
func TestTransform(t *testing.T) {
	t.Parallel()

	// Non-blank value gets transformed
	result, ok := Transform("hello", func(s string) string {
		return s + " world"
	})

	if !ok || result != "hello world" {
		t.Errorf("expected 'hello world', ok=true; got %q, ok=%v", result, ok)
	}

	// Blank value does not get transformed
	result2, ok2 := Transform("", func(s string) string {
		return "should not run"
	})

	if ok2 || result2 != "" {
		t.Errorf("expected blank result, got %q, ok=%v", result2, ok2)
	}

	// Blank with default
	result3, _ := Transform("", func(s string) string {
		return "nope"
	}, "default")

	if result3 != "default" {
		t.Errorf("expected 'default', got %q", result3)
	}
}

// Ref: @bedrock/code-0377
func TestE(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input  string
		expect string
	}{
		{"<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		{"Hello & World", "Hello &amp; World"},
		{"\"quoted\"", "&#34;quoted&#34;"},
		{"plain text", "plain text"},
	}

	for _, tc := range cases {
		got := E(tc.input)

		if got != tc.expect {
			t.Errorf("E(%q) = %q, want %q", tc.input, got, tc.expect)
		}
	}
}

// Ref: @bedrock/code-0377
func TestEnv(t *testing.T) {
	// NOT parallel — uses t.Setenv
	t.Setenv("TEST_SUPPORT_FOO", "bar")

	if got := Env("TEST_SUPPORT_FOO"); got != "bar" {
		t.Errorf("Env(set) = %q, want %q", got, "bar")
	}

	if got := Env("TEST_SUPPORT_NOTSET", "default"); got != "default" {
		t.Errorf("Env(unset) = %q, want 'default'", got)
	}
}

// Ref: @bedrock/code-0377
func TestEnvTrue(t *testing.T) {
	t.Parallel()

	os.Setenv("TEST_SUPPORT_BOOL_TRUE", "true")

	defer os.Unsetenv("TEST_SUPPORT_BOOL_TRUE")

	if got := Env("TEST_SUPPORT_BOOL_TRUE"); got != "true" {
		t.Errorf("expected 'true', got %q", got)
	}
}

// Ref: @bedrock/code-0377
func TestEnvFalse(t *testing.T) {
	t.Parallel()

	os.Setenv("TEST_SUPPORT_BOOL_FALSE", "false")

	defer os.Unsetenv("TEST_SUPPORT_BOOL_FALSE")

	if got := Env("TEST_SUPPORT_BOOL_FALSE"); got != "false" {
		t.Errorf("expected 'false', got %q", got)
	}
}

// Ref: @bedrock/code-0377
func TestEnvNull(t *testing.T) {
	t.Parallel()

	os.Setenv("TEST_SUPPORT_NULL", "null")

	defer os.Unsetenv("TEST_SUPPORT_NULL")

	if got := Env("TEST_SUPPORT_NULL"); got != "" {
		t.Errorf("expected empty string for null, got %q", got)
	}
}

// Ref: @bedrock/code-0377
func TestEnvEmpty(t *testing.T) {
	t.Parallel()

	os.Setenv("TEST_SUPPORT_EMPTY", "empty")

	defer os.Unsetenv("TEST_SUPPORT_EMPTY")

	if got := Env("TEST_SUPPORT_EMPTY"); got != "" {
		t.Errorf("expected empty string for 'empty', got %q", got)
	}
}

// Ref: @bedrock/code-0377
func TestEnvEscapedString(t *testing.T) {
	t.Parallel()

	os.Setenv("TEST_SUPPORT_QUOTED", "\"hello world\"")

	defer os.Unsetenv("TEST_SUPPORT_QUOTED")

	if got := Env("TEST_SUPPORT_QUOTED"); got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestInventoryHeadLastAndClassBasename(t *testing.T) {
	t.Parallel()

	// SupportHelpersTest::testHead
	// SupportHelpersTest::testLast
	// SupportHelpersTest::testClassBasename
	head, ok := Head([]string{"first", "second"})

	if !ok || head != "first" {
		t.Fatalf("Head = %q, %v", head, ok)
	}

	last, ok := Last([]string{"first", "second"})

	if !ok || last != "second" {
		t.Fatalf("Last = %q, %v", last, ok)
	}

	if _, ok := Head([]string{}); ok {
		t.Fatal("Head on empty slice should report false")
	}

	if _, ok := Last([]string{}); ok {
		t.Fatal("Last on empty slice should report false")
	}

	type localHelperType struct{}

	if got := ClassBasename("App\\Models\\User"); got != "User" {
		t.Fatalf("ClassBasename PHP class = %q", got)
	}

	if got := ClassBasename("/app/Models/User"); got != "User" {
		t.Fatalf("ClassBasename path = %q", got)
	}

	if got := ClassBasename(&localHelperType{}); got != "localHelperType" {
		t.Fatalf("ClassBasename type = %q", got)
	}
}
