package prompts

import (
	"os/exec"
	"runtime"
)

// NotifyOption configures a desktop notification.
type NotifyOption func(*notifyConfig)

type notifyConfig struct {
	body     string
	subtitle string
	sound    string
	icon     string
}

// NotifyWithBody sets the notification body.
func NotifyWithBody(s string) NotifyOption { return func(c *notifyConfig) { c.body = s } }

// NotifyWithSubtitle sets the notification subtitle (macOS only).
func NotifyWithSubtitle(s string) NotifyOption { return func(c *notifyConfig) { c.subtitle = s } }

// NotifyWithSound sets the notification sound (macOS only).
func NotifyWithSound(s string) NotifyOption { return func(c *notifyConfig) { c.sound = s } }

// NotifyWithIcon sets the notification icon (Linux only).
func NotifyWithIcon(s string) NotifyOption { return func(c *notifyConfig) { c.icon = s } }

// Notify sends a cross-platform desktop notification.
func Notify(title string, opts ...NotifyOption) {
	cfg := &notifyConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	switch runtime.GOOS {
	case "darwin":
		notifyMacOS(title, cfg)
	case "linux":
		notifyLinux(title, cfg)
	}
}

func notifyMacOS(title string, cfg *notifyConfig) {
	script := "display notification"

	if cfg.body != "" {
		script += " \"" + escapeAppleScript(cfg.body) + "\""
	} else {
		script += " \"\""
	}

	if title != "" {
		script += " with title \"" + escapeAppleScript(title) + "\""
	}

	if cfg.subtitle != "" {
		script += " subtitle \"" + escapeAppleScript(cfg.subtitle) + "\""
	}

	if cfg.sound != "" {
		script += " sound name \"" + escapeAppleScript(cfg.sound) + "\""
	}

	_ = exec.Command("osascript", "-e", script).Run()
}

func notifyLinux(title string, cfg *notifyConfig) {
	args := []string{title}

	if cfg.body != "" {
		args = append(args, cfg.body)
	}

	if cfg.icon != "" {
		args = append([]string{"-i", cfg.icon}, args...)
	}

	if _, err := exec.LookPath("notify-send"); err == nil {
		_ = exec.Command("notify-send", args...).Run()

		return
	}

	if _, err := exec.LookPath("kdialog"); err == nil {
		body := title

		if cfg.body != "" {
			body += "\n" + cfg.body
		}

		_ = exec.Command("kdialog", "--passivepopup", body, "5").Run()
	}
}

func escapeAppleScript(s string) string {
	result := make([]byte, 0, len(s))

	for i := 0; i < len(s); i++ {
		if s[i] == '"' || s[i] == '\\' {
			result = append(result, '\\')
		}

		result = append(result, s[i])
	}

	return string(result)
}
