package config

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestRepositoryDotPathAccessors(t *testing.T) {
	t.Parallel()

	repo := NewRepository(map[string]any{
		"auth": map[string]any{
			"session_lifetime": "24h",
			"cookies": map[string]any{
				"secure": true,
			},
			"middleware": []any{"web", "auth"},
		},
	})

	if !repo.Has("auth.session_lifetime") {
		t.Fatal("expected key to exist")
	}

	if got := repo.Get("auth.missing", "fallback"); got != "fallback" {
		t.Fatalf("unexpected fallback value: %v", got)
	}

	duration, err := repo.Duration("auth.session_lifetime")
	if err != nil {
		t.Fatalf("Duration: %v", err)
	}
	if duration != 24*time.Hour {
		t.Fatalf("unexpected duration: %v", duration)
	}

	secure, err := repo.Bool("auth.cookies.secure")
	if err != nil {
		t.Fatalf("Bool: %v", err)
	}
	if !secure {
		t.Fatal("expected secure cookie")
	}

	middleware, err := repo.StringSlice("auth.middleware")
	if err != nil {
		t.Fatalf("StringSlice: %v", err)
	}
	if !reflect.DeepEqual(middleware, []string{"web", "auth"}) {
		t.Fatalf("unexpected middleware: %#v", middleware)
	}
}

func TestRepositorySetMutatorsAndClone(t *testing.T) {
	t.Parallel()

	repo := NewRepository(nil)
	repo.Set("authflows.limiters.login", "login")
	repo.Set(map[string]any{
		"authflows.limiters.two-factor":                       "two-factor",
		"authflows.options.two-factor-authentication.confirm": true,
	})
	repo.Push("authflows.middleware", "web")
	repo.Push("authflows.middleware", "guest")
	repo.Prepend("authflows.middleware", "trim")

	mapped, err := repo.Map("authflows.limiters")
	if err != nil {
		t.Fatalf("Map: %v", err)
	}
	if mapped["login"] != "login" {
		t.Fatalf("unexpected limiter map: %#v", mapped)
	}
	if mapped["two-factor"] != "two-factor" {
		t.Fatalf("unexpected two-factor limiter: %#v", mapped)
	}

	middleware, err := repo.StringSlice("authflows.middleware")
	if err != nil {
		t.Fatalf("StringSlice: %v", err)
	}
	if !reflect.DeepEqual(middleware, []string{"trim", "web", "guest"}) {
		t.Fatalf("unexpected middleware order: %#v", middleware)
	}

	confirm, err := repo.Bool("authflows.options.two-factor-authentication.confirm")
	if err != nil {
		t.Fatalf("Bool: %v", err)
	}
	if !confirm {
		t.Fatal("expected confirm option to be enabled")
	}

	all := repo.All()
	inner := all["authflows"].(map[string]any)
	limits := inner["limiters"].(map[string]any)
	limits["login"] = "mutated"

	value, err := repo.String("authflows.limiters.login")
	if err != nil {
		t.Fatalf("String: %v", err)
	}
	if value != "login" {
		t.Fatalf("repository was mutated through clone: %s", value)
	}
}

func TestRepositoryLiteralDottedKeysTakePrecedence(t *testing.T) {
	t.Parallel()

	repo := NewRepository(map[string]any{
		"mail.mailers.smtp": "literal",
		"mail": map[string]any{
			"mailers": map[string]any{
				"smtp": "nested",
			},
		},
	})

	if got := repo.Get("mail.mailers.smtp", nil); got != "literal" {
		t.Fatalf("unexpected dotted value: %#v", got)
	}
}

func TestRepositoryGetManyEmpty(t *testing.T) {
	t.Parallel()

	repo := NewRepository(nil)
	values := repo.GetMany(nil)
	if len(values) != 0 {
		t.Fatalf("expected empty values, got %#v", values)
	}
}

func TestRepositoryGetEmptyKeyReturnsAll(t *testing.T) {
	t.Parallel()

	repo := NewRepository(map[string]any{"app": map[string]any{"name": "gollin"}})
	got := repo.Get("", nil)
	all := repo.All()
	if !reflect.DeepEqual(got, all) {
		t.Fatalf("unexpected root value: %#v", got)
	}
}

func TestRepositoryTypedGetterFailures(t *testing.T) {
	t.Parallel()

	repo := NewRepository(map[string]any{
		"authflows": map[string]any{
			"views":      "not-a-bool",
			"features":   []any{"registration", 42},
			"middleware": "web",
		},
	})

	if _, err := repo.Bool("authflows.views"); err == nil {
		t.Fatal("expected bool type error")
	}
	if _, err := repo.String("authflows.views"); err != nil {
		t.Fatalf("String: %v", err)
	}
	if _, err := repo.StringSlice("authflows.features"); err == nil {
		t.Fatal("expected string slice type error")
	}
	if _, err := repo.Map("authflows.middleware"); err == nil {
		t.Fatal("expected map type error")
	}
	if _, err := repo.Duration("authflows.middleware"); err == nil {
		t.Fatal("expected duration type error")
	}
}

func TestRepositoryMissingKeyErrors(t *testing.T) {
	t.Parallel()

	repo := NewRepository(nil)
	_, err := repo.String("missing")
	if err == nil {
		t.Fatal("expected missing key error")
	}
	if !errors.Is(err, err) {
		t.Fatal("expected non-nil error")
	}
}
