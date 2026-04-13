package mailx

import (
	"context"
	"fmt"
	"strings"

	clog "github.com/bedrock/packages/contracts/log"
	cmail "github.com/bedrock/packages/contracts/mail"
)

// LogTransport logs email messages instead of delivering them. It is
// intended for development and debugging.
type LogTransport struct {
	logger clog.Logger
}

var _ Transport = (*LogTransport)(nil)

// NewLogTransport creates a transport that logs messages via the given logger.
func NewLogTransport(logger clog.Logger) *LogTransport {
	return &LogTransport{logger: logger}
}

// Send logs the message details and returns a SentMessage.
func (t *LogTransport) Send(_ context.Context, message *cmail.Message) (*cmail.SentMessage, error) {
	t.logger.Debug(formatMessageForLog(message))

	return &cmail.SentMessage{
		Original:  message,
		MessageID: generateMessageID(),
	}, nil
}

// String returns the transport name.
func (t *LogTransport) String() string {
	return "log"
}

// Logger returns the underlying logger.
func (t *LogTransport) Logger() clog.Logger {
	return t.logger
}

func formatMessageForLog(m *cmail.Message) string {
	var b strings.Builder

	b.WriteString("Mail message\n")
	b.WriteString(fmt.Sprintf("  From:    %s\n", m.GetFrom()))
	b.WriteString(fmt.Sprintf("  Subject: %s\n", m.GetSubject()))

	if to := m.GetTo(); len(to) > 0 {
		b.WriteString(fmt.Sprintf("  To:      %s\n", formatAddresses(to)))
	}

	if cc := m.GetCC(); len(cc) > 0 {
		b.WriteString(fmt.Sprintf("  CC:      %s\n", formatAddresses(cc)))
	}

	if bcc := m.GetBCC(); len(bcc) > 0 {
		b.WriteString(fmt.Sprintf("  BCC:     %s\n", formatAddresses(bcc)))
	}

	if html := m.GetHTMLBody(); html != "" {
		b.WriteString(fmt.Sprintf("  HTML:    %s\n", html))
	}

	if text := m.GetTextBody(); text != "" {
		b.WriteString(fmt.Sprintf("  Text:    %s\n", text))
	}

	return b.String()
}

func formatAddresses(addrs []cmail.Address) string {
	parts := make([]string, len(addrs))

	for i, a := range addrs {
		parts[i] = a.String()
	}

	return strings.Join(parts, ", ")
}
