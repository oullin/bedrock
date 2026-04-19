// Port of Laravel\Ai\Tests\Feature\FileFakeTest
package ai_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/ai"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

// TestFileCanBePut mirrors test_file_can_be_stored.
func TestFileCanBePut(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	provider, err := m.FileProvider()

	if err != nil {
		t.Fatalf("FileProvider error: %v", err)
	}

	file := contractsgw.StorableFile{
		Content:  []byte("hello"),
		Filename: "hello.txt",
		MimeType: "text/plain",
	}
	result, putErr := provider.PutFile(context.Background(), file)

	if putErr != nil {
		t.Fatalf("PutFile error: %v", putErr)
	}

	if result == nil {
		t.Error("expected non-nil file put result")
	}

	rec.AssertFileStored(t, func(filename string) bool {
		return filename == "hello.txt"
	})
}

// TestFileAssertNothingStored mirrors the "nothing stored" assertion.
func TestFileAssertNothingStored(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	rec.AssertNothingFileStored(t)
}

// TestFileCanBeDeleted mirrors test_file_can_be_deleted.
func TestFileCanBeDeleted(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	provider, err := m.FileProvider()

	if err != nil {
		t.Fatalf("FileProvider error: %v", err)
	}

	delErr := provider.DeleteFile(context.Background(), "file-abc-123")

	if delErr != nil {
		t.Fatalf("DeleteFile error: %v", delErr)
	}

	rec.AssertFileDeleted(t, func(id string) bool {
		return id == "file-abc-123"
	})
}

// TestFileAssertNothingDeleted mirrors the "nothing deleted" assertion.
func TestFileAssertNothingDeleted(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	rec.AssertNothingFileDeleted(t)
}
