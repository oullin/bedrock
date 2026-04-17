package props_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/inertia/props"
)

// No X-Inertia-Partial-Data header — onlySet is empty.

type tryPropValue struct{ val any }

// --- Func Resolution ---

type tryPropError struct{}

func TestResolve_FullRequest(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/users", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	merged := protocol.Props{
		"users":  []string{"alice", "bob"},
		"title":  "Users Page",
		"lazy":   func() any { return "computed" },
		"always": props.Always("always-val"),
	}

	result, err := props.Resolve(r, "Users/Index", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["users"] == nil {
		t.Error("expected 'users' prop")
	}

	if result.Props["title"] != "Users Page" {
		t.Errorf("title = %v", result.Props["title"])
	}

	if result.Props["lazy"] != "computed" {
		t.Errorf("lazy = %v", result.Props["lazy"])
	}

	if result.Props["always"] != "always-val" {
		t.Errorf("always = %v", result.Props["always"])
	}
}

func TestResolve_DeferredExcludedOnInitial(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	merged := protocol.Props{
		"name":  "test",
		"stats": props.Defer(func() any { return "expensive" }, "sidebar"),
	}

	result, err := props.Resolve(r, "Dashboard", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["stats"]; ok {
		t.Error("deferred prop 'stats' should be excluded on initial load")
	}

	if result.Props["name"] != "test" {
		t.Errorf("name = %v", result.Props["name"])
	}

	sidebarGroup := result.DeferredProps["sidebar"]

	if len(sidebarGroup) != 1 || sidebarGroup[0] != "stats" {
		t.Errorf("DeferredProps[sidebar] = %v, want [stats]", sidebarGroup)
	}
}

func TestResolve_DeferredIncludedOnPartialReload(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Dashboard")
	r.Header.Set(protocol.HeaderPartialData, "stats")

	merged := protocol.Props{
		"name":  "test",
		"stats": props.Defer(func() any { return "expensive" }),
	}

	result, err := props.Resolve(r, "Dashboard", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["stats"] != "expensive" {
		t.Errorf("stats = %v, want %q", result.Props["stats"], "expensive")
	}

	if _, ok := result.Props["name"]; ok {
		t.Error("'name' should be excluded from partial reload requesting only 'stats'")
	}
}

func TestResolve_OptionalExcludedOnInitial(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	merged := protocol.Props{
		"name":     "test",
		"optional": props.Optional("opt-val"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["optional"]; ok {
		t.Error("optional prop should be excluded on initial load")
	}
}

func TestResolve_OptionalIncludedOnPartialReload(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "optional")

	merged := protocol.Props{
		"name":     "test",
		"optional": props.Optional("opt-val"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["optional"] != "opt-val" {
		t.Errorf("optional = %v, want %q", result.Props["optional"], "opt-val")
	}
}

func TestResolve_OnceExcludedWhenInExceptHeader(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderExceptOnceProps, "meta")

	merged := protocol.Props{
		"name": "test",
		"meta": props.Once("metadata"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["meta"]; ok {
		t.Error("once prop 'meta' should be excluded when in X-Inertia-Except-Once-Props")
	}
}

func TestResolve_OnceIncludedOnFirstVisit(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	merged := protocol.Props{
		"meta": props.Once("metadata"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["meta"] != "metadata" {
		t.Errorf("meta = %v, want %q", result.Props["meta"], "metadata")
	}
}

func TestResolve_MergeRecorded(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	merged := protocol.Props{
		"posts": props.Merge([]string{"a", "b"}),
		"deep":  props.DeepMerge(map[string]int{"x": 1}),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if len(result.MergeProps) != 1 || result.MergeProps[0] != "posts" {
		t.Errorf("MergeProps = %v, want [posts]", result.MergeProps)
	}

	if len(result.DeepMergeProps) != 1 || result.DeepMergeProps[0] != "deep" {
		t.Errorf("DeepMergeProps = %v, want [deep]", result.DeepMergeProps)
	}
}

func TestResolve_PartialExceptTakesPrecedence(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "a,b,c")
	r.Header.Set(protocol.HeaderPartialExcept, "b")

	merged := protocol.Props{
		"a": "val-a",
		"b": "val-b",
		"c": "val-c",
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["b"]; ok {
		t.Error("prop 'b' should be excluded by X-Inertia-Partial-Except")
	}

	if result.Props["a"] != "val-a" {
		t.Errorf("a = %v", result.Props["a"])
	}

	if result.Props["c"] != "val-c" {
		t.Errorf("c = %v", result.Props["c"])
	}
}

func TestResolve_AlwaysIncludedInPartialReload(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "name")

	merged := protocol.Props{
		"name":   "test",
		"always": props.Always("always-val"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["always"] != "always-val" {
		t.Error("AlwaysProp should be included even when not in partial data list")
	}
}

func TestResolve_LazyFuncWithError(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	merged := protocol.Props{
		"failing": func() (any, error) {
			return nil, http.ErrAbortHandler
		},
	}

	_, err := props.Resolve(r, "Page", merged)

	if err == nil {
		t.Error("expected error from lazy func")
	}
}

func TestResolve_OptionalExcludedOnPartialWithoutOnly(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")

	merged := protocol.Props{
		"name":     "test",
		"optional": props.Optional("opt-val"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["optional"]; ok {
		t.Error("optional prop should be excluded on partial reload without explicit request")
	}

	if result.Props["name"] != "test" {
		t.Errorf("name = %v", result.Props["name"])
	}
}

func TestResolve_DeferredMergeOnInitial(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	merged := protocol.Props{
		"items": props.Defer(func() any { return []string{"a"} }, "list").Merge(),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["items"]; ok {
		t.Error("deferred prop should be excluded on initial load")
	}

	if len(result.DeferredProps["list"]) != 1 || result.DeferredProps["list"][0] != "items" {
		t.Errorf("DeferredProps[list] = %v", result.DeferredProps["list"])
	}

	if len(result.MergeProps) != 1 || result.MergeProps[0] != "items" {
		t.Errorf("MergeProps = %v, want [items]", result.MergeProps)
	}
}

func TestResolve_DeferredMergeOnPartialReload(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "items")

	merged := protocol.Props{
		"items": props.Defer(func() any { return []string{"a"} }, "list").Merge(),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["items"] == nil {
		t.Error("deferred merge prop should be included in partial reload")
	}

	if len(result.MergeProps) != 1 || result.MergeProps[0] != "items" {
		t.Errorf("MergeProps = %v, want [items]", result.MergeProps)
	}
}

func TestResolve_OnceMetadataRecorded(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"notes": props.Once("ship notes"),
	})

	if err != nil {
		t.Fatal(err)
	}

	meta, ok := result.OnceProps["notes"]

	if !ok {
		t.Fatal("missing once metadata")
	}

	if meta.Prop != "notes" {
		t.Errorf("OnceProps[notes].Prop = %q, want %q", meta.Prop, "notes")
	}
}

func TestResolve_OnceMetadataExpiresAt(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"notes": props.Once("ship notes").ExpiresAt(1700000000),
	})

	if err != nil {
		t.Fatal(err)
	}

	meta, ok := result.OnceProps["notes"]

	if !ok {
		t.Fatal("missing once metadata")
	}

	if meta.ExpiresAt == nil || *meta.ExpiresAt != 1700000000 {
		t.Errorf("OnceProps[notes].ExpiresAt = %v, want 1700000000", meta.ExpiresAt)
	}
}

func TestResolve_OnceMetadataExpiresAtNilByDefault(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"notes": props.Once("ship notes"),
	})

	if err != nil {
		t.Fatal(err)
	}

	meta := result.OnceProps["notes"]

	if meta.ExpiresAt != nil {
		t.Errorf("OnceProps[notes].ExpiresAt = %v, want nil", meta.ExpiresAt)
	}
}

func TestResolve_ScrollMetadataRecorded(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Scroll(map[string]any{
			"data": []string{"one"},
		}, "feedPage", 1, nil, 2),
	})

	if err != nil {
		t.Fatal(err)
	}

	meta, ok := result.ScrollProps["feed"]

	if !ok {
		t.Fatal("missing scroll metadata")
	}

	if meta.PageName != "feedPage" {
		t.Errorf("ScrollProps[feed].PageName = %q, want %q", meta.PageName, "feedPage")
	}

	if meta.NextPage != 2 {
		t.Errorf("ScrollProps[feed].NextPage = %v, want %v", meta.NextPage, 2)
	}
}

func (tp tryPropValue) TryProp() (any, error) { return tp.val, nil }

func TestResolve_TryProper(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	merged := protocol.Props{
		"custom": tryPropValue{val: "resolved-via-try"},
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["custom"] != "resolved-via-try" {
		t.Errorf("custom = %v, want %q", result.Props["custom"], "resolved-via-try")
	}
}

func TestResolve_FuncReturningNil(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"nilProp": func() any { return nil },
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["nilProp"]; !ok {
		t.Error("key 'nilProp' should be present even when value is nil")
	}

	if result.Props["nilProp"] != nil {
		t.Errorf("nilProp = %v, want nil", result.Props["nilProp"])
	}
}

func TestResolve_FuncWithErrorReturningNil(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"nilProp": func() (any, error) { return nil, nil },
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["nilProp"]; !ok {
		t.Error("key 'nilProp' should be present")
	}
}

func (tp tryPropError) TryProp() (any, error) {
	return nil, http.ErrAbortHandler
}

func TestResolve_TryProperReturningError(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	_, err := props.Resolve(r, "Page", protocol.Props{
		"bad": tryPropError{},
	})

	if err == nil {
		t.Error("expected error from TryProp()")
	}
}

func TestResolve_TryProperWrappingLazyFunc(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	valFn := func() any { return "from-func" }

	result, err := props.Resolve(r, "Page", protocol.Props{
		"lazy": tryPropValue{val: valFn},
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["lazy"] != "from-func" {
		t.Errorf("lazy = %v, want %q", result.Props["lazy"], "from-func")
	}
}

func TestResolve_ProperChainUnwrapped(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	// AlwaysProp wrapping a value — the resolve loop unwraps Proper.
	result, err := props.Resolve(r, "Page", protocol.Props{
		"wrapped": props.Always("inner-value"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["wrapped"] != "inner-value" {
		t.Errorf("wrapped = %v, want %q", result.Props["wrapped"], "inner-value")
	}
}

// --- Partial Filtering ---

func TestResolve_PartialComponentMismatchBehavesLikeInitial(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "OtherComponent")
	r.Header.Set(protocol.HeaderPartialData, "name")

	merged := protocol.Props{
		"name":     "test",
		"deferred": props.Defer("val"),
		"optional": props.Optional("opt"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	// Component mismatch means isPartial=false, so behave like initial load.
	if result.Props["name"] != "test" {
		t.Error("regular prop should be included on component mismatch")
	}

	if _, ok := result.Props["deferred"]; ok {
		t.Error("deferred prop should be excluded (initial load behavior)")
	}

	if _, ok := result.Props["optional"]; ok {
		t.Error("optional prop should be excluded (initial load behavior)")
	}
}

func TestResolve_NonInertiaRequestBehavesLikeInitial(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	// No X-Inertia header.

	merged := protocol.Props{
		"name":     "test",
		"deferred": props.Defer("val"),
		"optional": props.Optional("opt"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["name"] != "test" {
		t.Error("regular prop should be included")
	}

	if _, ok := result.Props["deferred"]; ok {
		t.Error("deferred prop should be excluded on non-Inertia request")
	}

	if _, ok := result.Props["optional"]; ok {
		t.Error("optional prop should be excluded on non-Inertia request")
	}
}

func TestResolve_PartialExceptWithoutPartialData(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialExcept, "b")

	merged := protocol.Props{
		"a": "val-a",
		"b": "val-b",
		"c": "val-c",
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["b"]; ok {
		t.Error("'b' should be excluded by except header")
	}

	if result.Props["a"] != "val-a" {
		t.Errorf("a = %v", result.Props["a"])
	}

	if result.Props["c"] != "val-c" {
		t.Errorf("c = %v", result.Props["c"])
	}
}

func TestResolve_EmptyPartialDataHeader(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "")

	merged := protocol.Props{
		"a": "val-a",
		"b": "val-b",
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	// Empty partial data = onlySet is nil, so no "only" filtering.
	if result.Props["a"] != "val-a" {
		t.Errorf("a = %v", result.Props["a"])
	}

	if result.Props["b"] != "val-b" {
		t.Errorf("b = %v", result.Props["b"])
	}
}

func TestResolve_PartialDataWithWhitespace(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, " a , b ")

	merged := protocol.Props{
		"a": "val-a",
		"b": "val-b",
		"c": "val-c",
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["a"] != "val-a" {
		t.Errorf("a = %v", result.Props["a"])
	}

	if result.Props["b"] != "val-b" {
		t.Errorf("b = %v", result.Props["b"])
	}

	if _, ok := result.Props["c"]; ok {
		t.Error("'c' should be excluded from partial with whitespace-trimmed header")
	}
}

// --- AlwaysProp in Partial ---

func TestResolve_AlwaysNotExcludedByExceptHeader(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialExcept, "always")

	merged := protocol.Props{
		"always": props.Always("always-val"),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["always"] != "always-val" {
		t.Error("AlwaysProp should not be excludable by except header")
	}
}

func TestResolve_AlwaysWithScalarTypes(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"str":  props.Always("hello"),
		"num":  props.Always(42),
		"flag": props.Always(true),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["str"] != "hello" {
		t.Errorf("str = %v", result.Props["str"])
	}

	if result.Props["num"] != 42 {
		t.Errorf("num = %v", result.Props["num"])
	}

	if result.Props["flag"] != true {
		t.Errorf("flag = %v", result.Props["flag"])
	}
}

// --- MergeProp Filtering ---

func TestResolve_MergeExcludedByPartialExcept(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialExcept, "items")

	merged := protocol.Props{
		"items": props.Merge([]string{"a"}),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["items"]; ok {
		t.Error("MergeProp should be excluded by except header")
	}

	if len(result.MergeProps) > 0 {
		t.Errorf("MergeProps = %v, want empty", result.MergeProps)
	}
}

func TestResolve_MergeExcludedByPartialOnly(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "other")

	merged := protocol.Props{
		"other": "val",
		"items": props.Merge([]string{"a"}),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["items"]; ok {
		t.Error("MergeProp should be excluded when not in partial data")
	}
}

func TestResolve_DeepMergeExcludedByPartialExcept(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialExcept, "deep")

	merged := protocol.Props{
		"deep": props.DeepMerge(map[string]int{"x": 1}),
	}

	result, err := props.Resolve(r, "Page", merged)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["deep"]; ok {
		t.Error("DeepMergeProp should be excluded by except header")
	}

	if len(result.DeepMergeProps) > 0 {
		t.Errorf("DeepMergeProps = %v, want empty", result.DeepMergeProps)
	}
}

func TestResolve_MergeOnNonPartialInertiaRequest(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"items": props.Merge([]string{"a", "b"}),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["items"] == nil {
		t.Error("MergeProp should be included on non-partial request")
	}

	if len(result.MergeProps) != 1 || result.MergeProps[0] != "items" {
		t.Errorf("MergeProps = %v, want [items]", result.MergeProps)
	}
}

func TestResolve_MergeAndDeepMergeSimultaneous(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"shallow": props.Merge([]int{1}),
		"deep":    props.DeepMerge(map[string]int{"a": 1}),
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(result.MergeProps) != 1 || result.MergeProps[0] != "shallow" {
		t.Errorf("MergeProps = %v, want [shallow]", result.MergeProps)
	}

	if len(result.DeepMergeProps) != 1 || result.DeepMergeProps[0] != "deep" {
		t.Errorf("DeepMergeProps = %v, want [deep]", result.DeepMergeProps)
	}
}

// --- OnceProp Extended ---

func TestResolve_OnceIncludedOnPartialReload(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "meta")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"meta": props.Once("metadata"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["meta"] != "metadata" {
		t.Errorf("meta = %v, want %q", result.Props["meta"], "metadata")
	}
}

func TestResolve_OnceExcludedByPartialOnly(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "other")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"other": "val",
		"meta":  props.Once("metadata"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["meta"]; ok {
		t.Error("OnceProp should be excluded when not in partial data")
	}
}

func TestResolve_OnceMetadataRecordedWhenExcludedByExceptOnce(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderExceptOnceProps, "meta")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"meta": props.Once("metadata"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["meta"]; ok {
		t.Error("OnceProp value should be excluded")
	}

	// Metadata IS still recorded even when value is excluded.
	if _, ok := result.OnceProps["meta"]; !ok {
		t.Error("OnceProps metadata should still be recorded when excluded by except-once header")
	}
}

func TestResolve_OnceMultiplePropsPartiallyExcluded(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderExceptOnceProps, "cached")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"fresh":  props.Once("fresh-val"),
		"cached": props.Once("cached-val"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["fresh"] != "fresh-val" {
		t.Error("fresh OnceProp should be included")
	}

	if _, ok := result.Props["cached"]; ok {
		t.Error("cached OnceProp should be excluded")
	}
}

// --- ScrollProp Extended ---

func TestResolve_ScrollWithMergeIntentHeader(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderInfiniteScroll, "append")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Scroll([]int{1, 2}, "page", 1, nil, 2),
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(result.MergeProps) != 1 || result.MergeProps[0] != "feed" {
		t.Errorf("MergeProps = %v, want [feed] (triggered by merge intent header)", result.MergeProps)
	}
}

func TestResolve_ScrollWithoutMergeIntentOrMerge(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Scroll([]int{1, 2}, "page", 1, nil, 2),
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(result.MergeProps) > 0 {
		t.Errorf("MergeProps = %v, want empty (no merge intent or .Merge())", result.MergeProps)
	}
}

func TestResolve_ScrollResetFlagInMetadata(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Scroll([]int{1}, "page", 1, nil, 2).Reset(),
	})

	if err != nil {
		t.Fatal(err)
	}

	meta, ok := result.ScrollProps["feed"]

	if !ok {
		t.Fatal("missing scroll metadata")
	}

	if !meta.Reset {
		t.Error("ScrollMeta.Reset should be true")
	}
}

func TestResolve_ScrollExcludedByPartialExcept(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialExcept, "feed")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Scroll([]int{1}, "page", 1, nil, 2),
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["feed"]; ok {
		t.Error("ScrollProp should be excluded by except header")
	}

	if _, ok := result.ScrollProps["feed"]; ok {
		t.Error("ScrollProps metadata should not be recorded when excluded")
	}
}

// --- DeferProp Extended ---

func TestResolve_DeferredMultipleGroupsOnInitial(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"stats":    props.Defer("s", "sidebar"),
		"forecast": props.Defer("f", "sidebar"),
		"logs":     props.Defer("l", "footer"),
	})

	if err != nil {
		t.Fatal(err)
	}

	sidebar := result.DeferredProps["sidebar"]

	if len(sidebar) != 2 {
		t.Errorf("DeferredProps[sidebar] = %v, want 2 items", sidebar)
	}

	footer := result.DeferredProps["footer"]

	if len(footer) != 1 || footer[0] != "logs" {
		t.Errorf("DeferredProps[footer] = %v, want [logs]", footer)
	}
}

func TestResolve_DeferredExcludedByPartialExcept(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "stats")
	r.Header.Set(protocol.HeaderPartialExcept, "stats")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"stats": props.Defer(func() any { return "expensive" }),
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["stats"]; ok {
		t.Error("DeferProp should be excluded when in except header")
	}
}

func TestResolve_DeferredWithFuncValue(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "stats")

	called := false

	result, err := props.Resolve(r, "Page", protocol.Props{
		"stats": props.Defer(func() any {
			called = true

			return []string{"a", "b"}
		}),
	})

	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error("inner func should be invoked on partial reload")
	}

	if result.Props["stats"] == nil {
		t.Error("stats should contain resolved value")
	}
}

// --- Combined Interactions ---

func TestResolve_AllPropTypesOnInitial(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"regular":  "val",
		"always":   props.Always("always-val"),
		"optional": props.Optional("opt-val"),
		"deferred": props.Defer("defer-val"),
		"once":     props.Once("once-val"),
		"merge":    props.Merge([]int{1}),
		"scroll":   props.Scroll([]int{1}, "p", 1, nil, 2),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["regular"] != "val" {
		t.Error("regular prop should be included")
	}

	if result.Props["always"] != "always-val" {
		t.Error("AlwaysProp should be included")
	}

	if _, ok := result.Props["optional"]; ok {
		t.Error("OptionalProp should be excluded on initial")
	}

	if _, ok := result.Props["deferred"]; ok {
		t.Error("DeferProp should be excluded on initial")
	}

	if result.Props["once"] != "once-val" {
		t.Error("OnceProp should be included on initial")
	}

	if result.Props["merge"] == nil {
		t.Error("MergeProp should be included on initial")
	}

	if result.Props["scroll"] == nil {
		t.Error("ScrollProp should be included on initial")
	}

	if len(result.DeferredProps) == 0 {
		t.Error("DeferredProps metadata should be recorded")
	}

	if len(result.MergeProps) == 0 {
		t.Error("MergeProps should contain 'merge'")
	}

	if len(result.OnceProps) == 0 {
		t.Error("OnceProps metadata should be recorded")
	}

	if len(result.ScrollProps) == 0 {
		t.Error("ScrollProps metadata should be recorded")
	}
}

func TestResolve_AllPropTypesOnPartialReload(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "regular,deferred,optional,once,merge,scroll,always")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"regular":  "val",
		"always":   props.Always("always-val"),
		"optional": props.Optional("opt-val"),
		"deferred": props.Defer("defer-val"),
		"once":     props.Once("once-val"),
		"merge":    props.Merge([]int{1}),
		"scroll":   props.Scroll([]int{1}, "p", 1, nil, 2),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["regular"] != "val" {
		t.Error("regular prop should be included on partial")
	}

	if result.Props["always"] != "always-val" {
		t.Error("AlwaysProp should be included on partial")
	}

	if result.Props["optional"] != "opt-val" {
		t.Error("OptionalProp should be included when explicitly requested")
	}

	if result.Props["deferred"] != "defer-val" {
		t.Error("DeferProp should be included on partial reload")
	}

	if result.Props["once"] != "once-val" {
		t.Error("OnceProp should be included on partial")
	}

	if result.Props["merge"] == nil {
		t.Error("MergeProp should be included on partial")
	}

	if result.Props["scroll"] == nil {
		t.Error("ScrollProp should be included on partial")
	}
}

func TestResolve_MixedExceptAndOnceExcept(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialExcept, "excluded")
	r.Header.Set(protocol.HeaderExceptOnceProps, "cached_once")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"included":    "val",
		"excluded":    "val",
		"fresh_once":  props.Once("fresh"),
		"cached_once": props.Once("cached"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["included"] != "val" {
		t.Error("included should be present")
	}

	if _, ok := result.Props["excluded"]; ok {
		t.Error("excluded should be removed by partial except")
	}

	if result.Props["fresh_once"] != "fresh" {
		t.Error("fresh_once should be included")
	}

	if _, ok := result.Props["cached_once"]; ok {
		t.Error("cached_once should be excluded by except-once header")
	}
}

// --- Nested Wrapper Tests ---

func TestResolve_ScrollWrappingDefer_InitialLoad(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Scroll(
			props.Defer(func() any { return []int{1, 2} }, "sidebar"),
			"feedPage", 1, nil, 2,
		),
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["feed"]; ok {
		t.Error("Scroll(Defer(...)) should be excluded on initial load")
	}

	if len(result.DeferredProps["sidebar"]) != 1 || result.DeferredProps["sidebar"][0] != "feed" {
		t.Errorf("DeferredProps[sidebar] = %v, want [feed]", result.DeferredProps["sidebar"])
	}

	if _, ok := result.ScrollProps["feed"]; ok {
		t.Error("scroll metadata should not be recorded when prop is deferred on initial")
	}
}

func TestResolve_ScrollWrappingDefer_PartialReload(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "feed")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Scroll(
			props.Defer(func() any { return []int{1, 2} }, "sidebar"),
			"feedPage", 1, nil, 2,
		),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["feed"] == nil {
		t.Error("Scroll(Defer(...)) should be included on partial reload")
	}

	meta, ok := result.ScrollProps["feed"]

	if !ok {
		t.Fatal("scroll metadata should be recorded on partial reload")
	}

	if meta.PageName != "feedPage" {
		t.Errorf("ScrollProps[feed].PageName = %q, want %q", meta.PageName, "feedPage")
	}

	if meta.NextPage != 2 {
		t.Errorf("ScrollProps[feed].NextPage = %v, want 2", meta.NextPage)
	}
}

func TestResolve_DeferWrappingScroll_InitialLoad(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Defer(
			props.Scroll([]int{1, 2}, "feedPage", 1, nil, 2),
			"sidebar",
		),
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["feed"]; ok {
		t.Error("Defer(Scroll(...)) should be excluded on initial load")
	}

	if len(result.DeferredProps["sidebar"]) != 1 || result.DeferredProps["sidebar"][0] != "feed" {
		t.Errorf("DeferredProps[sidebar] = %v, want [feed]", result.DeferredProps["sidebar"])
	}

	if _, ok := result.ScrollProps["feed"]; ok {
		t.Error("scroll metadata should not be recorded when prop is deferred on initial")
	}
}

func TestResolve_DeferWrappingScroll_PartialReload(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "feed")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Defer(
			props.Scroll([]int{1, 2}, "feedPage", 1, nil, 2),
			"sidebar",
		),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["feed"] == nil {
		t.Error("Defer(Scroll(...)) should be included on partial reload")
	}

	meta, ok := result.ScrollProps["feed"]

	if !ok {
		t.Fatal("scroll metadata should be recorded on partial reload")
	}

	if meta.PageName != "feedPage" {
		t.Errorf("ScrollProps[feed].PageName = %q, want %q", meta.PageName, "feedPage")
	}
}

func TestResolve_ScrollWrappingDeferMerge_InitialLoad(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Scroll(
			props.Defer(func() any { return []int{1} }, "g").Merge(),
			"p", 1, nil, 2,
		),
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["feed"]; ok {
		t.Error("should be excluded on initial")
	}

	if len(result.DeferredProps["g"]) != 1 {
		t.Errorf("DeferredProps[g] = %v, want [feed]", result.DeferredProps["g"])
	}

	if len(result.MergeProps) != 1 || result.MergeProps[0] != "feed" {
		t.Errorf("MergeProps = %v, want [feed]", result.MergeProps)
	}
}

func TestResolve_DeferWrappingScrollMerge_PartialWithIntent(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "feed")
	r.Header.Set(protocol.HeaderInfiniteScroll, "append")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Defer(
			props.Scroll([]int{1}, "p", 1, nil, 2),
			"g",
		),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["feed"] == nil {
		t.Error("should be included on partial reload")
	}

	if _, ok := result.ScrollProps["feed"]; !ok {
		t.Error("scroll metadata should be recorded")
	}

	if len(result.MergeProps) != 1 || result.MergeProps[0] != "feed" {
		t.Errorf("MergeProps = %v, want [feed] (from merge intent header)", result.MergeProps)
	}
}

func TestResolve_AlwaysWrappingDefer_InitialLoad(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"info": props.Always(props.Defer("val", "g")),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["info"] != "val" {
		t.Errorf("Always(Defer(...)) should be included; got %v", result.Props["info"])
	}

	if len(result.DeferredProps) > 0 {
		t.Error("deferred group should not be recorded when Always overrides")
	}
}

func TestResolve_ScrollWrappingOnce_ExcludedByExceptOnce(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderExceptOnceProps, "feed")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Scroll(props.Once("val"), "p", 1, nil, 2),
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["feed"]; ok {
		t.Error("should be excluded by except-once header")
	}

	if _, ok := result.OnceProps["feed"]; !ok {
		t.Error("once metadata should still be recorded")
	}

	if _, ok := result.ScrollProps["feed"]; ok {
		t.Error("scroll metadata should not be recorded when excluded")
	}
}

func TestResolve_TripleNesting_DeferScrollMerge_Initial(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Defer(
			props.Scroll(props.Merge([]int{1}), "p", 1, nil, 2),
			"g",
		),
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["feed"]; ok {
		t.Error("should be excluded on initial")
	}

	if len(result.DeferredProps["g"]) != 1 {
		t.Errorf("DeferredProps[g] = %v, want [feed]", result.DeferredProps["g"])
	}
}

func TestResolve_TripleNesting_DeferScrollMerge_Partial(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "feed")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Defer(
			props.Scroll(props.Merge([]int{1}), "p", 1, nil, 2),
			"g",
		),
	})

	if err != nil {
		t.Fatal(err)
	}

	if result.Props["feed"] == nil {
		t.Error("should be included on partial")
	}

	if _, ok := result.ScrollProps["feed"]; !ok {
		t.Error("scroll metadata should be recorded on partial")
	}

	if len(result.MergeProps) < 1 {
		t.Error("merge metadata should be recorded on partial")
	}
}

func TestResolve_NestedWithPartialExcept(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderPartialComponent, "Page")
	r.Header.Set(protocol.HeaderPartialData, "feed")
	r.Header.Set(protocol.HeaderPartialExcept, "feed")

	result, err := props.Resolve(r, "Page", protocol.Props{
		"feed": props.Scroll(
			props.Defer(func() any { return []int{1} }, "g"),
			"p", 1, nil, 2,
		),
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result.Props["feed"]; ok {
		t.Error("should be excluded by partial except")
	}

	if _, ok := result.ScrollProps["feed"]; ok {
		t.Error("no metadata should be recorded when excluded by partial except")
	}

	if len(result.DeferredProps) > 0 {
		t.Error("no deferred metadata should be recorded when excluded by partial except")
	}
}
