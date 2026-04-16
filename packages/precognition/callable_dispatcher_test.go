package precognition_test

import (
	"testing"

	"github.com/bedrock/packages/precognition"
)

func TestCallableDispatcherResolvesAndAborts(t *testing.T) {
	t.Parallel()

	d := precognition.NewCallableDispatcher(nil)
	called := false

	callable := func(name string) string {
		called = true

		return "hello " + name
	}

	defer func() {
		v := recover()

		if v == nil {
			t.Fatal("expected SuccessResponse panic from Dispatch")
		}

		if _, ok := v.(precognition.SuccessResponse); !ok {
			t.Fatalf("expected SuccessResponse, got %T: %v", v, v)
		}

		if called {
			t.Fatal("callable should not have been executed")
		}
	}()

	d.Dispatch(nil, callable)
}

func TestCallableDispatcherWithRouteAccessor(t *testing.T) {
	t.Parallel()

	d := precognition.NewCallableDispatcher(nil)

	callable := func(name string) string {
		return "hello " + name
	}

	route := &fakeRoute{
		params:     map[string]string{"name": "world"},
		paramNames: []string{"name"},
	}

	defer func() {
		v := recover()

		if v == nil {
			t.Fatal("expected SuccessResponse panic from Dispatch")
		}

		if _, ok := v.(precognition.SuccessResponse); !ok {
			t.Fatalf("expected SuccessResponse, got %T: %v", v, v)
		}
	}()

	d.Dispatch(route, callable)
}

// fakeRoute satisfies the routeAccessor interface used by the dispatchers.
type fakeRoute struct {
	params     map[string]string
	paramNames []string
}

func (r *fakeRoute) ParametersWithoutNulls() map[string]string { return r.params }
func (r *fakeRoute) ParameterNames() []string                  { return r.paramNames }
