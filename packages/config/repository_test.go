package config

import (
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

func TestRepositorySetAndClone(t *testing.T) {
	t.Parallel()

	repo := NewRepository(nil)
	repo.Set("authflows.limiters.login", "login")

	mapped, err := repo.Map("authflows.limiters")
	if err != nil {
		t.Fatalf("Map: %v", err)
	}
	if mapped["login"] != "login" {
		t.Fatalf("unexpected limiter map: %#v", mapped)
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
