package pagination

// Paginator defines the contract for offset-based paginators.
type Paginator interface {
	Url(page int) string
	PreviousPageUrl() string
	NextPageUrl() string
	Appends(query map[string]string) Paginator
	WithQueryString(query map[string]string) Paginator
	Fragment(fragment ...string) string
	Items() []any
	FirstItem() *int
	LastItem() *int
	PerPage() int
	CurrentPage() int
	HasPages() bool
	HasMorePages() bool
	OnFirstPage() bool
	OnLastPage() bool
	Path() string
	IsEmpty() bool
	IsNotEmpty() bool
	GetPageName() string
	SetPageName(name string)
	GetOptions() map[string]any
	ToMap() map[string]any
	ToJSON() ([]byte, error)
	ToPrettyJSON() ([]byte, error)
	Count() int
}
