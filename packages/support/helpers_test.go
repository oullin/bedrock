package support

import (
	"os"
	"testing"
)

// Port of Illuminate\Tests\Support\SupportHelpersTest::testBlank
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

// Port of Illuminate\Tests\Support\SupportHelpersTest::testFilled
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

// Port of Illuminate\Tests\Support\SupportHelpersTest::testTap
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

// Port of Illuminate\Tests\Support\SupportHelpersTest::testWith
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

// Port of Illuminate\Tests\Support\SupportHelpersTest::testTransform
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

// Port of Illuminate\Tests\Support\SupportHelpersTest::testE
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

// Port of Illuminate\Tests\Support\SupportHelpersTest::testEnv
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

// Port of Illuminate\Tests\Support\SupportHelpersTest::testEnvTrue
func TestEnvTrue(t *testing.T) {
	t.Parallel()

	os.Setenv("TEST_SUPPORT_BOOL_TRUE", "true")
	defer os.Unsetenv("TEST_SUPPORT_BOOL_TRUE")

	if got := Env("TEST_SUPPORT_BOOL_TRUE"); got != "true" {
		t.Errorf("expected 'true', got %q", got)
	}
}

// Port of Illuminate\Tests\Support\SupportHelpersTest::testEnvFalse
func TestEnvFalse(t *testing.T) {
	t.Parallel()

	os.Setenv("TEST_SUPPORT_BOOL_FALSE", "false")
	defer os.Unsetenv("TEST_SUPPORT_BOOL_FALSE")

	if got := Env("TEST_SUPPORT_BOOL_FALSE"); got != "false" {
		t.Errorf("expected 'false', got %q", got)
	}
}

// Port of Illuminate\Tests\Support\SupportHelpersTest::testEnvNull
func TestEnvNull(t *testing.T) {
	t.Parallel()

	os.Setenv("TEST_SUPPORT_NULL", "null")
	defer os.Unsetenv("TEST_SUPPORT_NULL")

	if got := Env("TEST_SUPPORT_NULL"); got != "" {
		t.Errorf("expected empty string for null, got %q", got)
	}
}

// Port of Illuminate\Tests\Support\SupportHelpersTest::testEnvEmpty
func TestEnvEmpty(t *testing.T) {
	t.Parallel()

	os.Setenv("TEST_SUPPORT_EMPTY", "empty")
	defer os.Unsetenv("TEST_SUPPORT_EMPTY")

	if got := Env("TEST_SUPPORT_EMPTY"); got != "" {
		t.Errorf("expected empty string for 'empty', got %q", got)
	}
}

// Port of Illuminate\Tests\Support\SupportHelpersTest::testEnvEscapedString
func TestEnvEscapedString(t *testing.T) {
	t.Parallel()

	os.Setenv("TEST_SUPPORT_QUOTED", "\"hello world\"")
	defer os.Unsetenv("TEST_SUPPORT_QUOTED")

	if got := Env("TEST_SUPPORT_QUOTED"); got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}
