package watchers

import (
	"github.com/bedrock/packages/debugbar"
)

// MailMessage carries the data captured by MailWatcher for a single outbound
// email, mirroring the fields stored by Upstream's MailWatcher.
type MailMessage struct {
	// Mailable is the fully-qualified struct/type name of the mailable.
	Mailable string
	// Subject of the email.
	Subject string
	// From is the sender address(es).
	From []string
	// To is the primary recipient address(es).
	To []string
	// CC holds carbon-copy recipients.
	CC []string
	// BCC holds blind carbon-copy recipients.
	BCC []string
	// Queued reports whether the mail was queued rather than sent immediately.
	Queued bool
	// Tags are additional searchable tags.
	Tags []string
}

// MailWatcher monitors outbound email dispatch and records entries as
// DebugBar entries. It mirrors Upstream's MailWatcher class.
type MailWatcher struct {
	debugbar.BaseWatcher
}

// NewMailWatcher creates a MailWatcher with the given options.
func NewMailWatcher(t *debugbar.DebugBar, options map[string]any) *MailWatcher {
	w := &MailWatcher{}
	w.SetDebugBar(t)
	w.Options = options

	return w
}

// Register is a no-op for MailWatcher; callers drive it via Record.
func (w *MailWatcher) Register(_ any) error { return nil }

// Record records an outbound mail entry.
func (w *MailWatcher) Record(msg MailMessage) {
	content := map[string]any{
		"mailable": msg.Mailable,
		"subject":  msg.Subject,
		"from":     msg.From,
		"to":       msg.To,
		"cc":       msg.CC,
		"bcc":      msg.BCC,
		"queued":   msg.Queued,
	}

	entry := debugbar.NewEntry(debugbar.EntryTypeMail, content)
	entry.AddTags(msg.To...)
	entry.AddTags(msg.Tags...)

	w.Scope().RecordMail(entry)
}
