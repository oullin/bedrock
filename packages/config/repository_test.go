package config

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRepositoryLookupAndTypedAccessors(t *testing.T) {
	t.Parallel()

	repo := NewRepository(map[string]any{
		"foo":     "bar",
		"boolean": true,
		"integer": 1,
		"float":   1.25,
		"null":    nil,
		"auth": map[string]any{
			"session_lifetime": "24h",
			"session_duration": 2 * time.Hour,
			"cookies": map[string]any{
				"secure": true,
			},
			"middleware": []any{"web", "auth"},
			"providers":  []string{"email", "sms"},
			"count_text": "7",
			"flags":      "true",
		},
		"mail.mailers.smtp": "literal",
		"mail": map[string]any{
			"mailers": map[string]any{
				"smtp": "nested",
			},
		},
		"a": map[string]any{
			"b.c": "shadowed",
			"b": map[string]any{
				"c": "nested",
			},
		},
		"x": map[string]any{
			"z": "zoo",
		},
	})

	if !repo.Has("auth.session_lifetime") {
		t.Fatal("expected nested key to exist")
	}

	if repo.Has("auth.missing") {
		t.Fatal("expected missing nested key to be absent")
	}

	if got := repo.Get("foo", nil); got != "bar" {
		t.Fatalf("unexpected direct value: %#v", got)
	}

	if got := repo.Get("auth.missing", "fallback"); got != "fallback" {
		t.Fatalf("unexpected fallback value: %#v", got)
	}

	if got := repo.Get("mail.mailers.smtp", nil); got != "literal" {
		t.Fatalf("unexpected dotted key precedence: %#v", got)
	}

	if got := repo.Get("a.b.c", nil); got != "shadowed" {
		t.Fatalf("unexpected nested dotted key precedence: %#v", got)
	}

	if got := repo.Get("x.y.z", nil); got != nil {
		t.Fatalf("expected nil for missing deep key, got %#v", got)
	}

	if got := repo.Get(".", nil); got != nil {
		t.Fatalf("expected nil for invalid dot path, got %#v", got)
	}

	if got := repo.Get("null", "fallback"); got != nil {
		t.Fatalf("expected nil value to be preserved, got %#v", got)
	}

	if got := repo.Get("", nil); !reflect.DeepEqual(got, repo.All()) {
		t.Fatalf("expected empty key to return all config, got %#v", got)
	}

	values := repo.GetMany(map[string]any{
		"foo": "default",
		"x.y": "default",
		"x.z": "default",
		"baz": nil,
	})
	expectedValues := map[string]any{
		"foo": "bar",
		"x.y": "default",
		"x.z": "zoo",
		"baz": nil,
	}

	if !reflect.DeepEqual(values, expectedValues) {
		t.Fatalf("unexpected GetMany result: %#v", values)
	}

	if got := repo.GetMany(nil); len(got) != 0 {
		t.Fatalf("expected empty GetMany result, got %#v", got)
	}

	duration, err := repo.Duration("auth.session_lifetime")

	if err != nil {
		t.Fatalf("Duration: %v", err)
	}

	if duration != 24*time.Hour {
		t.Fatalf("unexpected duration: %v", duration)
	}

	duration, err = repo.Duration("auth.session_duration")

	if err != nil {
		t.Fatalf("Duration direct: %v", err)
	}

	if duration != 2*time.Hour {
		t.Fatalf("unexpected direct duration: %v", duration)
	}

	secure, err := repo.Bool("auth.cookies.secure")

	if err != nil {
		t.Fatalf("Bool: %v", err)
	}

	if !secure {
		t.Fatal("expected secure cookie")
	}

	boolFromString, err := repo.Bool("auth.flags")

	if err != nil {
		t.Fatalf("Bool string coercion: %v", err)
	}

	if !boolFromString {
		t.Fatal("expected bool string coercion to succeed")
	}

	number, err := repo.Int("integer")

	if err != nil {
		t.Fatalf("Int direct: %v", err)
	}

	if number != 1 {
		t.Fatalf("unexpected integer: %d", number)
	}

	number, err = repo.Int("auth.count_text")

	if err != nil {
		t.Fatalf("Int string coercion: %v", err)
	}

	if number != 7 {
		t.Fatalf("unexpected coerced integer: %d", number)
	}

	number, err = repo.Int("float")

	if err != nil {
		t.Fatalf("Int float coercion: %v", err)
	}

	if number != 1 {
		t.Fatalf("unexpected float-to-int coercion: %d", number)
	}

	middleware, err := repo.StringSlice("auth.middleware")

	if err != nil {
		t.Fatalf("StringSlice: %v", err)
	}

	if !reflect.DeepEqual(middleware, []string{"web", "auth"}) {
		t.Fatalf("unexpected middleware: %#v", middleware)
	}

	providers, err := repo.StringSlice("auth.providers")

	if err != nil {
		t.Fatalf("StringSlice []string: %v", err)
	}

	if !reflect.DeepEqual(providers, []string{"email", "sms"}) {
		t.Fatalf("unexpected []string providers: %#v", providers)
	}

	if got, err := repo.String("foo"); err != nil || got != "bar" {
		t.Fatalf("String: got=%q err=%v", got, err)
	}
}

func TestRepositorySetMutatorsAndCloneSemantics(t *testing.T) {
	t.Parallel()

	repo := NewRepository(nil)
	repo.Set("fortify.limiters.login", "login")
	repo.Set("fortify.nil-value")
	repo.Set(map[string]any{
		"fortify.limiters.two-factor":                       "two-factor",
		"fortify.options.two-factor-authentication.confirm": true,
		"fortify.alt-middleware":                            []string{"api", "signed"},
	})
	repo.Set("fortify.limiters", "replaced")
	repo.Set("fortify.limiters.login", "login")
	repo.Push("fortify.middleware", "web")
	repo.Push("fortify.middleware", "guest")
	repo.Prepend("fortify.middleware", "trim")
	repo.Prepend("fortify.new-middleware", "first")
	repo.Push("fortify.new-middleware", "second")
	repo.Push("fortify.alt-middleware", "verified")

	mapped, err := repo.Map("fortify.options")

	if err != nil {
		t.Fatalf("Map: %v", err)
	}

	inner, ok := mapped["two-factor-authentication"].(map[string]any)

	if !ok || inner["confirm"] != true {
		t.Fatalf("unexpected nested option map: %#v", mapped)
	}

	middleware, err := repo.StringSlice("fortify.middleware")

	if err != nil {
		t.Fatalf("StringSlice middleware: %v", err)
	}

	if !reflect.DeepEqual(middleware, []string{"trim", "web", "guest"}) {
		t.Fatalf("unexpected middleware order: %#v", middleware)
	}

	altMiddleware, err := repo.StringSlice("fortify.alt-middleware")

	if err != nil {
		t.Fatalf("StringSlice alt middleware: %v", err)
	}

	if !reflect.DeepEqual(altMiddleware, []string{"api", "signed", "verified"}) {
		t.Fatalf("unexpected alt middleware: %#v", altMiddleware)
	}

	newMiddleware, err := repo.StringSlice("fortify.new-middleware")

	if err != nil {
		t.Fatalf("StringSlice new middleware: %v", err)
	}

	if !reflect.DeepEqual(newMiddleware, []string{"first", "second"}) {
		t.Fatalf("unexpected new middleware: %#v", newMiddleware)
	}

	if got := repo.Get("fortify.nil-value", "fallback"); got != nil {
		t.Fatalf("expected nil value from single-arg Set, got %#v", got)
	}

	all := repo.All()
	fortify := all["fortify"].(map[string]any)
	mutatedOptions := fortify["options"].(map[string]any)
	mutatedOptions["added"] = "mutated"
	fortify["middleware"] = []any{"mutated"}

	options, err := repo.Map("fortify.options")

	if err != nil {
		t.Fatalf("Map after clone mutation: %v", err)
	}

	if _, ok := options["added"]; ok {
		t.Fatal("repository mutated through All clone")
	}

	currentMiddleware, err := repo.StringSlice("fortify.middleware")

	if err != nil {
		t.Fatalf("StringSlice after clone mutation: %v", err)
	}

	if !reflect.DeepEqual(currentMiddleware, []string{"trim", "web", "guest"}) {
		t.Fatalf("repository middleware mutated through All clone: %#v", currentMiddleware)
	}

	loginValue, err := repo.String("fortify.limiters.login")

	if err != nil {
		t.Fatalf("String limiter: %v", err)
	}

	if loginValue != "login" {
		t.Fatalf("unexpected login limiter: %q", loginValue)
	}

	limits, err := repo.Map("fortify.limiters")

	if err != nil {
		t.Fatalf("Map limiters: %v", err)
	}

	limits["login"] = "mutated"
	loginValue, err = repo.String("fortify.limiters.login")

	if err != nil {
		t.Fatalf("String limiter after map mutation: %v", err)
	}

	if loginValue != "login" {
		t.Fatalf("repository mutated through Map clone: %q", loginValue)
	}
}

func TestRepositoryTypeFailuresAndMissingKeys(t *testing.T) {
	t.Parallel()

	repo := NewRepository(map[string]any{
		"fortify": map[string]any{
			"views":            "not-a-bool",
			"features":         []any{"registration", 42},
			"middleware":       "web",
			"not-string":       true,
			"not-int":          "abc",
			"not-int-default":  true,
			"not-bool":         1,
			"not-duration":     "tomorrow",
			"not-duration-alt": true,
			"not-map":          []any{"web"},
			"not-slice":        true,
			"int64-value":      int64(9),
			"string-list-text": "web, auth , signed",
			"empty-list-text":  "   ",
		},
	})

	stringList, err := repo.StringSlice("fortify.string-list-text")

	if err != nil {
		t.Fatalf("StringSlice csv: %v", err)
	}

	if !reflect.DeepEqual(stringList, []string{"web", "auth", "signed"}) {
		t.Fatalf("unexpected csv slice: %#v", stringList)
	}

	emptyList, err := repo.StringSlice("fortify.empty-list-text")

	if err != nil {
		t.Fatalf("StringSlice empty csv: %v", err)
	}

	if len(emptyList) != 0 {
		t.Fatalf("expected empty slice, got %#v", emptyList)
	}

	int64Value, err := repo.Int("fortify.int64-value")

	if err != nil {
		t.Fatalf("Int int64 coercion: %v", err)
	}

	if int64Value != 9 {
		t.Fatalf("unexpected int64 coercion result: %d", int64Value)
	}

	for _, tc := range []struct {
		name string
		fn   func() error
		want string
	}{

		{name: "bool type", fn: func() error { _, err := repo.Bool("fortify.views"); return err }, want: `config: key "fortify.views" must be bool, got string`},

		{name: "string type", fn: func() error { _, err := repo.String("fortify.not-string"); return err }, want: `config: key "fortify.not-string" must be string, got bool`},

		{name: "int type", fn: func() error { _, err := repo.Int("fortify.not-int"); return err }, want: `config: key "fortify.not-int" must be int, got string`},

		{name: "int default type", fn: func() error { _, err := repo.Int("fortify.not-int-default"); return err }, want: `config: key "fortify.not-int-default" must be int, got bool`},

		{name: "bool default type", fn: func() error { _, err := repo.Bool("fortify.not-bool"); return err }, want: `config: key "fortify.not-bool" must be bool, got int`},

		{name: "duration type", fn: func() error { _, err := repo.Duration("fortify.not-duration"); return err }, want: `config: key "fortify.not-duration" must be duration, got string`},

		{name: "duration default type", fn: func() error { _, err := repo.Duration("fortify.not-duration-alt"); return err }, want: `config: key "fortify.not-duration-alt" must be duration, got bool`},

		{name: "slice type", fn: func() error { _, err := repo.StringSlice("fortify.features"); return err }, want: `config: key "fortify.features" must be []string, got []interface {}`},

		{name: "slice default type", fn: func() error { _, err := repo.StringSlice("fortify.not-slice"); return err }, want: `config: key "fortify.not-slice" must be []string, got bool`},

		{name: "map type", fn: func() error { _, err := repo.Map("fortify.not-map"); return err }, want: `config: key "fortify.not-map" must be map[string]any, got []interface {}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn()

			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != tc.want {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}

	for _, tc := range []struct {
		name string
		fn   func() error
		want string
	}{

		{name: "missing string", fn: func() error { _, err := repo.String("missing"); return err }, want: `config: missing key "missing"`},

		{name: "missing int", fn: func() error { _, err := repo.Int("missing"); return err }, want: `config: missing key "missing"`},

		{name: "missing bool", fn: func() error { _, err := repo.Bool("missing"); return err }, want: `config: missing key "missing"`},

		{name: "missing duration", fn: func() error { _, err := repo.Duration("missing"); return err }, want: `config: missing key "missing"`},

		{name: "missing slice", fn: func() error { _, err := repo.StringSlice("missing"); return err }, want: `config: missing key "missing"`},

		{name: "missing map", fn: func() error { _, err := repo.Map("missing"); return err }, want: `config: missing key "missing"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn()

			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != tc.want {
				t.Fatalf("unexpected missing-key error: %v", err)
			}
		})
	}
}

func TestInternalHelpers(t *testing.T) {
	t.Parallel()

	if got := splitKey("  auth.cookies.secure  "); !reflect.DeepEqual(got, []string{"auth", "cookies", "secure"}) {
		t.Fatalf("unexpected split key: %#v", got)
	}

	if got := splitKey("   "); got != nil {
		t.Fatalf("expected nil split for empty key, got %#v", got)
	}

	if err := typeError("foo", "string", 123); err.Error() != `config: key "foo" must be string, got int` {
		t.Fatalf("unexpected type error: %v", err)
	}

	source := map[string]any{
		"nested": map[string]any{
			"value": []any{
				map[string]any{"key": "value"},
				[]string{"a", "b"},
			},
		},
	}
	cloned := cloneMap(source)
	clonedNested := cloned["nested"].(map[string]any)
	clonedSlice := clonedNested["value"].([]any)
	clonedSlice[0].(map[string]any)["key"] = "mutated"
	clonedSlice[1].([]string)[0] = "mutated"

	originalNested := source["nested"].(map[string]any)
	originalSlice := originalNested["value"].([]any)

	if originalSlice[0].(map[string]any)["key"] != "value" {
		t.Fatal("expected nested map clone protection")
	}

	if originalSlice[1].([]string)[0] != "a" {
		t.Fatal("expected []string clone protection")
	}

	if got := cloneValue([]string{"x", "y"}).([]string); !reflect.DeepEqual(got, []string{"x", "y"}) {
		t.Fatalf("unexpected []string clone: %#v", got)
	}

	if got := cloneValue(99).(int); got != 99 {
		t.Fatalf("unexpected scalar clone result: %d", got)
	}

	repo := NewRepository(map[string]any{
		"slice":       []any{"a", "b"},
		"stringSlice": []string{"x", "y"},
		"scalar":      "value",
	})

	if got := repo.sliceValue("slice"); !reflect.DeepEqual(got, []any{"a", "b"}) {
		t.Fatalf("unexpected []any slice value: %#v", got)
	}

	if got := repo.sliceValue("stringSlice"); !reflect.DeepEqual(got, []any{"x", "y"}) {
		t.Fatalf("unexpected []string slice value: %#v", got)
	}

	if got := repo.sliceValue("scalar"); len(got) != 0 {
		t.Fatalf("expected empty slice for scalar, got %#v", got)
	}

	if got := repo.sliceValue("missing"); len(got) != 0 {
		t.Fatalf("expected empty slice for missing key, got %#v", got)
	}

	if got := repo.Get("scalar.deeper", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback for scalar intermediate path, got %#v", got)
	}

	repo.setPath("branch.leaf", "value")

	if got := repo.Get("branch.leaf", nil); got != "value" {
		t.Fatalf("unexpected setPath value: %#v", got)
	}

	repo.setPath("branch.leaf.deep", "shadowed")

	if got := repo.Get("branch.leaf.deep", nil); got != "shadowed" {
		t.Fatalf("expected setPath to replace scalar with nested map, got %#v", got)
	}

	before := repo.All()
	repo.setPath("   ", "ignored")
	after := repo.All()

	if !reflect.DeepEqual(before, after) {
		t.Fatalf("expected empty setPath to be ignored: before=%#v after=%#v", before, after)
	}
}

func TestLookupHelperPaths(t *testing.T) {
	t.Parallel()

	items := map[string]any{
		"direct": "value",
		"a": map[string]any{
			"b": map[string]any{
				"c": "nested",
			},
			"b.c": "shadowed",
		},
		"plain": map[string]any{
			"nested": map[string]any{
				"leaf": "value",
			},
		},
		"scalar": "leaf",
	}

	if got, ok := lookup(items, "direct"); !ok || got != "value" {
		t.Fatalf("unexpected direct lookup: got=%#v ok=%v", got, ok)
	}

	if got, ok := lookup(items, ""); !ok || !reflect.DeepEqual(got, items) {
		t.Fatalf("unexpected root lookup: got=%#v ok=%v", got, ok)
	}

	if got, ok := lookup(items, "a.b.c"); !ok || got != "shadowed" {
		t.Fatalf("unexpected shadowed lookup: got=%#v ok=%v", got, ok)
	}

	if got, ok := lookup(items, "a.b"); !ok || !reflect.DeepEqual(got, map[string]any{"c": "nested"}) {
		t.Fatalf("unexpected nested lookup: got=%#v ok=%v", got, ok)
	}

	if got, ok := lookup(items, "plain.nested.leaf"); !ok || got != "value" {
		t.Fatalf("unexpected terminal nested lookup: got=%#v ok=%v", got, ok)
	}

	if got, ok := lookup(items, "a.missing"); ok || got != nil {
		t.Fatalf("expected nested miss: got=%#v ok=%v", got, ok)
	}

	if got, ok := lookup(items, "scalar.deep"); ok || got != nil {
		t.Fatalf("expected scalar intermediate miss: got=%#v ok=%v", got, ok)
	}

	if got, ok := lookup(items, "scalar.deep.more"); ok || got != nil {
		t.Fatalf("expected deep scalar intermediate miss: got=%#v ok=%v", got, ok)
	}
}

func TestRepositoryErrorMessagesStayStable(t *testing.T) {
	t.Parallel()

	repo := NewRepository(map[string]any{
		"fortify": map[string]any{
			"views": "not-a-bool",
		},
	})

	_, err := repo.Bool("fortify.views")

	if err == nil {
		t.Fatal("expected bool error")
	}

	if !strings.Contains(err.Error(), `config: key "fortify.views" must be bool`) {
		t.Fatalf("unexpected bool error message: %v", err)
	}
}
