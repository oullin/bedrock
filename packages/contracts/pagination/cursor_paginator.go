package pagination

// CursorPaginator defines the contract for cursor-based paginators.
type CursorPaginator interface {
	Url(cursor *Cursor) string
	PreviousPageUrl() string
	NextPageUrl() string
	Appends(query map[string]string) CursorPaginator
	WithQueryString(query map[string]string) CursorPaginator
	Fragment(fragment ...string) string
	Items() []any
	PerPage() int
	Cursor() *Cursor
	NextCursor() *Cursor
	PreviousCursor() *Cursor
	HasPages() bool
	HasMorePages() bool
	OnFirstPage() bool
	OnLastPage() bool
	Path() string
	IsEmpty() bool
	IsNotEmpty() bool
	GetCursorName() string
	SetCursorName(name string)
	GetOptions() map[string]any
	ToMap() map[string]any
	ToJSON() ([]byte, error)
	ToPrettyJSON() ([]byte, error)
	Count() int
}
