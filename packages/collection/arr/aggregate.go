package arr

// Every report whether every element passes the callback test.
func Every[T any](items []T, callback func(T, int) bool) bool {
	for i, item := range items {
		if !callback(item, i) {
			return false
		}
	}

	return true
}

// Some reports whether any element passes the callback test.
func Some[T any](items []T, callback func(T, int) bool) bool {
	for i, item := range items {
		if callback(item, i) {
			return true
		}
	}

	return false
}
