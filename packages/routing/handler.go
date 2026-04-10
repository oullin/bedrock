package routing

// HandlerFunc handles a routed HTTP request. It returns an error which, if
// non-nil, results in a 500 Internal Server Error response.
type HandlerFunc func(*Context) error
