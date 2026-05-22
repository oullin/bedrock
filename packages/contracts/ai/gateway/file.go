package gateway

import "context"

// StorableFile represents a file that can be uploaded to a provider.
type StorableFile struct {
	Content  []byte
	Filename string
	MimeType string
}

// FileGetResult holds metadata for a retrieved file.
type FileGetResult struct {
	ID       string
	Filename string
	Size     int
	Usage    TokenUsage
	Meta     ResponseMeta
}

// FilePutResult holds the provider-assigned ID for an uploaded file.
type FilePutResult struct {
	ID       string
	Filename string
	Usage    TokenUsage
	Meta     ResponseMeta
}

// FileGateway manages file uploads and retrieval for a single provider.
type FileGateway interface {
	GetFile(ctx context.Context, id string) (*FileGetResult, error)
	PutFile(ctx context.Context, file StorableFile) (*FilePutResult, error)
	DeleteFile(ctx context.Context, id string) error
}
