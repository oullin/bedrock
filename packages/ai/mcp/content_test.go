package mcp_test

import (
	"encoding/base64"
	"testing"

	"github.com/bedrock/packages/ai/mcp"
)

// Port of Upstream\Mcp\Tests\ContentTest

func TestTextContentToToolFormat(t *testing.T) {
	t.Parallel()

	resp := mcp.Text("hello world")
	contents := resp.Contents()

	if len(contents) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(contents))
	}

	m := contents[0].ToTool()

	if m["type"] != "text" {
		t.Fatalf("expected type=text, got %v", m["type"])
	}

	if m["text"] != "hello world" {
		t.Fatalf("expected text=hello world, got %v", m["text"])
	}
}

func TestTextContentToPromptMatchesToTool(t *testing.T) {
	t.Parallel()

	resp := mcp.Text("prompt text")
	c := resp.Contents()[0]

	if c.ToTool()["text"] != c.ToPrompt()["text"] {
		t.Fatal("expected ToPrompt to equal ToTool for text content")
	}
}

func TestTextContentToResourceIncludesURI(t *testing.T) {
	t.Parallel()

	resp := mcp.Text("resource content")
	c := resp.Contents()[0]
	m := c.ToResource("file://path/to/resource")

	if m["uri"] != "file://path/to/resource" {
		t.Fatalf("expected uri field, got %v", m["uri"])
	}

	if m["text"] != "resource content" {
		t.Fatalf("expected text field, got %v", m["text"])
	}
}

func TestImageContentToToolFormat(t *testing.T) {
	t.Parallel()

	resp := mcp.Image("base64data==", "image/png")
	m := resp.Contents()[0].ToTool()

	if m["type"] != "image" {
		t.Fatalf("expected type=image, got %v", m["type"])
	}

	if m["data"] != "base64data==" {
		t.Fatalf("expected data, got %v", m["data"])
	}

	if m["mimeType"] != "image/png" {
		t.Fatalf("expected mimeType=image/png, got %v", m["mimeType"])
	}
}

func TestAudioContentToToolFormat(t *testing.T) {
	t.Parallel()

	resp := mcp.Audio("audiodata==", "audio/wav")
	m := resp.Contents()[0].ToTool()

	if m["type"] != "audio" {
		t.Fatalf("expected type=audio, got %v", m["type"])
	}

	if m["mimeType"] != "audio/wav" {
		t.Fatalf("expected mimeType=audio/wav, got %v", m["mimeType"])
	}
}

func TestBlobContentToResourceEncodesBase64(t *testing.T) {
	t.Parallel()

	raw := []byte("binary data")
	resp := mcp.Blob(raw, "application/octet-stream")
	m := resp.Contents()[0].ToResource("file://blob")
	encoded, _ := m["blob"].(string)
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		t.Fatalf("expected valid base64 blob, got error: %v", err)
	}

	if string(decoded) != "binary data" {
		t.Fatalf("expected decoded blob to equal original, got %q", decoded)
	}
}

func TestImageContentToResourceUsesBlob(t *testing.T) {
	t.Parallel()

	resp := mcp.Image("imgdata==", "image/jpeg")
	m := resp.Contents()[0].ToResource("file://img")

	if m["blob"] != "imgdata==" {
		t.Fatalf("expected blob field for image resource, got %v", m["blob"])
	}

	if m["uri"] != "file://img" {
		t.Fatalf("expected uri field, got %v", m["uri"])
	}
}
