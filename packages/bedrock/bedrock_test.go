package bedrock_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/bedrock"
	"github.com/bedrock/packages/container"
)

type sample struct{ x int }

func resetBedrock() {
	bedrock.SetApp(nil)
}

func TestSetApp_AndApp(t *testing.T) {
	t.Cleanup(resetBedrock)

	app := container.NewApplication()
	bedrock.SetApp(app)

	if got := bedrock.App(); got != app {
		t.Fatalf("App() returned %p, want %p", got, app)
	}
}

func TestApp_PanicsWhenNotInstalled(t *testing.T) {
	t.Cleanup(resetBedrock)
	bedrock.SetApp(nil)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when no app installed")
		}
	}()

	bedrock.App()
}

func TestHasApp(t *testing.T) {
	t.Cleanup(resetBedrock)
	bedrock.SetApp(nil)

	if bedrock.HasApp() {
		t.Fatal("HasApp should be false before SetApp")
	}

	bedrock.SetApp(container.NewApplication())

	if !bedrock.HasApp() {
		t.Fatal("HasApp should be true after SetApp")
	}
}

func TestMustMake_ResolvesValue(t *testing.T) {
	t.Cleanup(resetBedrock)

	app := container.NewApplication()
	app.Instance("answer", 42)
	bedrock.SetApp(app)

	if got := bedrock.MustMake("answer"); got != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

func TestMustMake_PanicsOnMiss(t *testing.T) {
	t.Cleanup(resetBedrock)

	bedrock.SetApp(container.NewApplication())

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on missing binding")
		}
	}()

	bedrock.MustMake("nope")
}

func TestResolve_GenericTypedAccess(t *testing.T) {
	t.Cleanup(resetBedrock)

	app := container.NewApplication()
	app.Instance("sample", &sample{x: 7})
	bedrock.SetApp(app)

	got := bedrock.Resolve[*sample]("sample")

	if got.x != 7 {
		t.Fatalf("expected x=7, got %d", got.x)
	}
}

func TestResolve_PanicsOnWrongType(t *testing.T) {
	t.Cleanup(resetBedrock)

	app := container.NewApplication()
	app.Instance("sample", "not a *sample")
	bedrock.SetApp(app)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on wrong type")
		}
	}()

	bedrock.Resolve[*sample]("sample")
}

func TestTryResolve_ReturnsErrorOnMiss(t *testing.T) {
	t.Cleanup(resetBedrock)

	bedrock.SetApp(container.NewApplication())

	_, err := bedrock.TryResolve[*sample]("nope")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, container.ErrNotBound) {
		t.Fatalf("expected ErrNotBound, got %v", err)
	}
}

func TestTryResolve_ReturnsErrorWhenNoApp(t *testing.T) {
	t.Cleanup(resetBedrock)
	bedrock.SetApp(nil)

	_, err := bedrock.TryResolve[*sample]("anything")

	if err == nil {
		t.Fatal("expected error when no app installed")
	}
}

func TestTryResolve_ReturnsErrorOnWrongType(t *testing.T) {
	t.Cleanup(resetBedrock)

	app := container.NewApplication()
	app.Instance("sample", 42)
	bedrock.SetApp(app)

	_, err := bedrock.TryResolve[*sample]("sample")

	if err == nil {
		t.Fatal("expected type error")
	}
}
