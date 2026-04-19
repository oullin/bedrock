package ai

import (
	"github.com/bedrock/packages/ai/sdk/fake"
)

// FakeStore injects a fake store gateway into the default provider.
func FakeStore() {
	globalManager().FakeStoreProvider()
}

// AssertStoreFileAdded fails t if no add_file operation matches fn (receives fileID).
func AssertStoreFileAdded(t fake.TestingT, fn func(fileID string) bool) {
	globalManager().Recorder().AssertStoreFileAdded(t, fn)
}

// AssertStoreFileRemoved fails t if no remove_file operation matches fn (receives fileID).
func AssertStoreFileRemoved(t fake.TestingT, fn func(fileID string) bool) {
	globalManager().Recorder().AssertStoreFileRemoved(t, fn)
}
