package pagination

import (
	"encoding/json"

	cpagination "github.com/bedrock/packages/contracts/pagination"
)

// CursorPaginator is a paginator that uses cursor-based pagination for
// efficient traversal of large datasets.
type CursorPaginator[T any] struct {
	abstractCursorPaginator[T]
	nextCursorValue *cpagination.Cursor
	prevCursorValue *cpagination.Cursor
}

// NewCursorPaginator creates a new cursor-based paginator. If items contains
// more than perPage elements, the extra items are trimmed and hasMore is set to
// true. Options may include: "path" (string), "cursorName" (string), "query"
// (map[string]string), "fragment" (string).
func NewCursorPaginator[T any](items []T, perPage int, cursor *cpagination.Cursor, options ...map[string]any) *CursorPaginator[T] {
	opts := mergeOptions(options)

	p := &CursorPaginator[T]{
		abstractCursorPaginator: abstractCursorPaginator[T]{
			perPage:    perPage,
			cursor:     cursor,
			cursorName: optString(opts, "cursorName", "cursor"),
			query:      optQuery(opts),
			fragment:   optString(opts, "fragment", ""),
			options:    opts,
		},
	}

	p.setPath(optString(opts, "path", "/"))

	hasMore := len(items) > perPage

	if hasMore {
		items = items[:perPage]
	}

	p.items = items
	p.hasMore = hasMore

	return p
}

func (p *CursorPaginator[T]) Url(cursor *cpagination.Cursor) string {
	return p.buildUrl(cursor)
}

func (p *CursorPaginator[T]) PreviousPageUrl() string {
	if p.prevCursorValue != nil {
		return p.buildUrl(p.prevCursorValue)
	}

	return ""
}

func (p *CursorPaginator[T]) NextPageUrl() string {
	if p.nextCursorValue != nil {
		return p.buildUrl(p.nextCursorValue)
	}

	return ""
}

func (p *CursorPaginator[T]) HasMorePages() bool {
	return p.hasMore
}

func (p *CursorPaginator[T]) HasPages() bool {
	return p.isNotEmpty() || p.cursor != nil
}

func (p *CursorPaginator[T]) OnFirstPage() bool {
	return p.cursor == nil
}

func (p *CursorPaginator[T]) OnLastPage() bool {
	return !p.hasMore
}

func (p *CursorPaginator[T]) Appends(query map[string]string) *CursorPaginator[T] {
	p.appendsQuery(query)

	return p
}

func (p *CursorPaginator[T]) WithQueryString(query map[string]string) *CursorPaginator[T] {
	p.appendsQuery(query)

	return p
}

func (p *CursorPaginator[T]) Fragment(fragment ...string) string {
	if len(fragment) > 0 {
		p.setFragment(fragment[0])
	}

	return p.getFragment()
}

func (p *CursorPaginator[T]) Items() []any {
	return p.itemsToAny()
}

func (p *CursorPaginator[T]) TypedItems() []T {
	return p.getItems()
}

func (p *CursorPaginator[T]) Through(callback func(T) T) *CursorPaginator[T] {
	p.transformItems(callback)

	return p
}

func (p *CursorPaginator[T]) PerPage() int {
	return p.getPerPage()
}

func (p *CursorPaginator[T]) Cursor() *cpagination.Cursor {
	return p.getCursor()
}

func (p *CursorPaginator[T]) NextCursor() *cpagination.Cursor {
	return p.nextCursorValue
}

func (p *CursorPaginator[T]) PreviousCursor() *cpagination.Cursor {
	return p.prevCursorValue
}

func (p *CursorPaginator[T]) SetNextCursor(cursor *cpagination.Cursor) *CursorPaginator[T] {
	p.nextCursorValue = cursor

	return p
}

func (p *CursorPaginator[T]) SetPreviousCursor(cursor *cpagination.Cursor) *CursorPaginator[T] {
	p.prevCursorValue = cursor

	return p
}

func (p *CursorPaginator[T]) Path() string {
	return p.getPath()
}

func (p *CursorPaginator[T]) IsEmpty() bool {
	return p.isEmpty()
}

func (p *CursorPaginator[T]) IsNotEmpty() bool {
	return p.isNotEmpty()
}

func (p *CursorPaginator[T]) Count() int {
	return p.count()
}

func (p *CursorPaginator[T]) GetCursorName() string {
	return p.getCursorName()
}

func (p *CursorPaginator[T]) SetCursorName(name string) {
	p.setCursorName(name)
}

func (p *CursorPaginator[T]) WithPath(path string) *CursorPaginator[T] {
	p.setPath(path)

	return p
}

func (p *CursorPaginator[T]) SetPath(path string) *CursorPaginator[T] {
	p.setPath(path)

	return p
}

func (p *CursorPaginator[T]) GetOptions() map[string]any {
	return p.getOptions()
}

func (p *CursorPaginator[T]) ToMap() map[string]any {
	var nextCursor any

	if p.nextCursorValue != nil {
		nextCursor = p.nextCursorValue.Encode()
	}

	var prevCursor any

	if p.prevCursorValue != nil {
		prevCursor = p.prevCursorValue.Encode()
	}

	return map[string]any{
		"data":          p.items,
		"path":          p.path,
		"per_page":      p.perPage,
		"next_cursor":   nextCursor,
		"next_page_url": nilIfEmpty(p.NextPageUrl()),
		"prev_cursor":   prevCursor,
		"prev_page_url": nilIfEmpty(p.PreviousPageUrl()),
	}
}

func (p *CursorPaginator[T]) ToJSON() ([]byte, error) {
	return json.Marshal(p.ToMap())
}

func (p *CursorPaginator[T]) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(p.ToMap(), "", "    ")
}
