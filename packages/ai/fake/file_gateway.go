package fake

import (
	"context"
	"fmt"
	"sync"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

// FileGateway is a fake implementation of gateway.FileGateway.
type FileGateway struct {
	mu       sync.Mutex
	prevent  bool
	recorder *Recorder
}

var _ contractsgw.FileGateway = (*FileGateway)(nil)

// NewFileGateway creates a FileGateway with a shared Recorder.
func NewFileGateway(recorder *Recorder) *FileGateway {
	return &FileGateway{recorder: recorder}
}

// PreventStray enables stray-call prevention.
func (g *FileGateway) PreventStray() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.prevent = true
}

// GetFile satisfies gateway.FileGateway.
func (g *FileGateway) GetFile(ctx context.Context, id string) (*contractsgw.FileGetResult, error) {
	g.recorder.recordFile("get", id, "")
	return &contractsgw.FileGetResult{ID: id, Filename: "fake_file.txt"}, nil
}

// PutFile satisfies gateway.FileGateway.
func (g *FileGateway) PutFile(ctx context.Context, file contractsgw.StorableFile) (*contractsgw.FilePutResult, error) {
	id := FakeFileID()
	g.recorder.recordFile("put", id, file.Filename)
	return &contractsgw.FilePutResult{ID: id, Filename: file.Filename}, nil
}

// DeleteFile satisfies gateway.FileGateway.
func (g *FileGateway) DeleteFile(ctx context.Context, id string) error {
	g.mu.Lock()
	prevent := g.prevent
	g.mu.Unlock()

	if prevent && id == "" {
		return fmt.Errorf("ai: unexpected call to faked file gateway")
	}
	g.recorder.recordFile("delete", id, "")
	return nil
}
