package prompts

import "testing"

// Port of Laravel\Prompts\Tests\Feature\NotifyPromptTest

// Port of Laravel\Prompts\Tests\Feature\NotifyPromptTest::test_escape_applescript
func TestEscapeAppleScript(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{`say "hi"`, `say \"hi\"`},
		{`path\to\file`, `path\\to\\file`},
	}

	for _, tt := range tests {
		got := escapeAppleScript(tt.input)

		if got != tt.expected {
			t.Errorf("escapeAppleScript(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// Port of Laravel\Prompts\Tests\Feature\NotifyPromptTest::test_notify_options
func TestNotifyOptionsSetter(t *testing.T) {
	t.Parallel()
	cfg := &notifyConfig{}
	NotifyWithBody("body")(cfg)
	NotifyWithSubtitle("sub")(cfg)
	NotifyWithSound("Basso")(cfg)
	NotifyWithIcon("icon.png")(cfg)

	if cfg.body != "body" {
		t.Errorf("expected body %q, got %q", "body", cfg.body)
	}

	if cfg.subtitle != "sub" {
		t.Errorf("expected subtitle %q, got %q", "sub", cfg.subtitle)
	}

	if cfg.sound != "Basso" {
		t.Errorf("expected sound %q, got %q", "Basso", cfg.sound)
	}

	if cfg.icon != "icon.png" {
		t.Errorf("expected icon %q, got %q", "icon.png", cfg.icon)
	}
}
