package pagination

import "encoding/json"

// Paginator is a simple paginator that does not know the total number of items.
// It only knows whether there are more pages after the current one.
type Paginator[T any] struct {
	abstractPaginator[T]
	hasMore bool
}

// NewPaginator creates a new simple paginator. If items contains more than
// perPage elements, the extra items are trimmed and hasMore is set to true.
// Options may include: "path" (string), "pageName" (string), "query"
// (map[string]string), "fragment" (string).
func NewPaginator[T any](items []T, perPage int, currentPage int, options ...map[string]any) *Paginator[T] {
	if currentPage < 1 {
		currentPage = 1
	}

	opts := mergeOptions(options)

	p := &Paginator[T]{
		abstractPaginator: abstractPaginator[T]{
			perPage:     perPage,
			currentPage: currentPage,
			pageName:    optString(opts, "pageName", "page"),
			query:       optQuery(opts),
			fragment:    optString(opts, "fragment", ""),
			onEachSide:  3,
			options:     opts,
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

func (p *Paginator[T]) HasMorePagesWhen(hasMore bool) *Paginator[T] {
	p.hasMore = hasMore

	return p
}

func (p *Paginator[T]) Url(page int) string {
	return p.buildUrl(page)
}

func (p *Paginator[T]) PreviousPageUrl() string {
	return p.previousPageUrl()
}

func (p *Paginator[T]) NextPageUrl() string {
	if p.hasMore {
		return p.buildUrl(p.currentPage + 1)
	}

	return ""
}

func (p *Paginator[T]) HasMorePages() bool {
	return p.hasMore
}

func (p *Paginator[T]) HasPages() bool {
	return p.currentPage != 1 || p.hasMore
}

func (p *Paginator[T]) OnFirstPage() bool {
	return p.onFirstPage()
}

func (p *Paginator[T]) OnLastPage() bool {
	return !p.hasMore
}

func (p *Paginator[T]) Appends(query map[string]string) *Paginator[T] {
	p.appendsQuery(query)

	return p
}

func (p *Paginator[T]) WithQueryString(query map[string]string) *Paginator[T] {
	p.appendsQuery(query)

	return p
}

func (p *Paginator[T]) Fragment(fragment ...string) string {
	if len(fragment) > 0 {
		p.setFragment(fragment[0])
	}

	return p.getFragment()
}

func (p *Paginator[T]) Items() []any {
	return p.itemsToAny()
}

func (p *Paginator[T]) TypedItems() []T {
	return p.getItems()
}

func (p *Paginator[T]) FirstItem() *int {
	return p.firstItem()
}

func (p *Paginator[T]) LastItem() *int {
	return p.lastItem()
}

func (p *Paginator[T]) Through(callback func(T) T) *Paginator[T] {
	p.transformItems(callback)

	return p
}

func (p *Paginator[T]) PerPage() int {
	return p.getPerPage()
}

func (p *Paginator[T]) CurrentPage() int {
	return p.getCurrentPage()
}

func (p *Paginator[T]) Path() string {
	return p.getPath()
}

func (p *Paginator[T]) IsEmpty() bool {
	return p.isEmpty()
}

func (p *Paginator[T]) IsNotEmpty() bool {
	return p.isNotEmpty()
}

func (p *Paginator[T]) Count() int {
	return p.count()
}

func (p *Paginator[T]) GetPageName() string {
	return p.getPageName()
}

func (p *Paginator[T]) SetPageName(name string) {
	p.setPageName(name)
}

func (p *Paginator[T]) WithPath(path string) *Paginator[T] {
	p.setPath(path)

	return p
}

func (p *Paginator[T]) SetPath(path string) *Paginator[T] {
	p.setPath(path)

	return p
}

func (p *Paginator[T]) GetOptions() map[string]any {
	return p.getOptions()
}

func (p *Paginator[T]) ToMap() map[string]any {
	result := map[string]any{
		"current_page":   p.currentPage,
		"data":           p.items,
		"first_page_url": p.buildUrl(1),
		"from":           p.firstItem(),
		"next_page_url":  nilIfEmpty(p.NextPageUrl()),
		"path":           p.path,
		"per_page":       p.perPage,
		"prev_page_url":  nilIfEmpty(p.PreviousPageUrl()),
		"to":             p.lastItem(),
	}

	return result
}

func (p *Paginator[T]) ToJSON() ([]byte, error) {
	return json.Marshal(p.ToMap())
}

func (p *Paginator[T]) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(p.ToMap(), "", "    ")
}

func mergeOptions(options []map[string]any) map[string]any {
	opts := make(map[string]any)

	if len(options) > 0 && options[0] != nil {
		for k, v := range options[0] {
			opts[k] = v
		}
	}

	return opts
}

func optString(opts map[string]any, key, fallback string) string {
	if v, ok := opts[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return fallback
}

func optQuery(opts map[string]any) map[string]string {
	if v, ok := opts["query"]; ok {
		if q, ok := v.(map[string]string); ok {
			result := make(map[string]string, len(q))

			for k, val := range q {
				result[k] = val
			}

			return result
		}
	}

	return make(map[string]string)
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}

	return s
}
