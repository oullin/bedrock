package envoy

import (
	"context"
	"errors"
	"testing"
)

func TestPlanRendersCommandsForHosts(t *testing.T) {
	t.Parallel()

	steps, err := Plan(Task{
		Name:     "deploy",
		Hosts:    []string{"web-2", "web-1"},
		Commands: []string{"cd {{ .release }}", "php artisan migrate --force"},
		Variables: map[string]any{
			"release": "/srv/app/current",
		},
	}, nil)

	if err != nil {
		t.Fatalf("Plan returned error: %v", err)
	}

	if len(steps) != 4 {
		t.Fatalf("len(steps) = %d, want 4", len(steps))
	}

	if steps[0].Host != "web-1" {
		t.Fatalf("first host = %q, want web-1", steps[0].Host)
	}

	if steps[0].Command != "cd /srv/app/current" {
		t.Fatalf("first command = %q", steps[0].Command)
	}

	if steps[1].Command != "php artisan migrate --force" {
		t.Fatalf("second command = %q", steps[1].Command)
	}
}

func TestRunUsesRunner(t *testing.T) {
	t.Parallel()

	var called bool
	results, err := Run(context.Background(), RunnerFunc(func(ctx context.Context, step Step) (Result, error) {
		called = true

		return Result{Step: step, Stdout: "ok"}, nil
	}), Task{
		Name:     "status",
		Commands: []string{"whoami"},
	}, nil)

	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if !called {
		t.Fatal("runner was not called")
	}

	if results[0].Stdout != "ok" {
		t.Fatalf("stdout = %q, want ok", results[0].Stdout)
	}
}

func TestPlanRequiresTaskName(t *testing.T) {
	t.Parallel()

	_, err := Plan(Task{Commands: []string{"echo ok"}}, nil)

	if !errors.Is(err, ErrMissingTaskName) {
		t.Fatalf("Plan error = %v, want ErrMissingTaskName", err)
	}
}
