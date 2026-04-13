package mailx

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"

	cmail "github.com/bedrock/packages/contracts/mail"
)

// ArrayTransport stores sent messages in memory. It is intended for
// testing and never actually delivers email.
type ArrayTransport struct {
	mu       sync.Mutex
	messages []*cmail.SentMessage
}

var _ Transport = (*ArrayTransport)(nil)

// NewArrayTransport creates a new in-memory transport.
func NewArrayTransport() *ArrayTransport {
	return &ArrayTransport{}
}

// Send stores the message in memory and returns a SentMessage.
func (t *ArrayTransport) Send(_ context.Context, message *cmail.Message) (*cmail.SentMessage, error) {
	sent := &cmail.SentMessage{
		Original:  message,
		MessageID: generateMessageID(),
	}

	t.mu.Lock()
	t.messages = append(t.messages, sent)
	t.mu.Unlock()

	return sent, nil
}

// String returns the transport name.
func (t *ArrayTransport) String() string {
	return "array"
}

// Messages returns a copy of all sent messages.
func (t *ArrayTransport) Messages() []*cmail.SentMessage {
	t.mu.Lock()

	defer t.mu.Unlock()

	cp := make([]*cmail.SentMessage, len(t.messages))
	copy(cp, t.messages)

	return cp
}

// Flush removes all stored messages.
func (t *ArrayTransport) Flush() {
	t.mu.Lock()
	t.messages = nil
	t.mu.Unlock()
}

func generateMessageID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	return fmt.Sprintf("<%x@mailx>", b)
}
