package compiler

import (
	"regexp"
	"testing"
)

// fakeRoute is a minimal SourceRoute used by the compiler tests; the real
// [routing.Route] type satisfies the same interface.
type fakeRoute struct {
	path         string
	host         string
	defaults     map[string]any
	requirements map[string]string
}

func (f fakeRoute) Path() string                    { return f.path }
func (f fakeRoute) Host() string                    { return f.host }
func (f fakeRoute) Requirements() map[string]string { return f.requirements }
func (f fakeRoute) HasDefault(name string) bool {
	_, ok := f.defaults[name]

	return ok
}

func mustCompile(t *testing.T, r SourceRoute) *CompiledRoute {
	t.Helper()
	c, err := Compile(r)

	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	return c
}

func TestCompile_StaticPath(t *testing.T) {
	c := mustCompile(t, fakeRoute{path: "/foo"})

	if c.StaticPrefix() != "/foo" {
		t.Errorf("static prefix = %q, want /foo", c.StaticPrefix())
	}

	if !regexp.MustCompile(c.Regex()).MatchString("/foo") {
		t.Errorf("regex %q does not match /foo", c.Regex())
	}
}

func TestCompile_SinglePathVariable(t *testing.T) {
	c := mustCompile(t, fakeRoute{path: "/users/{user}"})

	if got, want := c.PathVariables(), []string{"user"}; !equal(got, want) {
		t.Errorf("variables = %v, want %v", got, want)
	}

	re := regexp.MustCompile(c.Regex())
	m := re.FindStringSubmatch("/users/42")

	if m == nil {
		t.Fatalf("regex %q did not match /users/42", c.Regex())
	}

	if got := m[re.SubexpIndex("user")]; got != "42" {
		t.Errorf("captured user = %q, want 42", got)
	}
}

func TestCompile_OptionalVariableWithDefault(t *testing.T) {
	c := mustCompile(t, fakeRoute{
		path:     "/users/{user}",
		defaults: map[string]any{"user": "guest"},
	})
	re := regexp.MustCompile(c.Regex())

	if !re.MatchString("/users") {
		t.Errorf("regex %q should match /users (optional segment)", c.Regex())
	}

	if !re.MatchString("/users/alice") {
		t.Errorf("regex %q should match /users/alice", c.Regex())
	}
}

func TestCompile_VariableRequirement(t *testing.T) {
	c := mustCompile(t, fakeRoute{
		path:         "/users/{id}",
		requirements: map[string]string{"id": "[0-9]+"},
	})
	re := regexp.MustCompile(c.Regex())

	if !re.MatchString("/users/123") {
		t.Errorf("should match /users/123")
	}

	if re.MatchString("/users/abc") {
		t.Errorf("should not match /users/abc")
	}
}

func TestCompile_HostPattern(t *testing.T) {
	c := mustCompile(t, fakeRoute{path: "/", host: "{sub}.example.com"})

	if c.HostRegex() == "" {
		t.Fatal("host regex should not be empty")
	}

	re := regexp.MustCompile(c.HostRegex())
	m := re.FindStringSubmatch("api.example.com")

	if m == nil {
		t.Fatalf("host regex %q did not match api.example.com", c.HostRegex())
	}

	if got := m[re.SubexpIndex("sub")]; got != "api" {
		t.Errorf("captured sub = %q, want api", got)
	}
}

func TestCompile_DuplicateVariableErrors(t *testing.T) {
	if _, err := Compile(fakeRoute{path: "/{x}/{x}"}); err == nil {
		t.Error("expected duplicate variable error")
	}
}

func TestCompile_DigitLeadingNameErrors(t *testing.T) {
	if _, err := Compile(fakeRoute{path: "/{1bad}"}); err == nil {
		t.Error("expected digit-leading variable error")
	}
}

func TestCompile_VariableTooLongErrors(t *testing.T) {
	if _, err := Compile(fakeRoute{path: "/{aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa}"}); err == nil {
		t.Error("expected too-long variable error")
	}
}

func TestCompile_FragmentReservedErrors(t *testing.T) {
	if _, err := Compile(fakeRoute{path: "/{_fragment}"}); err == nil {
		t.Error("expected _fragment reserved error")
	}
}

func TestCompile_TwoVariables(t *testing.T) {
	c := mustCompile(t, fakeRoute{path: "/posts/{post}/comments/{comment}"})
	re := regexp.MustCompile(c.Regex())
	m := re.FindStringSubmatch("/posts/1/comments/2")

	if m == nil {
		t.Fatalf("regex %q did not match", c.Regex())
	}

	if m[re.SubexpIndex("post")] != "1" || m[re.SubexpIndex("comment")] != "2" {
		t.Errorf("captures wrong: %v", m)
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
