package pagination

// LengthAwarePaginator defines the contract for paginators that know the total
// number of items and can compute last page, URL ranges, etc.
type LengthAwarePaginator interface {
	Paginator
	Total() int
	LastPage() int
	GetUrlRange(start, end int) map[int]string
}
