package mailx

import (
	"context"

	cmail "github.com/bedrock/packages/contracts/mail"
)

// Transport defines the contract for email delivery backends.
type Transport interface {
	// Send delivers the message and returns the sent result.
	Send(ctx context.Context, message *cmail.Message) (*cmail.SentMessage, error)
	// String returns the transport name (e.g. "smtp", "log", "array").
	String() string
}
