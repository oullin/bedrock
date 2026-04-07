package console_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/bedrock/packages/anvil/console"
)

// Laravel: testCommandRegistration
func TestCommandRegistration(t *testing.T) {
	t.Parallel()

	k := console.New()
	k.Command("test:hello", "Says hello", func(_ context.Context, inv *console.Invocation) error {
		inv.Comment("hello world")
		return nil
	})

	commands := k.Commands()
	if len(commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(commands))
	}
	if commands[0].Name != "test:hello" {
		t.Fatalf("expected name 'test:hello', got %q", commands[0].Name)
	}
	if commands[0].Purpose != "Says hello" {
		t.Fatalf("expected purpose 'Says hello', got %q", commands[0].Purpose)
	}
}

// Laravel: testCommandExecution
func TestCommandExecution(t *testing.T) {
	t.Parallel()

	k := console.New()
	k.Command("greet", "Greets the user", func(_ context.Context, inv *console.Invocation) error {
		if len(inv.Args) > 0 {
			inv.Comment("hello " + inv.Args[0])
		} else {
			inv.Comment("hello world")
		}
		return nil
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := k.Handle([]string{"greet", "Alice"}, stdout, stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if got := strings.TrimSpace(stdout.String()); got != "hello Alice" {
		t.Fatalf("expected 'hello Alice', got %q", got)
	}
}

// Laravel: testCommandExitCodes
func TestCommandExitCodeOnError(t *testing.T) {
	t.Parallel()

	k := console.New()
	k.Command("fail", "Always fails", func(_ context.Context, _ *console.Invocation) error {
		return errors.New("something went wrong")
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := k.Handle([]string{"fail"}, stdout, stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "something went wrong") {
		t.Fatalf("expected error message in stderr, got %q", stderr.String())
	}
}

// Test unknown command
func TestUnknownCommandReturnsError(t *testing.T) {
	t.Parallel()

	k := console.New()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := k.Handle([]string{"nonexistent"}, stdout, stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("expected 'unknown command' in stderr, got %q", stderr.String())
	}
}

// Test listing commands when no args
func TestHandleWithNoArgsListsCommands(t *testing.T) {
	t.Parallel()

	k := console.New()
	k.Command("alpha", "First command", func(context.Context, *console.Invocation) error { return nil })
	k.Command("beta", "Second command", func(context.Context, *console.Invocation) error { return nil })

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := k.Handle([]string{}, stdout, stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	output := stdout.String()
	if !strings.Contains(output, "available commands") {
		t.Fatal("expected 'available commands' header")
	}
	if !strings.Contains(output, "alpha") || !strings.Contains(output, "beta") {
		t.Fatal("expected both commands in listing")
	}
}

// Test commands are sorted
func TestCommandsSortedByName(t *testing.T) {
	t.Parallel()

	k := console.New()
	k.Command("zulu", "", func(context.Context, *console.Invocation) error { return nil })
	k.Command("alpha", "", func(context.Context, *console.Invocation) error { return nil })
	k.Command("mike", "", func(context.Context, *console.Invocation) error { return nil })

	commands := k.Commands()
	if len(commands) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(commands))
	}
	if commands[0].Name != "alpha" || commands[1].Name != "mike" || commands[2].Name != "zulu" {
		t.Fatalf("expected sorted order [alpha, mike, zulu], got [%s, %s, %s]",
			commands[0].Name, commands[1].Name, commands[2].Name)
	}
}

// Test Invocation.Comment
func TestInvocationComment(t *testing.T) {
	t.Parallel()

	stdout := &bytes.Buffer{}
	inv := &console.Invocation{
		Args:   []string{"arg1"},
		Stdout: stdout,
		Stderr: &bytes.Buffer{},
	}

	inv.Comment("test message")
	if got := strings.TrimSpace(stdout.String()); got != "test message" {
		t.Fatalf("expected 'test message', got %q", got)
	}
}

// Test empty kernel
func TestEmptyKernelCommands(t *testing.T) {
	t.Parallel()

	k := console.New()
	if len(k.Commands()) != 0 {
		t.Fatal("expected empty commands list")
	}
}

// Test command overwrites
func TestCommandOverwrite(t *testing.T) {
	t.Parallel()

	k := console.New()
	k.Command("test", "v1", func(_ context.Context, inv *console.Invocation) error {
		inv.Comment("v1")
		return nil
	})
	k.Command("test", "v2", func(_ context.Context, inv *console.Invocation) error {
		inv.Comment("v2")
		return nil
	})

	commands := k.Commands()
	if len(commands) != 1 {
		t.Fatalf("expected 1 command after overwrite, got %d", len(commands))
	}
	if commands[0].Purpose != "v2" {
		t.Fatalf("expected overwritten purpose 'v2', got %q", commands[0].Purpose)
	}

	stdout := &bytes.Buffer{}
	k.Handle([]string{"test"}, stdout, &bytes.Buffer{})
	if got := strings.TrimSpace(stdout.String()); got != "v2" {
		t.Fatalf("expected 'v2' output, got %q", got)
	}
}
