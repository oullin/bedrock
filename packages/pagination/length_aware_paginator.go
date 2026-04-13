package pagination

import (
	"encoding/json"
	"math"
)

// LengthAwarePaginator is a paginator that knows the total number of items and
// can compute last page, URL ranges, and page link windows.
type LengthAwarePaginator[T any] struct {
	abstractPaginator[T]
	total    int
	lastPage int
}

// NewLengthAwarePaginator creates a new length-aware paginator. Options may
// include: "path" (string), "pageName" (string), "query" (map[string]string),
// "fragment" (string).
func NewLengthAwarePaginator[T any](items []T, total int, perPage int, currentPage int, options ...map[string]any) *LengthAwarePaginator[T] {
	opts := mergeOptions(options)

	lastPage := int(math.Max(math.Ceil(float64(total)/float64(perPage)), 1))

	if currentPage < 1 {
		currentPage = 1
	}

	p := &LengthAwarePaginator[T]{
		abstractPaginator: abstractPaginator[T]{
			perPage:     perPage,
			currentPage: currentPage,
			pageName:    optString(opts, "pageName", "page"),
			query:       optQuery(opts),
			fragment:    optString(opts, "fragment", ""),
			onEachSide:  3,
			options:     opts,
		},
		total:    total,
		lastPage: lastPage,
	}

	p.setPath(optString(opts, "path", "/"))

	if len(items) > perPage {
		items = items[:perPage]
	}

	p.items = items

	return p
}

func (p *LengthAwarePaginator[T]) Url(page int) string {
	return p.buildUrl(page)
}

func (p *LengthAwarePaginator[T]) PreviousPageUrl() string {
	return p.previousPageUrl()
}

func (p *LengthAwarePaginator[T]) NextPageUrl() string {
	if p.currentPage < p.lastPage {
		return p.buildUrl(p.currentPage + 1)
	}

	return ""
}

func (p *LengthAwarePaginator[T]) HasMorePages() bool {
	return p.currentPage < p.lastPage
}

func (p *LengthAwarePaginator[T]) HasPages() bool {
	return p.lastPage > 1
}

func (p *LengthAwarePaginator[T]) OnFirstPage() bool {
	return p.onFirstPage()
}

func (p *LengthAwarePaginator[T]) OnLastPage() bool {
	return p.currentPage >= p.lastPage
}

func (p *LengthAwarePaginator[T]) Total() int {
	return p.total
}

func (p *LengthAwarePaginator[T]) LastPage() int {
	return p.lastPage
}

func (p *LengthAwarePaginator[T]) GetUrlRange(start, end int) map[int]string {
	return p.getUrlRange(start, end)
}

func (p *LengthAwarePaginator[T]) OnEachSide(count int) *LengthAwarePaginator[T] {
	p.setOnEachSide(count)

	return p
}

func (p *LengthAwarePaginator[T]) Appends(query map[string]string) *LengthAwarePaginator[T] {
	p.appendsQuery(query)

	return p
}

func (p *LengthAwarePaginator[T]) WithQueryString(query map[string]string) *LengthAwarePaginator[T] {
	p.appendsQuery(query)

	return p
}

func (p *LengthAwarePaginator[T]) Fragment(fragment ...string) string {
	if len(fragment) > 0 {
		p.setFragment(fragment[0])
	}

	return p.getFragment()
}

func (p *LengthAwarePaginator[T]) Items() []any {
	return p.itemsToAny()
}

func (p *LengthAwarePaginator[T]) TypedItems() []T {
	return p.getItems()
}

func (p *LengthAwarePaginator[T]) FirstItem() *int {
	return p.firstItem()
}

func (p *LengthAwarePaginator[T]) LastItem() *int {
	return p.lastItem()
}

func (p *LengthAwarePaginator[T]) Through(callback func(T) T) *LengthAwarePaginator[T] {
	p.transformItems(callback)

	return p
}

func (p *LengthAwarePaginator[T]) PerPage() int {
	return p.getPerPage()
}

func (p *LengthAwarePaginator[T]) CurrentPage() int {
	return p.getCurrentPage()
}

func (p *LengthAwarePaginator[T]) Path() string {
	return p.getPath()
}

func (p *LengthAwarePaginator[T]) IsEmpty() bool {
	return p.isEmpty()
}

func (p *LengthAwarePaginator[T]) IsNotEmpty() bool {
	return p.isNotEmpty()
}

func (p *LengthAwarePaginator[T]) Count() int {
	return p.count()
}

func (p *LengthAwarePaginator[T]) GetPageName() string {
	return p.getPageName()
}

func (p *LengthAwarePaginator[T]) SetPageName(name string) {
	p.setPageName(name)
}

func (p *LengthAwarePaginator[T]) WithPath(path string) *LengthAwarePaginator[T] {
	p.setPath(path)

	return p
}

func (p *LengthAwarePaginator[T]) SetPath(path string) *LengthAwarePaginator[T] {
	p.setPath(path)

	return p
}

func (p *LengthAwarePaginator[T]) GetOptions() map[string]any {
	return p.getOptions()
}

func (p *LengthAwarePaginator[T]) ToMap() map[string]any {
	return map[string]any{
		"current_page":   p.currentPage,
		"data":           p.items,
		"first_page_url": p.buildUrl(1),
		"from":           p.firstItem(),
		"last_page":      p.lastPage,
		"last_page_url":  p.buildUrl(p.lastPage),
		"links":          p.linkCollection(),
		"next_page_url":  nilIfEmpty(p.NextPageUrl()),
		"path":           p.path,
		"per_page":       p.perPage,
		"prev_page_url":  nilIfEmpty(p.PreviousPageUrl()),
		"to":             p.lastItem(),
		"total":          p.total,
	}
}

func (p *LengthAwarePaginator[T]) linkCollection() []map[string]any {
	window := MakeUrlWindow(p)

	var links []map[string]any

	// Previous link.
	links = append(links, map[string]any{
		"url":    nilIfEmpty(p.PreviousPageUrl()),
		"label":  "&laquo; Previous",
		"active": false,
	})

	// Page number links from the window.
	for _, section := range []string{"first", "slider", "last"} {
		pages := window[section]

		if pages == nil {
			continue
		}

		pageNums := sortedKeys(pages)

		for _, page := range pageNums {
			links = append(links, map[string]any{
				"url":    pages[page],
				"label":  intToString(page),
				"active": page == p.currentPage,
			})
		}
	}

	// Next link.
	links = append(links, map[string]any{
		"url":    nilIfEmpty(p.NextPageUrl()),
		"label":  "Next &raquo;",
		"active": false,
	})

	return links
}

func (p *LengthAwarePaginator[T]) ToJSON() ([]byte, error) {
	return json.Marshal(p.ToMap())
}

func (p *LengthAwarePaginator[T]) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(p.ToMap(), "", "    ")
}

func sortedKeys(m map[int]string) []int {
	keys := make([]int, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}

	return keys
}
