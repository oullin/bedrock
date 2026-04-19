package app_test

import (
	"errors"
	"testing"

	bedrockapp "github.com/bedrock/app"
	"github.com/bedrock/packages/container"
)

type sample struct{ x int }

func resetBedrock() {
	bedrockapp.SetApp(nil)
}

func TestSetApp_AndApp(t *testing.T) {
	t.Cleanup(resetBedrock)

	application := container.NewApplication()
	bedrockapp.SetApp(application)

	if got := bedrockapp.App(); got != application {
		t.Fatalf("App() returned %p, want %p", got, application)
	}
}

func TestApp_PanicsWhenNotInstalled(t *testing.T) {
	t.Cleanup(resetBedrock)
	bedrockapp.SetApp(nil)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when no app installed")
		}
	}()

	bedrockapp.App()
}

func TestHasApp(t *testing.T) {
	t.Cleanup(resetBedrock)
	bedrockapp.SetApp(nil)

	if bedrockapp.HasApp() {
		t.Fatal("HasApp should be false before SetApp")
	}

	bedrockapp.SetApp(container.NewApplication())

	if !bedrockapp.HasApp() {
		t.Fatal("HasApp should be true after SetApp")
	}
}

func TestMustMake_ResolvesValue(t *testing.T) {
	t.Cleanup(resetBedrock)

	application := container.NewApplication()
	application.Instance("answer", 42)
	bedrockapp.SetApp(application)

	if got := bedrockapp.MustMake("answer"); got != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

func TestMustMake_PanicsOnMiss(t *testing.T) {
	t.Cleanup(resetBedrock)

	bedrockapp.SetApp(container.NewApplication())

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on missing binding")
		}
	}()

	bedrockapp.MustMake("nope")
}

func TestResolve_GenericTypedAccess(t *testing.T) {
	t.Cleanup(resetBedrock)

	application := container.NewApplication()
	application.Instance("sample", &sample{x: 7})
	bedrockapp.SetApp(application)

	got := bedrockapp.Resolve[*sample]("sample")

	if got.x != 7 {
		t.Fatalf("expected x=7, got %d", got.x)
	}
}

func TestResolve_PanicsOnWrongType(t *testing.T) {
	t.Cleanup(resetBedrock)

	application := container.NewApplication()
	application.Instance("sample", "not a *sample")
	bedrockapp.SetApp(application)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on wrong type")
		}
	}()

	bedrockapp.Resolve[*sample]("sample")
}

func TestTryResolve_ReturnsErrorOnMiss(t *testing.T) {
	t.Cleanup(resetBedrock)

	bedrockapp.SetApp(container.NewApplication())

	_, err := bedrockapp.TryResolve[*sample]("nope")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, container.ErrNotBound) {
		t.Fatalf("expected ErrNotBound, got %v", err)
	}
}

func TestTryResolve_ReturnsErrorWhenNoApp(t *testing.T) {
	t.Cleanup(resetBedrock)
	bedrockapp.SetApp(nil)

	_, err := bedrockapp.TryResolve[*sample]("anything")

	if err == nil {
		t.Fatal("expected error when no app installed")
	}
}

func TestTryResolve_ReturnsErrorOnWrongType(t *testing.T) {
	t.Cleanup(resetBedrock)

	application := container.NewApplication()
	application.Instance("sample", 42)
	bedrockapp.SetApp(application)

	_, err := bedrockapp.TryResolve[*sample]("sample")

	if err == nil {
		t.Fatal("expected type error")
	}
}
