package engines

import (
	"context"

	contract "github.com/bedrock/packages/contracts/scout"
)

// NullEngine is a no-op search engine used for testing and development.
// All methods return zero values and nil errors.
type NullEngine struct{}

// Compile-time interface check.

// NewNullEngine creates a new NullEngine.

// NullResult is the empty result returned by NullEngine search operations.
type NullResult struct{}

var _ contract.Engine = (*NullEngine)(nil)

func NewNullEngine() *NullEngine {
	return &NullEngine{}
}

func (e *NullEngine) Update(_ context.Context, _ []contract.Searchable) error {
	return nil
}

func (e *NullEngine) Delete(_ context.Context, _ []contract.Searchable) error {
	return nil
}

func (e *NullEngine) Search(_ context.Context, _ contract.SearchBuilder) (any, error) {
	return &NullResult{}, nil
}

func (e *NullEngine) Paginate(_ context.Context, _ contract.SearchBuilder, _ int, _ int) (any, error) {
	return &NullResult{}, nil
}

func (e *NullEngine) MapIds(_ any) []any {
	return []any{}
}

func (e *NullEngine) Map(_ context.Context, _ any, _ contract.Searchable) ([]contract.Searchable, error) {
	return []contract.Searchable{}, nil
}

func (e *NullEngine) LazyMap(_ context.Context, _ any, _ contract.Searchable) func(yield func(contract.Searchable) bool) {
	return func(_ func(contract.Searchable) bool) {}
}

func (e *NullEngine) GetTotalCount(_ any) int64 {
	return 0
}

func (e *NullEngine) Flush(_ context.Context, _ contract.Searchable) error {
	return nil
}

func (e *NullEngine) CreateIndex(_ context.Context, _ string, _ map[string]any) error {
	return nil
}

func (e *NullEngine) DeleteIndex(_ context.Context, _ string) error {
	return nil
}
