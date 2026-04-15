package provider

import (
	"context"

	"github.com/bedrock/packages/contracts/ai/gateway"
)

// StoreProvider is the provider-level contract for vector store management.
// Mirrors Upstream\Ai\Contracts\Providers\StoreProvider.
type StoreProvider interface {
	GetStore(ctx context.Context, id string) (*gateway.StoreData, error)
	CreateStore(ctx context.Context, req gateway.StoreCreateRequest) (*gateway.StoreData, error)
	AddFileToStore(ctx context.Context, storeID string, file gateway.StorableFile, metadata map[string]any) (*gateway.AddedDocumentData, error)
	RemoveFileFromStore(ctx context.Context, storeID, fileID string) error
	DeleteStore(ctx context.Context, id string) error
	StoreGateway() gateway.StoreGateway
	UseStoreGateway(gw gateway.StoreGateway) StoreProvider
}
