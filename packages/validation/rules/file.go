package rules

import (
	"mime/multipart"
	"path/filepath"
	"strings"
)

func init_file() {
	Register("File", validateFile)
	Register("Image", validateImage)
	Register("Mimes", validateMimes)
	Register("Mimetypes", validateMimetypes)
	Register("Extensions", validateExtensions)
}

// File validation works with *multipart.FileHeader values.
// When the value is not a file, the rule fails.

func validateFile(_ string, value any, _ []string, _ RuleContext) bool {
	_, ok := value.(*multipart.FileHeader)

	return ok
}

var imageExts = map[string]bool{
	"jpg": true, "jpeg": true, "png": true, "gif": true,
	"bmp": true, "svg": true, "webp": true, "avif": true,
}

func validateImage(_ string, value any, _ []string, _ RuleContext) bool {
	fh, ok := value.(*multipart.FileHeader)

	if !ok {
		return false
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fh.Filename), "."))

	return imageExts[ext]
}

// validateMimes: file extension must match one of the given MIME type groups.
// Params: [jpeg, png, pdf, ...]  (extensions, not full MIME types)
func validateMimes(_ string, value any, params []string, _ RuleContext) bool {
	fh, ok := value.(*multipart.FileHeader)

	if !ok {
		return false
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fh.Filename), "."))

	return containsString(params, ext)
}

// validateMimetypes: the header's content type must match exactly.
// Params: [image/jpeg, application/pdf, ...]
func validateMimetypes(_ string, value any, params []string, _ RuleContext) bool {
	fh, ok := value.(*multipart.FileHeader)

	if !ok {
		return false
	}

	ct := fh.Header.Get("Content-Type")

	for _, p := range params {
		if strings.EqualFold(ct, p) {
			return true
		}
	}

	return false
}

// validateExtensions: file must have one of the allowed extensions.
func validateExtensions(_ string, value any, params []string, _ RuleContext) bool {
	fh, ok := value.(*multipart.FileHeader)

	if !ok {
		return false
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fh.Filename), "."))

	return containsString(params, ext)
}
