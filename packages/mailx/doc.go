// Package mailx provides driver-based email sending with support for
// SMTP, log, and array (testing) transports. It mirrors the upstream Mail
// component, offering a unified API through the MailManager and
// individual mailers for each transport type. The package supports
// rich message construction including HTML and plain-text bodies,
// file attachments, inline embeds, custom headers, metadata, and
// tags. Events are dispatched before and after sending for
// observability and interception.
package mailx
