// FileFakeTest::test_files_can_be_faked
// FileFakeTest::test_can_assert_file_was_stored
// FileFakeTest::test_can_assert_no_files_were_stored
// FileFakeTest::test_can_assert_file_was_deleted
// FileFakeTest::test_can_assert_no_files_were_deleted
package ai_test

import (
	"context"
	"testing"

	ai "github.com/bedrock/packages/ai/sdk"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

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

func TestFileAssertNothingStored(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	rec.AssertNothingFileStored(t)
}

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

func TestFileAssertNothingDeleted(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	rec.AssertNothingFileDeleted(t)
}
