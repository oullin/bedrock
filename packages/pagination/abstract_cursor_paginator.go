package pagination

import (
	"net/url"
	"sort"
	"strings"

	cpagination "github.com/bedrock/packages/contracts/pagination"
)

// abstractCursorPaginator provides shared functionality for cursor-based paginators.
type abstractCursorPaginator[T any] struct {
	items      []T
	perPage    int
	cursor     *cpagination.Cursor
	path       string
	query      map[string]string
	fragment   string
	cursorName string
	hasMore    bool
	options    map[string]any
}

func (p *abstractCursorPaginator[T]) buildUrl(cursor *cpagination.Cursor) string {
	path := p.path

	params := url.Values{}
	keys := make([]string, 0, len(p.query))

	for k := range p.query {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		params.Set(k, p.query[k])
	}

	if cursor != nil {
		params.Set(p.cursorName, cursor.Encode())
	}

	result := path

	if len(params) > 0 {
		result += "?" + encodeQuery(params)
	}

	if p.fragment != "" {
		result += "#" + p.fragment
	}

	return result
}

func (p *abstractCursorPaginator[T]) appendsQuery(query map[string]string) {
	for k, v := range query {
		p.query[k] = v
	}
}

func (p *abstractCursorPaginator[T]) getFragment() string {
	return p.fragment
}

func (p *abstractCursorPaginator[T]) setFragment(f string) {
	p.fragment = f
}

func (p *abstractCursorPaginator[T]) getItems() []T {
	return p.items
}

func (p *abstractCursorPaginator[T]) transformItems(callback func(T) T) {
	for i, item := range p.items {
		p.items[i] = callback(item)
	}
}

func (p *abstractCursorPaginator[T]) count() int {
	return len(p.items)
}

func (p *abstractCursorPaginator[T]) isEmpty() bool {
	return len(p.items) == 0
}

func (p *abstractCursorPaginator[T]) isNotEmpty() bool {
	return len(p.items) > 0
}

func (p *abstractCursorPaginator[T]) getPerPage() int {
	return p.perPage
}

func (p *abstractCursorPaginator[T]) getCursor() *cpagination.Cursor {
	return p.cursor
}

func (p *abstractCursorPaginator[T]) getCursorName() string {
	return p.cursorName
}

func (p *abstractCursorPaginator[T]) setCursorName(name string) {
	p.cursorName = name
}

func (p *abstractCursorPaginator[T]) getPath() string {
	return p.path
}

func (p *abstractCursorPaginator[T]) setPath(path string) {
	p.path = strings.TrimRight(path, "/")
}

func (p *abstractCursorPaginator[T]) getOptions() map[string]any {
	return p.options
}

func (p *abstractCursorPaginator[T]) itemsToAny() []any {
	result := make([]any, len(p.items))

	for i, item := range p.items {
		result[i] = item
	}

	return result
}
