package props_test

import (
	"testing"

	"github.com/bedrock/packages/inertia/props"
)

func TestAlwaysProp(t *testing.T) {
	t.Parallel()

	p := props.Always("hello")

	if got := p.Prop(); got != "hello" {
		t.Errorf("Always.Prop() = %v, want %q", got, "hello")
	}
}

func TestOptionalProp(t *testing.T) {
	t.Parallel()

	p := props.Optional(42)

	if got := p.Prop(); got != 42 {
		t.Errorf("Optional.Prop() = %v, want 42", got)
	}
}

func TestDeferProp(t *testing.T) {
	t.Parallel()

	t.Run("default group", func(t *testing.T) {
		t.Parallel()

		p := props.Defer("data")

		if p.Group != "default" {
			t.Errorf("Defer.Group = %q, want %q", p.Group, "default")
		}

		if p.IsMerge() {
			t.Error("Defer.IsMerge() = true, want false")
		}

		if got := p.Prop(); got != "data" {
			t.Errorf("Defer.Prop() = %v, want %q", got, "data")
		}
	})

	t.Run("custom group", func(t *testing.T) {
		t.Parallel()

		p := props.Defer("data", "sidebar")

		if p.Group != "sidebar" {
			t.Errorf("Defer.Group = %q, want %q", p.Group, "sidebar")
		}
	})

	t.Run("merge", func(t *testing.T) {
		t.Parallel()

		p := props.Defer("data").Merge()

		if !p.IsMerge() {
			t.Error("Defer.Merge().IsMerge() = false, want true")
		}
	})
}

func TestOnceProp(t *testing.T) {
	t.Parallel()

	p := props.Once("once-val")

	if got := p.Prop(); got != "once-val" {
		t.Errorf("Once.Prop() = %v, want %q", got, "once-val")
	}
}

func TestMergeProp(t *testing.T) {
	t.Parallel()

	t.Run("shallow", func(t *testing.T) {
		t.Parallel()

		p := props.Merge([]int{1, 2})

		if p.IsDeep() {
			t.Error("Merge.IsDeep() = true, want false")
		}
	})

	t.Run("deep", func(t *testing.T) {
		t.Parallel()

		p := props.DeepMerge(map[string]int{"a": 1})

		if !p.IsDeep() {
			t.Error("DeepMerge.IsDeep() = false, want true")
		}
	})
}

func TestScrollProp(t *testing.T) {
	t.Parallel()

	p := props.Scroll([]int{1, 2}, "feedPage", 1, nil, 2).Merge().Reset()

	if !p.IsMerge() {
		t.Error("Scroll.Merge().IsMerge() = false, want true")
	}

	if !p.IsReset() {
		t.Error("Scroll.Reset().IsReset() = false, want true")
	}

	if p.PageName != "feedPage" {
		t.Errorf("Scroll.PageName = %q, want %q", p.PageName, "feedPage")
	}
}

// --- AlwaysProp Extended ---

func TestAlwaysProp_WithFunc(t *testing.T) {
	t.Parallel()

	fn := func() any { return "lazy" }

	p := props.Always(fn)

	// Prop() should return the func itself, not invoke it.
	if p.Prop() == nil {
		t.Error("Always wrapping func should not return nil")
	}
}

func TestAlwaysProp_WithNil(t *testing.T) {
	t.Parallel()

	p := props.Always(nil)

	if p.Prop() != nil {
		t.Errorf("Always(nil).Prop() = %v, want nil", p.Prop())
	}
}

func TestAlwaysProp_WithStruct(t *testing.T) {
	t.Parallel()

	type info struct{ Name string }

	p := props.Always(info{Name: "test"})

	got, ok := p.Prop().(info)

	if !ok {
		t.Fatal("expected info struct")
	}

	if got.Name != "test" {
		t.Errorf("Name = %q", got.Name)
	}
}

// --- OptionalProp Extended ---

func TestOptionalProp_WithFunc(t *testing.T) {
	t.Parallel()

	fn := func() any { return 42 }

	p := props.Optional(fn)

	if p.Prop() == nil {
		t.Error("Optional wrapping func should not return nil")
	}
}

func TestOptionalProp_WithNil(t *testing.T) {
	t.Parallel()

	p := props.Optional(nil)

	if p.Prop() != nil {
		t.Errorf("Optional(nil).Prop() = %v, want nil", p.Prop())
	}
}

// --- DeferProp Extended ---

func TestDeferProp_EmptyGroupDefaultsToDefault(t *testing.T) {
	t.Parallel()

	p := props.Defer("val", "")

	if p.Group != "default" {
		t.Errorf("Defer with empty group = %q, want %q", p.Group, "default")
	}
}

func TestDeferProp_MergeReturnsCopy(t *testing.T) {
	t.Parallel()

	original := props.Defer("val")
	merged := original.Merge()

	if original.IsMerge() {
		t.Error("original should not be affected by Merge()")
	}

	if !merged.IsMerge() {
		t.Error("merged copy should have IsMerge() = true")
	}
}

func TestDeferProp_IsMergeDefaultFalse(t *testing.T) {
	t.Parallel()

	p := props.Defer("val")

	if p.IsMerge() {
		t.Error("fresh DeferProp should have IsMerge() = false")
	}
}

func TestDeferProp_WithFunc(t *testing.T) {
	t.Parallel()

	fn := func() any { return "lazy" }

	p := props.Defer(fn)

	if p.Prop() == nil {
		t.Error("Defer wrapping func should not return nil")
	}
}

func TestDeferProp_PropReturnsValue(t *testing.T) {
	t.Parallel()

	p := props.Defer("data")

	if p.Prop() != "data" {
		t.Errorf("Prop() = %v, want %q", p.Prop(), "data")
	}
}

// --- OnceProp Extended ---

func TestOnceProp_WithFunc(t *testing.T) {
	t.Parallel()

	fn := func() any { return "lazy" }

	p := props.Once(fn)

	if p.Prop() == nil {
		t.Error("Once wrapping func should not return nil")
	}
}

func TestOnceProp_WithNil(t *testing.T) {
	t.Parallel()

	p := props.Once(nil)

	if p.Prop() != nil {
		t.Errorf("Once(nil).Prop() = %v, want nil", p.Prop())
	}
}

func TestOnceProp_ExpiresAt(t *testing.T) {
	t.Parallel()

	original := props.Once("data")
	withExpiry := original.ExpiresAt(1700000000)

	if original.GetExpiresAt() != nil {
		t.Error("original should not be affected by ExpiresAt()")
	}

	if withExpiry.GetExpiresAt() == nil || *withExpiry.GetExpiresAt() != 1700000000 {
		t.Errorf("ExpiresAt = %v, want 1700000000", withExpiry.GetExpiresAt())
	}
}

func TestOnceProp_ExpiresAtPreservesValue(t *testing.T) {
	t.Parallel()

	p := props.Once("data").ExpiresAt(1700000000)

	if got := p.Prop(); got != "data" {
		t.Errorf("Prop() = %v, want %q", got, "data")
	}
}

func TestOnceProp_DefaultExpiresAtIsNil(t *testing.T) {
	t.Parallel()

	p := props.Once("data")

	if p.GetExpiresAt() != nil {
		t.Errorf("default GetExpiresAt() = %v, want nil", p.GetExpiresAt())
	}
}

// --- MergeProp Extended ---

func TestMergeProp_Prop(t *testing.T) {
	t.Parallel()

	p := props.Merge([]int{1, 2, 3})

	got, ok := p.Prop().([]int)

	if !ok {
		t.Fatal("expected []int")
	}

	if len(got) != 3 {
		t.Errorf("len = %d, want 3", len(got))
	}
}

func TestMergeProp_DeepMergeProp(t *testing.T) {
	t.Parallel()

	p := props.DeepMerge(map[string]int{"a": 1})

	if !p.IsDeep() {
		t.Error("DeepMerge should have IsDeep() = true")
	}

	got, ok := p.Prop().(map[string]int)

	if !ok {
		t.Fatal("expected map[string]int")
	}

	if got["a"] != 1 {
		t.Errorf("a = %v", got["a"])
	}
}

func TestMergeProp_ShallowIsNotDeep(t *testing.T) {
	t.Parallel()

	p := props.Merge([]int{1})

	if p.IsDeep() {
		t.Error("Merge should have IsDeep() = false")
	}
}

// --- ScrollProp Extended ---

func TestScrollProp_DefaultMergeIsFalse(t *testing.T) {
	t.Parallel()

	p := props.Scroll([]int{1}, "p", 1, nil, 2)

	if p.IsMerge() {
		t.Error("fresh ScrollProp should have IsMerge() = false")
	}
}

func TestScrollProp_DefaultResetIsFalse(t *testing.T) {
	t.Parallel()

	p := props.Scroll([]int{1}, "p", 1, nil, 2)

	if p.IsReset() {
		t.Error("fresh ScrollProp should have IsReset() = false")
	}
}

func TestScrollProp_MergeReturnsCopy(t *testing.T) {
	t.Parallel()

	original := props.Scroll([]int{1}, "p", 1, nil, 2)
	merged := original.Merge()

	if original.IsMerge() {
		t.Error("original should not be affected by Merge()")
	}

	if !merged.IsMerge() {
		t.Error("merged copy should have IsMerge() = true")
	}
}

func TestScrollProp_ResetReturnsCopy(t *testing.T) {
	t.Parallel()

	original := props.Scroll([]int{1}, "p", 1, nil, 2)
	resetted := original.Reset()

	if original.IsReset() {
		t.Error("original should not be affected by Reset()")
	}

	if !resetted.IsReset() {
		t.Error("reset copy should have IsReset() = true")
	}
}

func TestScrollProp_AllMetadataFields(t *testing.T) {
	t.Parallel()

	p := props.Scroll([]int{1, 2}, "myPage", 5, 4, 6)

	if p.PageName != "myPage" {
		t.Errorf("PageName = %q", p.PageName)
	}

	if p.CurrentPage != 5 {
		t.Errorf("CurrentPage = %v", p.CurrentPage)
	}

	if p.PreviousPage != 4 {
		t.Errorf("PreviousPage = %v", p.PreviousPage)
	}

	if p.NextPage != 6 {
		t.Errorf("NextPage = %v", p.NextPage)
	}

	got, ok := p.Prop().([]int)

	if !ok || len(got) != 2 {
		t.Errorf("Prop() = %v", p.Prop())
	}
}
