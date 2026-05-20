package mcp_test

import (
	"encoding/base64"
	"testing"

	"github.com/bedrock/packages/ai/mcp"
)

// Port of \Mcp\Tests\ContentTest

func TestTextContentToToolFormat(t *testing.T) {
	t.Parallel()

	// TextTest::it_does_not_include_meta_if_null
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

	if _, ok := m["_meta"]; ok {
		t.Fatalf("expected no _meta when text metadata is unset, got %#v", m)
	}
}

func TestTextContentStringReturnsRawText(t *testing.T) {
	t.Parallel()

	// TextTest::it_casts_to_string_as_raw_text
	content := mcp.Text("hello world").Contents()[0]
	stringer, ok := content.(interface{ String() string })

	if !ok {
		t.Fatalf("expected text content to implement String")
	}

	if stringer.String() != "hello world" {
		t.Fatalf("expected raw text string, got %q", stringer.String())
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

	// ImageTest::it_does_not_include_meta_if_null
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

	if _, ok := m["_meta"]; ok {
		t.Fatalf("expected no _meta when image metadata is unset, got %#v", m)
	}
}

func TestImageContentDefaultsAndString(t *testing.T) {
	t.Parallel()

	// ImageTest::it_casts_to_string_as_raw_data
	// ImageTest::it_defaults_mimetype_to_image_png
	content := mcp.Image("base64data==", "").Contents()[0]
	stringer, ok := content.(interface{ String() string })

	if !ok {
		t.Fatalf("expected image content to implement String")
	}

	if stringer.String() != "base64data==" {
		t.Fatalf("expected raw image data string, got %q", stringer.String())
	}

	if got := content.ToTool()["mimeType"]; got != "image/png" {
		t.Fatalf("expected default image/png MIME type, got %v", got)
	}
}

func TestAudioContentToToolFormat(t *testing.T) {
	t.Parallel()

	// AudioTest::it_does_not_include_meta_if_null
	resp := mcp.Audio("audiodata==", "audio/wav")
	m := resp.Contents()[0].ToTool()

	if m["type"] != "audio" {
		t.Fatalf("expected type=audio, got %v", m["type"])
	}

	if m["mimeType"] != "audio/wav" {
		t.Fatalf("expected mimeType=audio/wav, got %v", m["mimeType"])
	}

	if _, ok := m["_meta"]; ok {
		t.Fatalf("expected no _meta when audio metadata is unset, got %#v", m)
	}
}

func TestAudioContentDefaultsAndString(t *testing.T) {
	t.Parallel()

	// AudioTest::it_casts_to_string_as_raw_data
	// AudioTest::it_defaults_mimetype_to_audio_wav
	content := mcp.Audio("audiodata==", "").Contents()[0]
	stringer, ok := content.(interface{ String() string })

	if !ok {
		t.Fatalf("expected audio content to implement String")
	}

	if stringer.String() != "audiodata==" {
		t.Fatalf("expected raw audio data string, got %q", stringer.String())
	}

	if got := content.ToTool()["mimeType"]; got != "audio/wav" {
		t.Fatalf("expected default audio/wav MIME type, got %v", got)
	}
}

func TestBlobContentToResourceEncodesBase64(t *testing.T) {
	t.Parallel()

	// BlobTest::it_casts_to_string_as_raw_content
	// BlobTest::it_does_not_include_meta_if_null
	raw := []byte("binary data")
	resp := mcp.Blob(raw, "application/octet-stream")
	content := resp.Contents()[0]
	stringer, ok := content.(interface{ String() string })

	if !ok {
		t.Fatalf("expected blob content to implement String")
	}

	if stringer.String() != "binary data" {
		t.Fatalf("expected raw blob string, got %q", stringer.String())
	}

	m := content.ToResource("file://blob")
	encoded, _ := m["blob"].(string)
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		t.Fatalf("expected valid base64 blob, got error: %v", err)
	}

	if string(decoded) != "binary data" {
		t.Fatalf("expected decoded blob to equal original, got %q", decoded)
	}

	if _, ok := m["_meta"]; ok {
		t.Fatalf("expected no _meta when blob metadata is unset, got %#v", m)
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
