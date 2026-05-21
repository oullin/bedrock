package gateway

import "context"

// StoreCreateRequest carries parameters for creating a vector store.
type StoreCreateRequest struct {
	Name        *string
	Description *string
	Files       []StorableFile
	ExpiresIn   *int // seconds
}

// StoreData holds metadata for a vector store.
type StoreData struct {
	ID                    string
	Name                  *string
	FilesCount            int
	VectorStoreFilesCount int
	Ready                 bool
}

// AddedDocumentData holds the result of adding a file to a store.
type AddedDocumentData struct {
	ID       string
	Filename string
	Usage    TokenUsage
	Meta     ResponseMeta
}

// StoreGateway manages vector stores for a single provider.
type StoreGateway interface {
	GetStore(ctx context.Context, id string) (*StoreData, error)
	CreateStore(ctx context.Context, req StoreCreateRequest) (*StoreData, error)
	AddFile(ctx context.Context, storeID string, file StorableFile, metadata map[string]any) (*AddedDocumentData, error)
	RemoveFile(ctx context.Context, storeID, fileID string) error
	DeleteStore(ctx context.Context, id string) error
}
