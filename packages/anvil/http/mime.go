package http

import (
	"mime"
	"path/filepath"
	"strings"
)

// mimeTypes maps file extensions to MIME types as a fallback when
// the standard library's mime package does not recognise an extension.
var mimeTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
	".bmp":  "image/bmp",
	".svg":  "image/svg+xml",
	".ico":  "image/x-icon",
	".tif":  "image/tiff",
	".tiff": "image/tiff",
	".mp4":  "video/mp4",
	".webm": "video/webm",
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
	".ogg":  "audio/ogg",
	".pdf":  "application/pdf",
	".zip":  "application/zip",
	".gz":   "application/gzip",
	".tar":  "application/x-tar",
	".json": "application/json",
	".xml":  "application/xml",
	".csv":  "text/csv",
	".txt":  "text/plain",
	".html": "text/html",
	".htm":  "text/html",
	".css":  "text/css",
	".js":   "application/javascript",
	".wasm": "application/wasm",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
}

// mimeExtensions is the reverse lookup table built from mimeTypes.
// For MIME types that map to multiple extensions the first one wins.
var mimeExtensions map[string]string

func init() {
	mimeExtensions = make(map[string]string, len(mimeTypes))
	for ext, mt := range mimeTypes {
		if _, exists := mimeExtensions[mt]; !exists {
			mimeExtensions[mt] = ext[1:] // strip leading dot
		}
	}
}

// MimeFrom returns the MIME type for the given filename based on its
// extension. Returns "application/octet-stream" when the extension is
// unknown.
func MimeFrom(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return "application/octet-stream"
	}

	return MimeGet(ext[1:])
}

// MimeGet returns the MIME type for the given extension (without the
// leading dot). Returns "application/octet-stream" when the extension
// is unknown.
func MimeGet(ext string) string {
	dotExt := "." + strings.ToLower(ext)

	if mt, ok := mimeTypes[dotExt]; ok {
		return mt
	}

	if mt := mime.TypeByExtension(dotExt); mt != "" {
		// Strip parameters (e.g. "; charset=utf-8").
		if i := strings.IndexByte(mt, ';'); i >= 0 {
			mt = strings.TrimSpace(mt[:i])
		}
		return mt
	}

	return "application/octet-stream"
}

// MimeSearch returns the file extension (without the leading dot) for
// the given MIME type. Returns an empty string when no match is found.
func MimeSearch(mimeType string) string {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))

	if ext, ok := mimeExtensions[mimeType]; ok {
		return ext
	}

	exts, _ := mime.ExtensionsByType(mimeType)
	if len(exts) > 0 {
		return exts[0][1:] // strip leading dot
	}

	return ""
}
