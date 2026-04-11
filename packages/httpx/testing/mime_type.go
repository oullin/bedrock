package testing

import (
	"path/filepath"
	"strings"
)

// mimeTypes maps file extensions to MIME types.
var mimeTypes = map[string]string{
	"jpg":   "image/jpeg",
	"jpeg":  "image/jpeg",
	"png":   "image/png",
	"gif":   "image/gif",
	"bmp":   "image/bmp",
	"svg":   "image/svg+xml",
	"webp":  "image/webp",
	"ico":   "image/x-icon",
	"tif":   "image/tiff",
	"tiff":  "image/tiff",
	"pdf":   "application/pdf",
	"csv":   "text/csv",
	"txt":   "text/plain",
	"html":  "text/html",
	"htm":   "text/html",
	"json":  "application/json",
	"xml":   "application/xml",
	"zip":   "application/zip",
	"gz":    "application/gzip",
	"tar":   "application/x-tar",
	"mp3":   "audio/mpeg",
	"mp4":   "video/mp4",
	"avi":   "video/x-msvideo",
	"doc":   "application/msword",
	"docx":  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"xls":   "application/vnd.ms-excel",
	"xlsx":  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"ppt":   "application/vnd.ms-powerpoint",
	"pptx":  "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"css":   "text/css",
	"js":    "application/javascript",
	"woff":  "font/woff",
	"woff2": "font/woff2",
	"ttf":   "font/ttf",
	"otf":   "font/otf",
}

// MimeType returns the MIME type for a file name or extension. Returns
// "application/octet-stream" for unknown types.
func MimeType(nameOrExt string) string {
	ext := strings.TrimPrefix(filepath.Ext(nameOrExt), ".")

	if ext == "" {
		ext = strings.ToLower(nameOrExt)
	}

	if mt, ok := mimeTypes[strings.ToLower(ext)]; ok {
		return mt
	}

	return "application/octet-stream"
}

// ExtensionForMime returns the first matching extension for a MIME type.
// Returns an empty string if no match is found.
func ExtensionForMime(mimeType string) string {
	mimeType = strings.ToLower(mimeType)

	for ext, mt := range mimeTypes {
		if mt == mimeType {
			return ext
		}
	}

	return ""
}
