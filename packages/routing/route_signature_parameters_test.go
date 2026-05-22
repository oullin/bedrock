package routing

import (
	"reflect"
	"testing"
)

// Ref: @bedrock/code-0395
// RouteSignatureParametersTest::test_it_can_extract_the_route_action_signature_parameters

type testEnum string

type testMarker interface {
	mark()
}

func TestRouteSignatureParameters(t *testing.T) {
	t.Run("test_parameters_can_be_retrieved_from_func", func(t *testing.T) {
		fn := func(a string, b int) {}
		a, err := ParseAction("/", fn)

		if err != nil {
			t.Fatal(err)
		}

		params := RouteSignatureParameters{}.FromAction(a, nil)

		if len(params) != 2 {
			t.Fatalf("got %d params, want 2", len(params))
		}

		if params[0].Type.Kind() != reflect.String || params[1].Type.Kind() != reflect.Int {
			t.Errorf("types = %v, %v", params[0].Type, params[1].Type)
		}
	})

	t.Run("test_filtering_by_backed_enum", func(t *testing.T) {
		fn := func(s string, e testEnum) {}
		a, _ := ParseAction("/", fn)
		params := RouteSignatureParameters{}.FromAction(a, map[string]any{"backedEnum": true})

		if len(params) != 1 {
			t.Fatalf("got %d params, want 1", len(params))
		}
	})

	t.Run("test_filtering_by_subclass_interface", func(t *testing.T) {
		fn := func(s string, m testMarker) {}
		a, _ := ParseAction("/", fn)
		params := RouteSignatureParameters{}.FromAction(a, map[string]any{
			"subClass": reflect.TypeOf((*testMarker)(nil)).Elem(),
		})

		if len(params) != 1 {
			t.Fatalf("got %d params, want 1", len(params))
		}
	})

	t.Run("test_non_func_action_returns_nil", func(t *testing.T) {
		a := &Action{Uses: "Controller@show"}
		params := RouteSignatureParameters{}.FromAction(a, nil)

		if params != nil {
			t.Errorf("expected nil, got %v", params)
		}
	})
}

func (e testEnum) BackingValue() string { return string(e) }
