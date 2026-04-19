package mcp

import (
	"encoding/base64"
	"encoding/json"
)

// cursorPayload is the structure encoded inside a pagination cursor.
type cursorPayload struct {
	Offset int `json:"offset"`
}

// CursorPaginator implements cursor-based pagination identical to Laravel's
// CursorPaginator. The cursor is a base64-encoded JSON object {"offset":N}.
type CursorPaginator struct {
	items   []any
	perPage int
	cursor  string
}

// NewCursorPaginator creates a paginator for items. perPage is the page size
// and cursor is the raw base64 cursor from the client (empty string = first
// page).
func NewCursorPaginator(items []any, perPage int, cursor string) *CursorPaginator {
	return &CursorPaginator{items: items, perPage: perPage, cursor: cursor}
}

// decodeCursor extracts the offset from a base64 cursor string. Returns 0 on
// any error, matching Laravel's silent fallback behaviour.
func decodeCursor(cursor string) int {
	if cursor == "" {
		return 0
	}

	b, err := base64.StdEncoding.DecodeString(cursor)

	if err != nil {
		return 0
	}

	var p cursorPayload

	if err := json.Unmarshal(b, &p); err != nil {
		return 0
	}

	if p.Offset < 0 {
		return 0
	}

	return p.Offset
}

// encodeCursor encodes an offset as a base64 cursor string.
func encodeCursor(offset int) string {
	b, _ := json.Marshal(cursorPayload{Offset: offset})

	return base64.StdEncoding.EncodeToString(b)
}

// Paginate returns a map with the page of items under key and, when more
// items remain, a "nextCursor" entry.
func (p *CursorPaginator) Paginate(key string) map[string]any {
	offset := decodeCursor(p.cursor)

	if offset >= len(p.items) {
		return map[string]any{key: []any{}}
	}

	end := offset + p.perPage
	hasMore := end < len(p.items)

	if end > len(p.items) {
		end = len(p.items)
	}

	page := p.items[offset:end]
	result := map[string]any{key: page}

	if hasMore {
		result["nextCursor"] = encodeCursor(end)
	}

	return result
}
