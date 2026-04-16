package query

import (
	"context"

	cpagination "github.com/bedrock/packages/contracts/pagination"
	"github.com/bedrock/packages/pagination"
)

// Paginate paginates the query results with total count.
func (b *Builder) Paginate(ctx context.Context, perPage, page int, columns ...string) (*pagination.LengthAwarePaginator[map[string]any], error) {
	if page < 1 {
		page = 1
	}

	total, err := b.CloneWithout("columns", "orders", "limit", "offset").Count(ctx)
	if err != nil {
		return nil, err
	}

	results, err := b.ForPage(page, perPage).Get(ctx, columns...)
	if err != nil {
		return nil, err
	}

	return pagination.NewLengthAwarePaginator(results, int(total), perPage, page), nil
}

// SimplePaginate paginates without total count.
func (b *Builder) SimplePaginate(ctx context.Context, perPage, page int, columns ...string) (*pagination.Paginator[map[string]any], error) {
	if page < 1 {
		page = 1
	}

	results, err := b.ForPage(page, perPage+1).Get(ctx, columns...)
	if err != nil {
		return nil, err
	}

	return pagination.NewPaginator(results, perPage, page), nil
}

// CursorPaginate paginates using cursor-based pagination.
func (b *Builder) CursorPaginate(ctx context.Context, perPage int, cursor *cpagination.Cursor, columns ...string) (*pagination.CursorPaginator[map[string]any], error) {
	if perPage < 1 {
		perPage = 15
	}

	results, err := b.Limit(perPage + 1).Get(ctx, columns...)
	if err != nil {
		return nil, err
	}

	return pagination.NewCursorPaginator(results, perPage, cursor), nil
}
