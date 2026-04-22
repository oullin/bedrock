package process

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRunShellCommand(t *testing.T) {
	t.Parallel()

	result, err := New().Run(context.Background(), Shell("printf hello"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Successful() || result.Output() != "hello" {
		t.Fatalf("expected successful hello output, got code=%d output=%q", result.ExitCode(), result.Output())
	}
}

func TestRunCommandWithInputEnvAndPath(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("shell snippets are POSIX-specific")
	}

	result, err := New().
		Command(Shell("printf \"$BEDROCK_PROCESS:\" && pwd && cat")).
		Env(map[string]string{"BEDROCK_PROCESS": "ok"}).
		Path(t.TempDir()).
		Input("input").
		Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result.Output(), "ok:") || !strings.HasSuffix(result.Output(), "input") {
		t.Fatalf("expected env, path, and input in output, got %q", result.Output())
	}
}

func TestFakeProcessesAndAssertions(t *testing.T) {
	t.Parallel()

	manager := New().Fake("php artisan *", NewResult(Command{}, 0, "done", ""))

	result, err := manager.Run(context.Background(), Shell("php artisan queue:work"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Output() != "done" {
		t.Fatalf("expected fake output, got %q", result.Output())
	}

	if err := manager.AssertRan(Shell("php artisan queue:work")); err != nil {
		t.Fatalf("expected assertion to pass: %v", err)
	}
}

func TestPreventStrayProcesses(t *testing.T) {
	t.Parallel()

	_, err := New().PreventStrayProcesses().Run(context.Background(), Shell("printf no"))

	if !errors.Is(err, ErrStrayProcess) {
		t.Fatalf("expected ErrStrayProcess, got %v", err)
	}
}

func TestFakeSequenceExhaustion(t *testing.T) {
	t.Parallel()

	manager := New().Sequence("cmd", NewResult(Command{}, 0, "first", ""), NewResult(Command{}, 0, "second", ""))

	first, err := manager.Run(context.Background(), Shell("cmd"))

	if err != nil {
		t.Fatalf("unexpected first error: %v", err)
	}

	second, err := manager.Run(context.Background(), Shell("cmd"))

	if err != nil {
		t.Fatalf("unexpected second error: %v", err)
	}

	if first.Output() != "first" || second.Output() != "second" {
		t.Fatalf("expected ordered sequence, got %q then %q", first.Output(), second.Output())
	}

	_, err = manager.Run(context.Background(), Shell("cmd"))

	if !errors.Is(err, ErrSequenceEmpty) {
		t.Fatalf("expected ErrSequenceEmpty, got %v", err)
	}
}

func TestThrowOnFailedProcess(t *testing.T) {
	t.Parallel()

	_, err := New().
		Fake("fail", NewResult(Command{}, 2, "", "bad")).
		Command(Shell("fail")).
		Throw().
		Run(context.Background())

	if !errors.Is(err, ErrProcessFailed) {
		t.Fatalf("expected ErrProcessFailed, got %v", err)
	}
}

func TestStartAndWaitUntil(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("shell snippets are POSIX-specific")
	}

	invoked, err := New().Start(context.Background(), Shell("printf ready"))

	if err != nil {
		t.Fatalf("unexpected start error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	result, err := invoked.WaitUntil(ctx, func(output string, _ string) bool {
		return strings.Contains(output, "ready")
	})

	if err != nil {
		t.Fatalf("unexpected wait error: %v", err)
	}

	if result.Output() != "ready" {
		t.Fatalf("expected ready output, got %q", result.Output())
	}
}

func TestPoolAndPipe(t *testing.T) {
	t.Parallel()

	manager := New().
		Fake("one", NewResult(Command{}, 0, "1", "")).
		Fake("two", NewResult(Command{}, 0, "2", "")).
		Fake("cat", NewResult(Command{}, 0, "piped", ""))

	results, err := manager.Pool().Command("a", Shell("one")).Command("b", Shell("two")).Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected pool error: %v", err)
	}

	if results["a"].Output() != "1" || results["b"].Output() != "2" {
		t.Fatalf("unexpected pool results: %#v", results)
	}

	result, err := manager.Pipe(Shell("cat")).Run(context.Background())

	if err != nil {
		t.Fatalf("unexpected pipe error: %v", err)
	}

	if result.Output() != "piped" {
		t.Fatalf("expected piped output, got %q", result.Output())
	}
}
