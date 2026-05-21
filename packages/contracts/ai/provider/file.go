package provider

import (
	"context"

	"github.com/bedrock/packages/contracts/ai/gateway"
)

// FileProvider is the provider-level contract for file management.
type FileProvider interface {
	GetFile(ctx context.Context, id string) (*gateway.FileGetResult, error)
	PutFile(ctx context.Context, file gateway.StorableFile) (*gateway.FilePutResult, error)
	DeleteFile(ctx context.Context, id string) error
	FileGateway() gateway.FileGateway
	UseFileGateway(gw gateway.FileGateway) FileProvider
}
