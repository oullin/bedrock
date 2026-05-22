package routing

// ViewFactory is the minimum surface ViewController needs from the bedrock
// view layer. The service provider in M11 wires in a real implementation.
type ViewFactory interface {
	Make(view string, data map[string]any) any
}

// ViewController is the invokable controller used by [Router.View]. It
// renders the named view with the supplied data, merging in any extra route
// parameters as additional data.
//
// Ref: @bedrock/code-0347
type ViewController struct {
	Controller
	view ViewFactory
}

// NewViewController wires the controller to a view factory.
func NewViewController(view ViewFactory) *ViewController {
	return &ViewController{view: view}
}

// Invoke dispatches the view with the route parameters merged into data.
func (c *ViewController) Invoke(route *Route, viewName string, data map[string]any) any {
	if data == nil {
		data = map[string]any{}
	}

	for k, v := range route.Parameters {
		if k == "view" || k == "data" || k == "status" || k == "headers" {
			continue
		}

		data[k] = v
	}

	if c.view == nil {
		return nil
	}

	return c.view.Make(viewName, data)
}
