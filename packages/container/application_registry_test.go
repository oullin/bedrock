package container_test

import (
	"testing"

	"github.com/bedrock/packages/container"
)

type sample struct{ Value string }

func TestSetAppAndApp(t *testing.T) {
	t.Cleanup(func() {
		container.SetApp(nil)
	})

	container.SetApp(nil)

	application := container.NewApplication()
	container.SetApp(application)

	if got := container.App(); got != application {
		t.Fatalf("App() = %p, want %p", got, application)
	}
}

func TestAppPanicsWhenUnset(t *testing.T) {
	t.Cleanup(func() {
		container.SetApp(nil)
	})

	container.SetApp(nil)

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when app is unset")
		}
	}()

	container.App()
}

func TestHasApp(t *testing.T) {
	t.Cleanup(func() {
		container.SetApp(nil)
	})

	container.SetApp(nil)

	if container.HasApp() {
		t.Fatal("HasApp should be false after SetApp(nil)")
	}

	container.SetApp(container.NewApplication())

	if !container.HasApp() {
		t.Fatal("HasApp should be true after installing an app")
	}
}

func TestMustMake(t *testing.T) {
	t.Cleanup(func() {
		container.SetApp(nil)
	})

	application := container.NewApplication()
	application.Instance("answer", 42)
	container.SetApp(application)

	if got := container.MustMake("answer"); got != 42 {
		t.Fatalf("MustMake = %v, want 42", got)
	}
}

func TestMustMakePanicsWhenMissing(t *testing.T) {
	t.Cleanup(func() {
		container.SetApp(nil)
	})

	container.SetApp(container.NewApplication())

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for missing binding")
		}
	}()

	container.MustMake("nope")
}

func TestResolve(t *testing.T) {
	t.Cleanup(func() {
		container.SetApp(nil)
	})

	application := container.NewApplication()
	application.Instance("sample", &sample{Value: "ok"})
	container.SetApp(application)

	got := container.Resolve[*sample]("sample")

	if got.Value != "ok" {
		t.Fatalf("Resolve returned %+v", got)
	}
}

func TestResolvePanicsOnWrongType(t *testing.T) {
	t.Cleanup(func() {
		container.SetApp(nil)
	})

	application := container.NewApplication()
	application.Instance("sample", 123)
	container.SetApp(application)

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for wrong type")
		}
	}()

	container.Resolve[*sample]("sample")
}

func TestTryResolveMissing(t *testing.T) {
	t.Cleanup(func() {
		container.SetApp(nil)
	})

	container.SetApp(container.NewApplication())

	_, err := container.TryResolve[*sample]("nope")

	if err == nil {
		t.Fatal("expected error for missing binding")
	}
}

func TestTryResolveNoApp(t *testing.T) {
	t.Cleanup(func() {
		container.SetApp(nil)
	})

	container.SetApp(nil)

	_, err := container.TryResolve[*sample]("anything")

	if err == nil {
		t.Fatal("expected error when app is unset")
	}
}

func TestTryResolveWrongType(t *testing.T) {
	t.Cleanup(func() {
		container.SetApp(nil)
	})

	application := container.NewApplication()
	application.Instance("sample", 123)
	container.SetApp(application)

	_, err := container.TryResolve[*sample]("sample")

	if err == nil {
		t.Fatal("expected error for wrong type")
	}
}
