package notifications_test

import (
	"testing"

	"github.com/bedrock/packages/notifications"
)

func TestSimpleMessageDefaultLevel(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage()

	if m.GetLevel() != "info" {
		t.Fatalf("default level = %q, want %q", m.GetLevel(), "info")
	}
}

func TestSimpleMessageSuccess(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().Success()

	if m.GetLevel() != "success" {
		t.Fatalf("level = %q, want %q", m.GetLevel(), "success")
	}
}

func TestSimpleMessageError(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().Error()

	if m.GetLevel() != "error" {
		t.Fatalf("level = %q, want %q", m.GetLevel(), "error")
	}
}

func TestSimpleMessageLevel(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().Level("warning")

	if m.GetLevel() != "warning" {
		t.Fatalf("level = %q, want %q", m.GetLevel(), "warning")
	}
}

func TestSimpleMessageSubject(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().Subject("Hello")

	if m.GetSubject() != "Hello" {
		t.Fatalf("subject = %q, want %q", m.GetSubject(), "Hello")
	}
}

func TestSimpleMessageGreeting(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().Greeting("Hi there!")

	if m.GetGreeting() != "Hi there!" {
		t.Fatalf("greeting = %q, want %q", m.GetGreeting(), "Hi there!")
	}
}

func TestSimpleMessageSalutation(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().Salutation("Regards")

	if m.GetSalutation() != "Regards" {
		t.Fatalf("salutation = %q, want %q", m.GetSalutation(), "Regards")
	}
}

func TestSimpleMessageLinesBeforeAction(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().
		Line("First line").
		Line("Second line")

	if len(m.GetIntroLines()) != 2 {
		t.Fatalf("intro lines count = %d, want 2", len(m.GetIntroLines()))
	}

	if len(m.GetOutroLines()) != 0 {
		t.Fatalf("outro lines count = %d, want 0", len(m.GetOutroLines()))
	}
}

func TestSimpleMessageLinesAfterAction(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().
		Line("Intro").
		Action("Click", "https://example.com").
		Line("Outro")

	if len(m.GetIntroLines()) != 1 {
		t.Fatalf("intro lines count = %d, want 1", len(m.GetIntroLines()))
	}

	if len(m.GetOutroLines()) != 1 {
		t.Fatalf("outro lines count = %d, want 1", len(m.GetOutroLines()))
	}
}

func TestSimpleMessageLineIf(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().
		LineIf(true, "Visible").
		LineIf(false, "Hidden")

	if len(m.GetIntroLines()) != 1 {
		t.Fatalf("intro lines count = %d, want 1", len(m.GetIntroLines()))
	}

	if m.GetIntroLines()[0] != "Visible" {
		t.Fatalf("line = %q, want %q", m.GetIntroLines()[0], "Visible")
	}
}

func TestSimpleMessageLines(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().Lines([]string{"A", "B", "C"})

	if len(m.GetIntroLines()) != 3 {
		t.Fatalf("intro lines count = %d, want 3", len(m.GetIntroLines()))
	}
}

func TestSimpleMessageLinesIf(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().
		LinesIf(false, []string{"A", "B"})

	if len(m.GetIntroLines()) != 0 {
		t.Fatalf("intro lines count = %d, want 0", len(m.GetIntroLines()))
	}
}

func TestSimpleMessageWith(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().With("Test line")

	if len(m.GetIntroLines()) != 1 {
		t.Fatalf("intro lines count = %d, want 1", len(m.GetIntroLines()))
	}
}

func TestSimpleMessageAction(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().Action("Go", "https://go.dev")

	if m.GetActionText() != "Go" {
		t.Fatalf("action text = %q, want %q", m.GetActionText(), "Go")
	}

	if m.GetActionURL() != "https://go.dev" {
		t.Fatalf("action URL = %q, want %q", m.GetActionURL(), "https://go.dev")
	}
}

func TestSimpleMessageMailer(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().Mailer("postmark")

	if m.GetMailer() != "postmark" {
		t.Fatalf("mailer = %q, want %q", m.GetMailer(), "postmark")
	}
}

func TestSimpleMessageToMap(t *testing.T) {
	t.Parallel()

	m := notifications.NewSimpleMessage().
		Subject("Test").
		Level("success").
		Line("Hello")

	data := m.ToMap()

	if data["subject"] != "Test" {
		t.Fatalf("subject = %v, want %q", data["subject"], "Test")
	}

	if data["level"] != "success" {
		t.Fatalf("level = %v, want %q", data["level"], "success")
	}
}
