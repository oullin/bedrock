package routing

import (
	"errors"
	"testing"
)

// Ref: @bedrock/code-0392
// RouteBindingTest::test_it_can_resolve_the_explicit_model_for_the_given_route
// RouteBindingTest::test_it_cannot_resolve_the_explicit_soft_deleted_model_for_the_given_route
// RouteBindingTest::test_it_can_resolve_the_explicit_soft_deleted_model_for_the_given_route_with_trashed

type explicitBindingModel struct {
	id      string
	trashed bool
}

type explicitBindingContainer struct {
	model *explicitBindingModel
}

func (m *explicitBindingModel) ResolveRouteBinding(value, field string) (any, error) {
	if m.trashed {
		return nil, nil
	}

	if value == "1" {
		return &explicitBindingModel{id: "1"}, nil
	}

	return nil, nil
}

func (m *explicitBindingModel) ResolveSoftDeletableRouteBinding(value, field string) (any, error) {
	if value == "1" {
		return &explicitBindingModel{id: "1", trashed: true}, nil
	}

	return nil, nil
}

func (m *explicitBindingModel) IsSoftDeletable() bool { return true }

func (c explicitBindingContainer) Make(abstract string) (any, error) {
	return c.model, nil
}

func TestRouteBinding_ForModel(t *testing.T) {
	// RouteBindingTest::test_it_can_resolve_the_explicit_model_for_the_given_route
	t.Run("test_resolves_the_explicit_model", func(t *testing.T) {
		route := NewRoute("GET", "/users/{user}", func() {})
		resolver := ForModel(explicitBindingContainer{model: &explicitBindingModel{}}, "User", nil)

		got, err := resolver("1", route)

		if err != nil {
			t.Fatal(err)
		}

		model, ok := got.(*explicitBindingModel)

		if !ok {
			t.Fatalf("model type = %T", got)
		}

		if model.id != "1" {
			t.Errorf("id = %q", model.id)
		}
	})

	// RouteBindingTest::test_it_cannot_resolve_the_explicit_soft_deleted_model_for_the_given_route
	t.Run("test_rejects_soft_deleted_model_without_trashed", func(t *testing.T) {
		route := NewRoute("GET", "/users/{user}", func() {})
		resolver := ForModel(explicitBindingContainer{model: &explicitBindingModel{trashed: true}}, "User", nil)

		got, err := resolver("1", route)

		if got != nil {
			t.Fatalf("model = %#v, want nil", got)
		}

		var notFound *ModelNotFoundError

		if !errors.As(err, &notFound) {
			t.Fatalf("err = %v, want ModelNotFoundError", err)
		}
	})

	// RouteBindingTest::test_it_can_resolve_the_explicit_soft_deleted_model_for_the_given_route_with_trashed
	t.Run("test_resolves_soft_deleted_model_with_trashed", func(t *testing.T) {
		route := NewRoute("GET", "/users/{user}", func() {}).WithTrashed()
		resolver := ForModel(explicitBindingContainer{model: &explicitBindingModel{trashed: true}}, "User", nil)

		got, err := resolver("1", route)

		if err != nil {
			t.Fatal(err)
		}

		model, ok := got.(*explicitBindingModel)

		if !ok {
			t.Fatalf("model type = %T", got)
		}

		if !model.trashed {
			t.Error("expected soft-deleted model")
		}
	})
}
