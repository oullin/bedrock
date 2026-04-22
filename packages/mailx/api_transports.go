package mailx

import (
	"context"

	cmail "github.com/bedrock/packages/contracts/mail"
)

// APIEmailPayload is the provider-neutral email payload used by API transports.
type APIEmailPayload struct {
	Message *cmail.Message
}

// APIEmailResult is the provider-neutral result returned by API transports.
type APIEmailResult struct {
	MessageID string
}

// APIEmailClient sends an email through a provider API.
type APIEmailClient interface {
	SendEmail(ctx context.Context, payload APIEmailPayload) (*APIEmailResult, error)
}

type apiTransport struct {
	name   string
	client APIEmailClient
}

// NewSESTransport creates a typed AWS SES API transport.
func NewSESTransport(client APIEmailClient) Transport {
	return &apiTransport{name: "ses", client: client}
}

// NewSESV2Transport creates a typed AWS SESv2 API transport.
func NewSESV2Transport(client APIEmailClient) Transport {
	return &apiTransport{name: "ses-v2", client: client}
}

// NewCloudflareTransport creates a typed Cloudflare email API transport.
func NewCloudflareTransport(client APIEmailClient) Transport {
	return &apiTransport{name: "cloudflare", client: client}
}

func (t *apiTransport) Send(ctx context.Context, message *cmail.Message) (*cmail.SentMessage, error) {
	result, err := t.client.SendEmail(ctx, APIEmailPayload{Message: message})

	if err != nil {
		return nil, err
	}

	messageID := ""

	if result != nil {
		messageID = result.MessageID
	}

	return &cmail.SentMessage{Original: message, MessageID: messageID}, nil
}

func (t *apiTransport) String() string {
	return t.name
}
