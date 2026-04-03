package arr

// Flatten flattens a slice of slices into a single slice.
func Flatten[T any](items [][]T) []T {
	result := make([]T, 0)

	for _, inner := range items {
		result = append(result, inner...)
	}

	return result
}

// Collapse merges a slice of slices into a single slice.
// It is an alias for [Flatten].
func Collapse[T any](items [][]T) []T {
	return Flatten(items)
}

// CrossJoin returns the Cartesian product of the given slices.
func CrossJoin[T any](lists ...[]T) [][]T {
	results := [][]T{{}}

	for _, list := range lists {
		var newResults [][]T

		for _, result := range results {
			for _, item := range list {
				newResult := make([]T, len(result)+1)
				copy(newResult, result)
				newResult[len(result)] = item
				newResults = append(newResults, newResult)
			}
		}

		results = newResults
	}

	return results
}
