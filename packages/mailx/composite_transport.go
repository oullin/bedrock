package mailx

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	cmail "github.com/bedrock/packages/contracts/mail"
)

// FailoverTransport attempts each transport until one succeeds.
type FailoverTransport struct {
	transports []Transport
}

// NewFailoverTransport creates a failover transport.

// Send sends through the first successful transport.

// String returns the transport name.

// RoundRobinTransport sends each message through the next transport.
type RoundRobinTransport struct {
	transports []Transport
	next       atomic.Uint64
}

func NewFailoverTransport(transports ...Transport) *FailoverTransport {
	return &FailoverTransport{transports: append([]Transport(nil), transports...)}
}

func (t *FailoverTransport) Send(ctx context.Context, message *cmail.Message) (*cmail.SentMessage, error) {
	if len(t.transports) == 0 {
		return nil, errors.New("mailx: failover transport has no transports")
	}

	var errs []error

	for _, transport := range t.transports {
		sent, err := transport.Send(ctx, message)

		if err == nil {
			return sent, nil
		}

		errs = append(errs, fmt.Errorf("%s: %w", transport.String(), err))
	}

	return nil, errors.Join(errs...)
}

func (t *FailoverTransport) String() string {
	return "failover"
}

// NewRoundRobinTransport creates a round-robin transport.
func NewRoundRobinTransport(transports ...Transport) *RoundRobinTransport {
	return &RoundRobinTransport{transports: append([]Transport(nil), transports...)}
}

// Send sends through the next transport.
func (t *RoundRobinTransport) Send(ctx context.Context, message *cmail.Message) (*cmail.SentMessage, error) {
	if len(t.transports) == 0 {
		return nil, errors.New("mailx: round-robin transport has no transports")
	}

	index := int(t.next.Add(1)-1) % len(t.transports)

	return t.transports[index].Send(ctx, message)
}

// String returns the transport name.
func (t *RoundRobinTransport) String() string {
	return "round-robin"
}
