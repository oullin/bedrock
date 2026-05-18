package multisteps_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/workflow/multisteps"
)

func TestArgKindNames(t *testing.T) {
	cases := []struct {
		arg  multisteps.Arg
		want string
	}{
		{multisteps.Literal(42), "literal"},
		{multisteps.Variable("name"), "variable"},
		{multisteps.Response("job", "field"), "response"},
	}

	for _, c := range cases {
		if got := c.arg.Kind(); got != c.want {
			t.Errorf("Arg.Kind = %q, want %q", got, c.want)
		}
	}
}

func TestArgDependsOnJob(t *testing.T) {
	if got := multisteps.Response("job", "f").DependsOnJob(); got != "job" {
		t.Errorf("Response DependsOnJob = %q, want job", got)
	}

	if got := multisteps.Literal(1).DependsOnJob(); got != "" {
		t.Errorf("Literal DependsOnJob = %q, want empty", got)
	}

	if got := multisteps.Variable("v").DependsOnJob(); got != "" {
		t.Errorf("Variable DependsOnJob = %q, want empty", got)
	}
}

func TestAsTypedExtraction(t *testing.T) {
	r := multisteps.Result{
		Responses: map[string]any{
			"create": map[string]any{
				"id":   int64(7),
				"name": "ok",
			},
		},
	}

	id, err := multisteps.As[int64](r, "create", "id")

	if err != nil {
		t.Fatalf("As: %v", err)
	}

	if id != 7 {
		t.Errorf("id = %d, want 7", id)
	}

	name, err := multisteps.As[string](r, "create", "name")

	if err != nil || name != "ok" {
		t.Errorf("As string: name=%q err=%v", name, err)
	}
}

func TestAsReturnsErrorOnMissingJob(t *testing.T) {
	r := multisteps.Result{Responses: map[string]any{}}
	_, err := multisteps.As[int](r, "missing", "id")

	if err == nil {
		t.Fatal("expected error for missing job response")
	}
}

func TestAsReturnsErrorOnMissingField(t *testing.T) {
	r := multisteps.Result{Responses: map[string]any{"job": map[string]any{"a": 1}}}
	_, err := multisteps.As[int](r, "job", "missing")

	if err == nil {
		t.Fatal("expected error for missing field")
	}

	var ure *multisteps.UnresolvedResponseError

	if !errors.As(err, &ure) {
		t.Errorf("err = %v, want *UnresolvedResponseError", err)
	}
}

type fixture struct {
	ID   int
	Name string
}

func TestAsResolvesStructFieldsViaReflection(t *testing.T) {
	r := multisteps.Result{
		Responses: map[string]any{
			"create": fixture{ID: 42, Name: "answer"},
		},
	}

	id, err := multisteps.As[int](r, "create", "ID")

	if err != nil {
		t.Fatalf("As struct field: %v", err)
	}

	if id != 42 {
		t.Errorf("id = %d, want 42", id)
	}
}

func TestAsResolvesPointerStructFields(t *testing.T) {
	r := multisteps.Result{
		Responses: map[string]any{
			"create": &fixture{ID: 1, Name: "ptr"},
		},
	}

	name, err := multisteps.As[string](r, "create", "Name")

	if err != nil {
		t.Fatalf("As pointer struct field: %v", err)
	}

	if name != "ptr" {
		t.Errorf("name = %q, want ptr", name)
	}
}

func TestAsReturnsErrorOnTypeMismatch(t *testing.T) {
	r := multisteps.Result{Responses: map[string]any{"job": map[string]any{"id": "not-an-int"}}}
	_, err := multisteps.As[int](r, "job", "id")

	if err == nil {
		t.Fatal("expected error for type mismatch")
	}
}

func TestAsResolvesWholeResponseWhenFieldEmpty(t *testing.T) {
	r := multisteps.Result{Responses: map[string]any{"job": 42}}
	val, err := multisteps.As[int](r, "job", "")

	if err != nil {
		t.Fatalf("As empty field: %v", err)
	}

	if val != 42 {
		t.Errorf("val = %d, want 42", val)
	}
}
