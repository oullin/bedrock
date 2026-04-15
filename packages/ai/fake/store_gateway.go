package fake

import (
	"context"
	"fmt"
	"sync"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

// StoreGateway is a fake implementation of gateway.StoreGateway.
type StoreGateway struct {
	mu       sync.Mutex
	prevent  bool
	recorder *Recorder
}

var _ contractsgw.StoreGateway = (*StoreGateway)(nil)

// NewStoreGateway creates a StoreGateway with a shared Recorder.
func NewStoreGateway(recorder *Recorder) *StoreGateway {
	return &StoreGateway{recorder: recorder}
}

// PreventStray enables stray-call prevention.
func (g *StoreGateway) PreventStray() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.prevent = true
}

// GetStore satisfies gateway.StoreGateway.
func (g *StoreGateway) GetStore(ctx context.Context, id string) (*contractsgw.StoreData, error) {
	g.recorder.recordStore("get", id, "")
	return &contractsgw.StoreData{ID: id, Ready: true}, nil
}

// CreateStore satisfies gateway.StoreGateway.
func (g *StoreGateway) CreateStore(ctx context.Context, req contractsgw.StoreCreateRequest) (*contractsgw.StoreData, error) {
	id := fmt.Sprintf("store_%d", fakeIDCounter.Add(1))
	g.recorder.recordStore("create", id, "")
	return &contractsgw.StoreData{ID: id, Ready: true}, nil
}

// AddFile satisfies gateway.StoreGateway.
func (g *StoreGateway) AddFile(ctx context.Context, storeID string, file contractsgw.StorableFile, metadata map[string]any) (*contractsgw.AddedDocumentData, error) {
	docID := FakeFileID()
	g.recorder.recordStore("add_file", storeID, docID)
	return &contractsgw.AddedDocumentData{ID: docID, Filename: file.Filename}, nil
}

// RemoveFile satisfies gateway.StoreGateway.
func (g *StoreGateway) RemoveFile(ctx context.Context, storeID, fileID string) error {
	g.recorder.recordStore("remove_file", storeID, fileID)
	return nil
}

// DeleteStore satisfies gateway.StoreGateway.
func (g *StoreGateway) DeleteStore(ctx context.Context, id string) error {
	g.recorder.recordStore("delete", id, "")
	return nil
}
