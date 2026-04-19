package responses

import "github.com/bedrock/packages/ai/sdk/data"

// FileResponse holds metadata for a retrieved file.
// Mirrors Laravel\Ai\Responses\FileResponse.
type FileResponse struct {
	ID       string
	Filename string
	Size     int
	Usage    data.Usage
	Meta     data.Meta
}

// StoredFileResponse holds metadata for an uploaded file.
// Mirrors Laravel\Ai\Responses\StoredFileResponse.
type StoredFileResponse struct {
	ID       string
	Filename string
	Usage    data.Usage
	Meta     data.Meta
}

// AddedDocumentResponse holds the result of adding a document to a vector store.
// Mirrors Laravel\Ai\Responses\AddedDocumentResponse.
type AddedDocumentResponse struct {
	ID       string
	Filename string
	Usage    data.Usage
	Meta     data.Meta
}
