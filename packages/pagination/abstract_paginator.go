package pagination

import (
	"net/url"
	"sort"
	"strings"
)

// abstractPaginator provides shared functionality for offset-based paginators.
type abstractPaginator[T any] struct {
	items       []T
	perPage     int
	currentPage int
	path        string
	query       map[string]string
	fragment    string
	pageName    string
	onEachSide  int
	options     map[string]any
}

func (p *abstractPaginator[T]) isValidPageNumber(page int) bool {
	return page >= 1
}

func (p *abstractPaginator[T]) buildUrl(page int) string {
	if page <= 0 {
		page = 1
	}

	path := p.path

	params := url.Values{}
	// Sort keys for deterministic output.
	keys := make([]string, 0, len(p.query))

	for k := range p.query {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		params.Set(k, p.query[k])
	}

	params.Set(p.pageName, intToString(page))

	result := path + "?" + encodeQuery(params)

	if p.fragment != "" {
		result += "#" + p.fragment
	}

	return result
}

func (p *abstractPaginator[T]) previousPageUrl() string {
	if p.currentPage > 1 {
		return p.buildUrl(p.currentPage - 1)
	}

	return ""
}

func (p *abstractPaginator[T]) appendsQuery(query map[string]string) {
	for k, v := range query {
		p.query[k] = v
	}
}

func (p *abstractPaginator[T]) getFragment() string {
	return p.fragment
}

func (p *abstractPaginator[T]) setFragment(f string) {
	p.fragment = f
}

func (p *abstractPaginator[T]) getItems() []T {
	return p.items
}

func (p *abstractPaginator[T]) firstItem() *int {
	if len(p.items) == 0 {
		return nil
	}

	v := (p.currentPage-1)*p.perPage + 1

	return &v
}

func (p *abstractPaginator[T]) lastItem() *int {
	if len(p.items) == 0 {
		return nil
	}

	first := (p.currentPage-1)*p.perPage + 1
	v := first + len(p.items) - 1

	return &v
}

func (p *abstractPaginator[T]) transformItems(callback func(T) T) {
	for i, item := range p.items {
		p.items[i] = callback(item)
	}
}

func (p *abstractPaginator[T]) count() int {
	return len(p.items)
}

func (p *abstractPaginator[T]) isEmpty() bool {
	return len(p.items) == 0
}

func (p *abstractPaginator[T]) isNotEmpty() bool {
	return len(p.items) > 0
}

func (p *abstractPaginator[T]) onFirstPage() bool {
	return p.currentPage <= 1
}

func (p *abstractPaginator[T]) getPerPage() int {
	return p.perPage
}

func (p *abstractPaginator[T]) getCurrentPage() int {
	return p.currentPage
}

func (p *abstractPaginator[T]) getPageName() string {
	return p.pageName
}

func (p *abstractPaginator[T]) setPageName(name string) {
	p.pageName = name
}

func (p *abstractPaginator[T]) getPath() string {
	return p.path
}

func (p *abstractPaginator[T]) setPath(path string) {
	p.path = strings.TrimRight(path, "/")
}

func (p *abstractPaginator[T]) getOptions() map[string]any {
	return p.options
}

func (p *abstractPaginator[T]) getOnEachSide() int {
	return p.onEachSide
}

func (p *abstractPaginator[T]) setOnEachSide(count int) {
	p.onEachSide = count
}

func (p *abstractPaginator[T]) getUrlRange(start, end int) map[int]string {
	urls := make(map[int]string)

	for i := start; i <= end; i++ {
		urls[i] = p.buildUrl(i)
	}

	return urls
}

func (p *abstractPaginator[T]) itemsToAny() []any {
	result := make([]any, len(p.items))

	for i, item := range p.items {
		result[i] = item
	}

	return result
}
