package ai

import (
	"github.com/bedrock/packages/ai/fake"
)

// FakeFiles injects a fake file gateway into the default provider.
func FakeFiles() {
	globalManager().FakeFileProvider()
}

// AssertFileStored fails t if no stored file matches fn (receives filename).
func AssertFileStored(t fake.TestingT, fn func(filename string) bool) {
	globalManager().Recorder().AssertFileStored(t, fn)
}

// AssertNothingFileStored fails t if any file was stored.
func AssertNothingFileStored(t fake.TestingT) {
	globalManager().Recorder().AssertNothingFileStored(t)
}

// AssertFileDeleted fails t if no deleted file matches fn (receives id).
func AssertFileDeleted(t fake.TestingT, fn func(id string) bool) {
	globalManager().Recorder().AssertFileDeleted(t, fn)
}

// AssertNothingFileDeleted fails t if any file was deleted.
func AssertNothingFileDeleted(t fake.TestingT) {
	globalManager().Recorder().AssertNothingFileDeleted(t)
}
