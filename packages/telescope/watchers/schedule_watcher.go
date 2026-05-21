package watchers

import (
	"github.com/bedrock/packages/telescope"
)

// ScheduleWatcher monitors scheduled task execution and records entries as
// Telescope entries. It mirrors the the underlying behavior class.
type ScheduleWatcher struct {
	telescope.BaseWatcher
}

// NewScheduleWatcher creates a ScheduleWatcher with the given options.

// Register is a no-op for ScheduleWatcher; callers drive it via Record.

// ScheduledTask carries metadata about a scheduled task execution.
type ScheduledTask struct {
	// Description is a human-readable description of the task.
	Description string
	// Command is the CLI command or callable identifier.
	Command string
	// Expression is the cron expression (e.g. "0 * * * *").
	Expression string
	// Timezone is the timezone in which the cron expression is evaluated.
	Timezone string
	// User is the system user that runs the command (optional).
	User string
	// Output captured from the task (optional).
	Output string
	// ExitCode is the exit code of the command.
	ExitCode int
}

func NewScheduleWatcher(t *telescope.Telescope, options map[string]any) *ScheduleWatcher {
	w := &ScheduleWatcher{}
	w.SetTelescope(t)
	w.Options = options

	return w
}

func (w *ScheduleWatcher) Register(_ any) error { return nil }

// Record records a scheduled task execution entry.
func (w *ScheduleWatcher) Record(task ScheduledTask) {
	content := map[string]any{
		"description": task.Description,
		"command":     task.Command,
		"expression":  task.Expression,
		"timezone":    task.Timezone,
		"user":        task.User,
		"output":      task.Output,
		"exit_code":   task.ExitCode,
	}

	entry := telescope.NewEntry(telescope.EntryTypeScheduledTask, content)

	w.Scope().RecordScheduledTask(entry)
}
