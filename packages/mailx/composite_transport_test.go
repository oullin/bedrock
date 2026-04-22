package mailx

import (
	"context"
	"errors"
	"testing"

	cmail "github.com/bedrock/packages/contracts/mail"
)

type recordingTransport struct {
	name  string
	err   error
	calls int
}

type fakeAPIEmailClient struct {
	payload APIEmailPayload
}

func (t *recordingTransport) Send(_ context.Context, message *cmail.Message) (*cmail.SentMessage, error) {
	t.calls++

	if t.err != nil {
		return nil, t.err
	}

	return &cmail.SentMessage{Original: message, MessageID: t.name}, nil
}

func (t *recordingTransport) String() string {
	return t.name
}

func TestFailoverTransportUsesNextTransportAfterFailure(t *testing.T) {
	t.Parallel()

	first := &recordingTransport{name: "first", err: errors.New("down")}
	second := &recordingTransport{name: "second"}

	sent, err := NewFailoverTransport(first, second).Send(context.Background(), cmail.NewMessage())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sent.MessageID != "second" || first.calls != 1 || second.calls != 1 {
		t.Fatalf("expected failover to second transport, got id=%q first=%d second=%d", sent.MessageID, first.calls, second.calls)
	}
}

func TestRoundRobinTransportRotatesTransports(t *testing.T) {
	t.Parallel()

	first := &recordingTransport{name: "first"}
	second := &recordingTransport{name: "second"}
	transport := NewRoundRobinTransport(first, second)

	a, err := transport.Send(context.Background(), cmail.NewMessage())

	if err != nil {
		t.Fatalf("unexpected first error: %v", err)
	}

	b, err := transport.Send(context.Background(), cmail.NewMessage())

	if err != nil {
		t.Fatalf("unexpected second error: %v", err)
	}

	if a.MessageID != "first" || b.MessageID != "second" {
		t.Fatalf("expected rotation, got %q then %q", a.MessageID, b.MessageID)
	}
}

func (c *fakeAPIEmailClient) SendEmail(_ context.Context, payload APIEmailPayload) (*APIEmailResult, error) {
	c.payload = payload

	return &APIEmailResult{MessageID: "api-id"}, nil
}

func TestAPITransportsUseInjectedClients(t *testing.T) {
	t.Parallel()

	client := &fakeAPIEmailClient{}
	msg := cmail.NewMessage().Subject("Hello")

	sent, err := NewSESV2Transport(client).Send(context.Background(), msg)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sent.MessageID != "api-id" || client.payload.Message != msg {
		t.Fatal("expected injected API client to receive message and return id")
	}
}
