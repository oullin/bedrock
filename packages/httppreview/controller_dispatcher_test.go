package httppreview_test

import (
	"strings"
	"testing"

	"github.com/bedrock/packages/httppreview"
)

// testController is a fake controller for testing the controller dispatcher.
type testController struct{}

// testController does not implement HasMiddleware, so should return empty.

// middlewareController implements HasMiddleware for testing.
type middlewareController struct{}

type testMiddlewareDef struct {
	Middleware any
	Only       []string
	Except     []string
}

func (c *testController) Store(name string) string {
	return "stored: " + name
}

func TestControllerDispatcherResolvesAndAborts(t *testing.T) {
	t.Parallel()

	d := httppreview.NewControllerDispatcher(nil)
	ctrl := &testController{}

	defer func() {
		v := recover()

		if v == nil {
			t.Fatal("expected SuccessResponse panic from Dispatch")
		}

		if _, ok := v.(httppreview.SuccessResponse); !ok {
			t.Fatalf("expected SuccessResponse, got %T: %v", v, v)
		}
	}()

	d.Dispatch(nil, ctrl, "Store")
}

func TestControllerDispatcherWithRouteAccessor(t *testing.T) {
	t.Parallel()

	d := httppreview.NewControllerDispatcher(nil)
	ctrl := &testController{}

	route := &fakeRoute{
		params:     map[string]string{"name": "world"},
		paramNames: []string{"name"},
	}

	defer func() {
		v := recover()

		if v == nil {
			t.Fatal("expected SuccessResponse panic from Dispatch")
		}

		if _, ok := v.(httppreview.SuccessResponse); !ok {
			t.Fatalf("expected SuccessResponse, got %T: %v", v, v)
		}
	}()

	d.Dispatch(route, ctrl, "Store")
}

func TestControllerDispatcherEnsureMethodExists(t *testing.T) {
	t.Parallel()

	d := httppreview.NewControllerDispatcher(nil)
	ctrl := &testController{}

	defer func() {
		v := recover()

		if v == nil {
			t.Fatal("expected panic for missing method")
		}

		msg, ok := v.(string)

		if !ok {
			t.Fatalf("expected string panic, got %T: %v", v, v)
		}

		if !strings.Contains(msg, "NoSuchMethod") {
			t.Fatalf("expected panic message to contain method name, got %q", msg)
		}

		if !strings.Contains(msg, "not defined") {
			t.Fatalf("expected panic message to contain 'not defined', got %q", msg)
		}
	}()

	d.Dispatch(nil, ctrl, "NoSuchMethod")
}

func TestControllerDispatcherGetMiddleware(t *testing.T) {
	t.Parallel()

	d := httppreview.NewControllerDispatcher(nil)

	mw := d.GetMiddleware(&testController{}, "Store")

	if len(mw) != 0 {
		t.Fatalf("expected empty middleware list, got %d items", len(mw))
	}
}

func (c *middlewareController) Store()   {}
func (c *middlewareController) Update()  {}
func (c *middlewareController) Destroy() {}
