package notifications_test

import (
	"testing"

	"github.com/bedrock/packages/notifications"
)

func TestNewAction(t *testing.T) {
	t.Parallel()

	a := notifications.NewAction("Click Me", "https://example.com")

	if a.Text != "Click Me" {
		t.Fatalf("Text = %q, want %q", a.Text, "Click Me")
	}

	if a.URL != "https://example.com" {
		t.Fatalf("URL = %q, want %q", a.URL, "https://example.com")
	}
}

func TestActionStructLiteral(t *testing.T) {
	t.Parallel()

	a := notifications.Action{Text: "Submit", URL: "/submit"}

	if a.Text != "Submit" || a.URL != "/submit" {
		t.Fatalf("unexpected Action: %+v", a)
	}
}
