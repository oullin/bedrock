package routing

// RedirectController is the invokable controller used by [Router.Redirect]
// and [Router.PermanentRedirect]. It reads the destination and status from
// the matched route's parameters and produces a [RedirectResponse].
//
// Ref: @bedrock/code-0328
type RedirectController struct{ Controller }

// Invoke is the dispatch entry point — controller dispatchers detect the
// "Invoke" method by reflection and call it like an invokable.
func (c *RedirectController) Invoke(route *Route) *RedirectResponse {
	destination := route.Parameter("destination", "")
	status := 302

	if s := route.Parameter("status", ""); s != "" {
		// best-effort parse; default keeps 302
		var n int

		for _, ch := range s {
			if ch < '0' || ch > '9' {
				n = 0

				break
			}

			n = n*10 + int(ch-'0')
		}

		if n != 0 {
			status = n
		}
	}

	return &RedirectResponse{URL: destination, Status: status}
}
