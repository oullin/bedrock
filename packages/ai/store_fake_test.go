// Port of Upstream\Ai\Tests\Feature\StoreFakeTest
package ai_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/ai"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

// TestStoreFileCanBeAdded mirrors test_file_can_be_added_to_store.
func TestStoreFileCanBeAdded(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	provider, err := m.StoreProvider()
	if err != nil {
		t.Fatalf("StoreProvider error: %v", err)
	}

	storeID := "store-abc"
	file := contractsgw.StorableFile{
		Content:  []byte("content"),
		Filename: "data.txt",
		MimeType: "text/plain",
	}
	result, addErr := provider.AddFileToStore(context.Background(), storeID, file, nil)
	if addErr != nil {
		t.Fatalf("AddFileToStore error: %v", addErr)
	}
	if result == nil {
		t.Error("expected non-nil add result")
	}

	rec.AssertStoreFileAdded(t, func(fileID string) bool {
		return fileID != "" // any non-empty file ID is valid
	})
}

// TestStoreFileCanBeRemoved mirrors test_file_can_be_removed_from_store.
func TestStoreFileCanBeRemoved(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	provider, err := m.StoreProvider()
	if err != nil {
		t.Fatalf("StoreProvider error: %v", err)
	}

	removeErr := provider.RemoveFileFromStore(context.Background(), "store-abc", "file-xyz")
	if removeErr != nil {
		t.Fatalf("RemoveFileFromStore error: %v", removeErr)
	}

	rec.AssertStoreFileRemoved(t, func(fileID string) bool {
		return fileID == "file-xyz"
	})
}
