package pagination

import (
	"net/http"
	"strconv"

	cpagination "github.com/bedrock/packages/contracts/pagination"
)

// ResolveCurrentPage extracts the current page number from the request query
// string. Returns 1 if the parameter is absent or invalid.
func ResolveCurrentPage(r *http.Request, pageName ...string) int {
	name := "page"

	if len(pageName) > 0 {
		name = pageName[0]
	}

	page, err := strconv.Atoi(r.URL.Query().Get(name))

	if err != nil || page < 1 {
		return 1
	}

	return page
}

// ResolveCurrentPath extracts the request path.
func ResolveCurrentPath(r *http.Request) string {
	return r.URL.Path
}

// ResolveCurrentCursor decodes the cursor from the request query string.
// Returns nil, nil if the parameter is absent.
func ResolveCurrentCursor(r *http.Request, cursorName ...string) (*cpagination.Cursor, error) {
	name := "cursor"

	if len(cursorName) > 0 {
		name = cursorName[0]
	}

	encoded := r.URL.Query().Get(name)

	if encoded == "" {
		return nil, nil
	}

	return cpagination.DecodeCursor(encoded)
}
