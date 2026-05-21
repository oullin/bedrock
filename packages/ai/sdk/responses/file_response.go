package responses

import "github.com/bedrock/packages/ai/sdk/data"

// FileResponse holds metadata for a retrieved file.
type FileResponse struct {
	ID       string
	Filename string
	Size     int
	Usage    data.Usage
	Meta     data.Meta
}

// StoredFileResponse holds metadata for an uploaded file.
type StoredFileResponse struct {
	ID       string
	Filename string
	Usage    data.Usage
	Meta     data.Meta
}

// AddedDocumentResponse holds the result of adding a document to a vector store.
type AddedDocumentResponse struct {
	ID       string
	Filename string
	Usage    data.Usage
	Meta     data.Meta
}
