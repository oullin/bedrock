package bootstrap_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/bootstrap"
	"github.com/bedrock/packages/container"
)

type sample struct{ x int }

func resetBedrock() {
	bootstrap.SetApp(nil)
}

func TestSetApp_AndApp(t *testing.T) {
	t.Cleanup(resetBedrock)

	application := container.NewApplication()
	bootstrap.SetApp(application)

	if got := bootstrap.App(); got != application {
		t.Fatalf("App() returned %p, want %p", got, application)
	}
}

func TestApp_PanicsWhenNotInstalled(t *testing.T) {
	t.Cleanup(resetBedrock)
	bootstrap.SetApp(nil)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when no app installed")
		}
	}()

	bootstrap.App()
}

func TestHasApp(t *testing.T) {
	t.Cleanup(resetBedrock)
	bootstrap.SetApp(nil)

	if bootstrap.HasApp() {
		t.Fatal("HasApp should be false before SetApp")
	}

	bootstrap.SetApp(container.NewApplication())

	if !bootstrap.HasApp() {
		t.Fatal("HasApp should be true after SetApp")
	}
}

func TestMustMake_ResolvesValue(t *testing.T) {
	t.Cleanup(resetBedrock)

	application := container.NewApplication()
	application.Instance("answer", 42)
	bootstrap.SetApp(application)

	if got := bootstrap.MustMake("answer"); got != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

func TestMustMake_PanicsOnMiss(t *testing.T) {
	t.Cleanup(resetBedrock)

	bootstrap.SetApp(container.NewApplication())

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on missing binding")
		}
	}()

	bootstrap.MustMake("nope")
}

func TestResolve_GenericTypedAccess(t *testing.T) {
	t.Cleanup(resetBedrock)

	application := container.NewApplication()
	application.Instance("sample", &sample{x: 7})
	bootstrap.SetApp(application)

	got := bootstrap.Resolve[*sample]("sample")

	if got.x != 7 {
		t.Fatalf("expected x=7, got %d", got.x)
	}
}

func TestResolve_PanicsOnWrongType(t *testing.T) {
	t.Cleanup(resetBedrock)

	application := container.NewApplication()
	application.Instance("sample", "not a *sample")
	bootstrap.SetApp(application)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on wrong type")
		}
	}()

	bootstrap.Resolve[*sample]("sample")
}

func TestTryResolve_ReturnsErrorOnMiss(t *testing.T) {
	t.Cleanup(resetBedrock)

	bootstrap.SetApp(container.NewApplication())

	_, err := bootstrap.TryResolve[*sample]("nope")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, container.ErrNotBound) {
		t.Fatalf("expected ErrNotBound, got %v", err)
	}
}

func TestTryResolve_ReturnsErrorWhenNoApp(t *testing.T) {
	t.Cleanup(resetBedrock)
	bootstrap.SetApp(nil)

	_, err := bootstrap.TryResolve[*sample]("anything")

	if err == nil {
		t.Fatal("expected error when no app installed")
	}
}

func TestTryResolve_ReturnsErrorOnWrongType(t *testing.T) {
	t.Cleanup(resetBedrock)

	application := container.NewApplication()
	application.Instance("sample", 42)
	bootstrap.SetApp(application)

	_, err := bootstrap.TryResolve[*sample]("sample")

	if err == nil {
		t.Fatal("expected type error")
	}
}
