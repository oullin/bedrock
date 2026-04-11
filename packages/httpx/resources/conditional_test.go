package resources_test

import (
	"testing"

	"github.com/bedrock/packages/httpx/resources"
)

func TestWhenTrue(t *testing.T) {
	t.Parallel()

	v := resources.When(true, "visible")

	if v.IsMissing() {
		t.Fatal("expected When(true) to not be missing")
	}

	if v.Value != "visible" {
		t.Fatalf("expected visible, got %v", v.Value)
	}
}

func TestWhenFalse(t *testing.T) {
	t.Parallel()

	v := resources.When(false, "hidden")

	if !v.IsMissing() {
		t.Fatal("expected When(false) to be missing")
	}
}

func TestUnlessTrue(t *testing.T) {
	t.Parallel()

	v := resources.Unless(true, "hidden")

	if !v.IsMissing() {
		t.Fatal("expected Unless(true) to be missing")
	}
}

func TestUnlessFalse(t *testing.T) {
	t.Parallel()

	v := resources.Unless(false, "visible")

	if v.IsMissing() {
		t.Fatal("expected Unless(false) to not be missing")
	}
}

func TestMergeWhenTrue(t *testing.T) {
	t.Parallel()

	v := resources.MergeWhen(true, map[string]any{"key": "val"})

	mv, ok := v.(resources.MergeValue)

	if !ok {
		t.Fatal("expected MergeValue")
	}

	if mv.Data["key"] != "val" {
		t.Fatal("expected merged data")
	}
}

func TestMergeWhenFalse(t *testing.T) {
	t.Parallel()

	v := resources.MergeWhen(false, map[string]any{"key": "val"})

	_, ok := v.(resources.MissingValue)

	if !ok {
		t.Fatal("expected MissingValue when condition is false")
	}
}

func TestMissingValue(t *testing.T) {
	t.Parallel()

	mv := resources.MissingValue{}

	if !mv.IsMissing() {
		t.Fatal("MissingValue should always be missing")
	}
}

func TestMergeValueIsMissing(t *testing.T) {
	t.Parallel()

	mv := resources.MergeValue{Data: map[string]any{"a": 1}}

	if mv.IsMissing() {
		t.Fatal("MergeValue should not be missing")
	}
}
