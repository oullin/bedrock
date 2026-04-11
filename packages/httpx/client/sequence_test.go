package client_test

import (
	"testing"

	"github.com/bedrock/packages/httpx/client"
)

func TestResponseSequence(t *testing.T) {
	t.Parallel()

	seq := client.NewResponseSequence(
		client.ResponseStub{Status: 200, Body: "first"},
		client.ResponseStub{Status: 201, Body: "second"},
	)

	resp1 := client.NewResponse(seq.Next())

	if resp1.Body() != "first" {
		t.Fatalf("expected first, got %s", resp1.Body())
	}

	resp2 := client.NewResponse(seq.Next())

	if resp2.Body() != "second" {
		t.Fatalf("expected second, got %s", resp2.Body())
	}

	if !seq.IsEmpty() {
		t.Fatal("expected sequence to be empty")
	}
}

func TestResponseSequenceWhenEmpty(t *testing.T) {
	t.Parallel()

	seq := client.NewResponseSequence(
		client.ResponseStub{Status: 200, Body: "only"},
	).WhenEmpty(client.ResponseStub{Status: 500, Body: "fallback"})

	_ = seq.Next() // consume "only"

	resp := client.NewResponse(seq.Next())

	if resp.Status() != 500 {
		t.Fatalf("expected 500 fallback, got %d", resp.Status())
	}

	if resp.Body() != "fallback" {
		t.Fatalf("expected fallback, got %s", resp.Body())
	}
}

func TestResponseSequencePush(t *testing.T) {
	t.Parallel()

	seq := client.NewResponseSequence().Push(
		client.ResponseStub{Status: 200, Body: "pushed"},
	)

	resp := client.NewResponse(seq.Next())

	if resp.Body() != "pushed" {
		t.Fatalf("expected pushed, got %s", resp.Body())
	}
}

func TestResponseSequenceDefaultEmpty(t *testing.T) {
	t.Parallel()

	seq := client.NewResponseSequence()

	resp := client.NewResponse(seq.Next())

	if resp.Status() != 200 {
		t.Fatalf("expected 200 default, got %d", resp.Status())
	}
}
