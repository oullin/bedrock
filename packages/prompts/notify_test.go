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

// NotifyPromptTest::it_builds_the_correct_macos_command
func TestMacOSNotificationScript(t *testing.T) {
	t.Parallel()

	got := macOSNotificationScript("Deploy", &notifyConfig{body: "Done"})
	want := `display notification "Done" with title "Deploy"`

	if got != want {
		t.Fatalf("macOS script = %q, want %q", got, want)
	}
}

// NotifyPromptTest::it_includes_subtitle_and_sound_in_macos_command
func TestMacOSNotificationScriptIncludesSubtitleAndSound(t *testing.T) {
	t.Parallel()

	got := macOSNotificationScript("Deploy", &notifyConfig{
		body:     "Done",
		subtitle: "Release",
		sound:    "Basso",
	})

	if got != `display notification "Done" with title "Deploy" subtitle "Release" sound name "Basso"` {
		t.Fatalf("macOS script = %q", got)
	}
}

// NotifyPromptTest::it_sets_linux_options
func TestLinuxNotificationArgs(t *testing.T) {
	t.Parallel()

	got := linuxNotificationArgs("Deploy", &notifyConfig{body: "Done", icon: "app.png"})
	want := []string{"-i", "app.png", "Deploy", "Done"}

	if len(got) != len(want) {
		t.Fatalf("linux args = %#v, want %#v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("linux args = %#v, want %#v", got, want)
		}
	}
}
