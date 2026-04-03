package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestBuilderMergesBaseOverlayAndEnv(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "auth.yml"), "session_lifetime: 24h\ncookies:\n  secure: true\n")
	writeFile(t, filepath.Join(dir, "authflows.yml"), "views: true\nlimiters:\n  login: login\nfeatures:\n  - registration\n")
	writeFile(t, filepath.Join(dir, "local", "auth.yml"), "session_lifetime: 12h\ncookies:\n  secure: false\n")

	repo, err := NewBuilder(dir).
		WithAppEnv("local").
		WithEnv(map[string]string{
			"AUTH_SESSION_LIFETIME": "48h",
			"FORTIFY_VIEWS":         "false",
		}).
		Build(context.Background())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	duration, err := repo.Duration("auth.session_lifetime")
	if err != nil {
		t.Fatalf("Duration: %v", err)
	}
	if duration.String() != "48h0m0s" {
		t.Fatalf("unexpected duration: %v", duration)
	}

	secure, err := repo.Bool("auth.cookies.secure")
	if err != nil {
		t.Fatalf("Bool: %v", err)
	}
	if secure {
		t.Fatal("expected overlay bool to apply")
	}

	views, err := repo.Bool("authflows.views")
	if err != nil {
		t.Fatalf("Bool: %v", err)
	}
	if views {
		t.Fatal("expected env override to disable views")
	}
}

func TestBuilderUsesAPPENVWhenWithAppEnvIsUnset(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "app.yml"), "name: base\n")
	writeFile(t, filepath.Join(dir, "production", "app.yml"), "name: production\n")

	t.Setenv("APP_ENV", "production")

	repo, err := NewBuilder(dir).WithEnv(nil).Build(context.Background())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	name, err := repo.String("app.name")
	if err != nil {
		t.Fatalf("String: %v", err)
	}
	if name != "production" {
		t.Fatalf("unexpected app name: %q", name)
	}
}

func TestBuilderMissingDirectoriesAndInvalidPaths(t *testing.T) {
	t.Parallel()

	missing := filepath.Join(t.TempDir(), "missing")
	repo, err := NewBuilder(missing).Build(context.Background())
	if err != nil {
		t.Fatalf("Build missing dir: %v", err)
	}
	if got := repo.All(); len(got) != 0 {
		t.Fatalf("expected empty repository for missing dir, got %#v", got)
	}

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "app.yml"), "name: base\n")
	repo, err = NewBuilder(dir).WithAppEnv("missing-overlay").Build(context.Background())
	if err != nil {
		t.Fatalf("Build missing overlay: %v", err)
	}
	name, err := repo.String("app.name")
	if err != nil {
		t.Fatalf("String: %v", err)
	}
	if name != "base" {
		t.Fatalf("unexpected base value with missing overlay: %q", name)
	}

	file := filepath.Join(t.TempDir(), "not-a-dir.yml")
	if err := os.WriteFile(file, []byte("name: value\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err = NewBuilder(file).Build(context.Background())
	if err == nil {
		t.Fatal("expected invalid path error")
	}
	if !strings.Contains(err.Error(), "is not a directory") {
		t.Fatalf("unexpected invalid path error: %v", err)
	}

	overlayFile := filepath.Join(dir, "broken-overlay")
	if err := os.WriteFile(overlayFile, []byte("name: value\n"), 0o644); err != nil {
		t.Fatalf("WriteFile overlay: %v", err)
	}

	_, err = NewBuilder(dir).WithAppEnv("broken-overlay").Build(context.Background())
	if err == nil {
		t.Fatal("expected invalid overlay path error")
	}
	if !strings.Contains(err.Error(), "is not a directory") {
		t.Fatalf("unexpected invalid overlay error: %v", err)
	}
}

func TestBuilderInvalidYAML(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "auth.yml"), "session_lifetime: [\n")

	_, err := NewBuilder(dir).Build(context.Background())
	if err == nil {
		t.Fatal("expected invalid yaml error")
	}
}

func TestBuilderNormalizeAndCoerceEnvValues(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "types.yml"), strings.Join([]string{
		"boolean_true: true",
		"integer_value: 1",
		"float_value: 1.5",
		"int64_like: 6",
		"tags:",
		"  - alpha",
		"nested:",
		"  providers:",
		"    foo: bar",
		"  list:",
		"    - key: one",
		"      value: two",
		"raw_string: base",
	}, "\n")+"\n")

	repo, err := NewBuilder(dir).WithEnv(map[string]string{
		"TYPES_BOOLEAN_TRUE":  "false",
		"TYPES_INTEGER_VALUE": "9",
		"TYPES_FLOAT_VALUE":   "3.25",
		"TYPES_INT64_LIKE":    "11",
		"TYPES_TAGS":          "gamma, delta",
		"TYPES_RAW_STRING":    "override",
	}).Build(context.Background())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	booleanValue, err := repo.Bool("types.boolean_true")
	if err != nil {
		t.Fatalf("Bool: %v", err)
	}
	if booleanValue {
		t.Fatal("expected bool env override to coerce to false")
	}

	integerValue, err := repo.Int("types.integer_value")
	if err != nil {
		t.Fatalf("Int integer value: %v", err)
	}
	if integerValue != 9 {
		t.Fatalf("unexpected integer env coercion: %d", integerValue)
	}

	floatValue, err := repo.Get("types.float_value", nil).(float64), error(nil)
	if err != nil {
		t.Fatalf("unexpected float retrieval error: %v", err)
	}
	if floatValue != 3.25 {
		t.Fatalf("unexpected float env coercion: %v", floatValue)
	}

	int64Like, err := repo.Int("types.int64_like")
	if err != nil {
		t.Fatalf("Int int64-like value: %v", err)
	}
	if int64Like != 11 {
		t.Fatalf("unexpected int64-like env coercion: %d", int64Like)
	}

	tags, err := repo.StringSlice("types.tags")
	if err != nil {
		t.Fatalf("StringSlice tags: %v", err)
	}
	if !reflect.DeepEqual(tags, []string{"gamma", "delta"}) {
		t.Fatalf("unexpected tag env coercion: %#v", tags)
	}

	rawString, err := repo.String("types.raw_string")
	if err != nil {
		t.Fatalf("String raw_string: %v", err)
	}
	if rawString != "override" {
		t.Fatalf("unexpected raw string override: %q", rawString)
	}

	nested, err := repo.Map("types.nested")
	if err != nil {
		t.Fatalf("Map nested: %v", err)
	}
	providers, ok := nested["providers"].(map[string]any)
	if !ok || providers["foo"] != "bar" {
		t.Fatalf("unexpected normalized nested providers: %#v", nested)
	}
	list, ok := nested["list"].([]any)
	if !ok || len(list) != 1 {
		t.Fatalf("unexpected normalized nested list: %#v", nested)
	}
	first, ok := list[0].(map[string]any)
	if !ok || first["key"] != "one" || first["value"] != "two" {
		t.Fatalf("unexpected normalized list item: %#v", list[0])
	}
}

func TestBuilderHelpers(t *testing.T) {
	builder := NewBuilder(".")
	if got := builder.WithAppEnv(" local "); got != builder {
		t.Fatal("expected WithAppEnv to return same builder")
	}
	if builder.appEnv != "local" {
		t.Fatalf("unexpected trimmed appEnv: %q", builder.appEnv)
	}
	if got := builder.WithEnv(map[string]string{"APP_ENV": "testing"}); got != builder {
		t.Fatal("expected WithEnv to return same builder")
	}
	if builder.env["APP_ENV"] != "testing" {
		t.Fatalf("unexpected copied env map: %#v", builder.env)
	}

	clonedEnv := map[string]string{"APP_ENV": "local"}
	builder.WithEnv(clonedEnv)
	clonedEnv["APP_ENV"] = "mutated"
	if builder.env["APP_ENV"] != "local" {
		t.Fatalf("expected WithEnv to copy input map, got %#v", builder.env)
	}

	builder.WithEnv(nil)
	t.Setenv("BUILDER_LOOKUP_ENV", "from-process")
	if got := builder.lookupEnv("BUILDER_LOOKUP_ENV"); got != "from-process" {
		t.Fatalf("unexpected process env lookup: %q", got)
	}

	builder.WithEnv(map[string]string{"BUILDER_LOOKUP_ENV": "from-builder"})
	if got := builder.lookupEnv("BUILDER_LOOKUP_ENV"); got != "from-builder" {
		t.Fatalf("unexpected builder env lookup: %q", got)
	}

	items := map[string]any{
		"types": map[string]any{
			"boolean_true":  true,
			"integer_value": 1,
			"float_value":   float64(1.5),
			"tags":          []string{"alpha", "beta"},
			"raw_string":    "base",
			"nested": map[string]any{
				"value": "one",
			},
		},
	}
	applyEnvOverrides(items, func(key string) string {
		switch key {
		case "TYPES_BOOLEAN_TRUE":
			return "false"
		case "TYPES_INTEGER_VALUE":
			return "2"
		case "TYPES_FLOAT_VALUE":
			return "9.5"
		case "TYPES_TAGS":
			return "gamma,delta"
		case "TYPES_RAW_STRING":
			return "override"
		default:
			return ""
		}
	})

	typed := items["types"].(map[string]any)
	if typed["boolean_true"] != false {
		t.Fatalf("unexpected bool override: %#v", typed["boolean_true"])
	}
	if typed["integer_value"] != 2 {
		t.Fatalf("unexpected int override: %#v", typed["integer_value"])
	}
	if typed["float_value"] != 9.5 {
		t.Fatalf("unexpected float override: %#v", typed["float_value"])
	}
	if !reflect.DeepEqual(typed["tags"], []string{"gamma", "delta"}) {
		t.Fatalf("unexpected slice override: %#v", typed["tags"])
	}
	if typed["raw_string"] != "override" {
		t.Fatalf("unexpected raw string override: %#v", typed["raw_string"])
	}

	keys := flattenedKeys(items, "")
	if !reflect.DeepEqual(keys, []string{
		"types.boolean_true",
		"types.float_value",
		"types.integer_value",
		"types.nested.value",
		"types.raw_string",
		"types.tags",
	}) {
		t.Fatalf("unexpected flattened keys: %#v", keys)
	}

	if got := getByPath(items, "types.nested.value"); got != "one" {
		t.Fatalf("unexpected getByPath value: %#v", got)
	}
	if got := getByPath(items, "types.missing.value"); got != nil {
		t.Fatalf("expected nil for missing path, got %#v", got)
	}
	if got := getByPath(items, "types.raw_string.value"); got != nil {
		t.Fatalf("expected nil for non-map intermediate path, got %#v", got)
	}

	setByPath(items, "types.nested.extra", "two")
	nested := items["types"].(map[string]any)["nested"].(map[string]any)
	if nested["extra"] != "two" {
		t.Fatalf("unexpected setByPath nested result: %#v", nested)
	}

	setByPath(items, "types.raw_string.deeper", "replaced")
	replaced := items["types"].(map[string]any)["raw_string"].(map[string]any)
	if replaced["deeper"] != "replaced" {
		t.Fatalf("unexpected setByPath replacement result: %#v", replaced)
	}

	if got := coerceEnvValue(true, "false"); got != false {
		t.Fatalf("unexpected bool coercion: %#v", got)
	}
	if got := coerceEnvValue(1, "7"); got != 7 {
		t.Fatalf("unexpected int coercion: %#v", got)
	}
	if got := coerceEnvValue(int64(1), "8"); got != int64(8) {
		t.Fatalf("unexpected int64 coercion: %#v", got)
	}
	if got := coerceEnvValue(float64(1.5), "4.25"); got != 4.25 {
		t.Fatalf("unexpected float coercion: %#v", got)
	}
	if got := coerceEnvValue([]string{"x"}, "a,b"); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("unexpected []string coercion: %#v", got)
	}
	if got := coerceEnvValue([]any{"x"}, "c,d"); !reflect.DeepEqual(got, []string{"c", "d"}) {
		t.Fatalf("unexpected []any coercion: %#v", got)
	}
	if got := coerceEnvValue(true, "not-bool"); got != "not-bool" {
		t.Fatalf("expected failed bool parse to keep raw string, got %#v", got)
	}

	normalized := normalizeValue(map[any]any{
		"nested": map[any]any{
			"key": "value",
		},
		"list": []any{map[any]any{"foo": "bar"}},
	}).(map[string]any)
	if normalized["nested"].(map[string]any)["key"] != "value" {
		t.Fatalf("unexpected normalized nested map: %#v", normalized)
	}
	if normalized["list"].([]any)[0].(map[string]any)["foo"] != "bar" {
		t.Fatalf("unexpected normalized nested list entry: %#v", normalized)
	}

	merged := mergeValue(
		map[string]any{"a": map[string]any{"b": "base", "c": "keep"}},
		map[string]any{"a": map[string]any{"b": "override"}},
	).(map[string]any)
	mergedA := merged["a"].(map[string]any)
	if mergedA["b"] != "override" || mergedA["c"] != "keep" {
		t.Fatalf("unexpected mergeValue result: %#v", merged)
	}
	if got := mergeValue("base", "override"); got != "override" {
		t.Fatalf("unexpected scalar mergeValue result: %#v", got)
	}
}

func TestBuilderLoadDirDirectly(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "app.yml"), "name: yaml\n")
	writeFile(t, filepath.Join(dir, "mail.yaml"), "driver: smtp\n")
	writeFile(t, filepath.Join(dir, "notes.txt"), "ignored")
	writeFile(t, filepath.Join(dir, "nested", "auth.yml"), "driver: session\n")

	target := map[string]any{}
	builder := NewBuilder(dir)
	if err := builder.loadDir(target, dir); err != nil {
		t.Fatalf("loadDir: %v", err)
	}

	if target["app"].(map[string]any)["name"] != "yaml" {
		t.Fatalf("unexpected app namespace: %#v", target["app"])
	}
	if target["mail"].(map[string]any)["driver"] != "smtp" {
		t.Fatalf("unexpected mail namespace: %#v", target["mail"])
	}
	if _, ok := target["notes"]; ok {
		t.Fatalf("expected non-YAML file to be ignored: %#v", target)
	}
	if _, ok := target["nested"]; ok {
		t.Fatalf("expected nested directory to be ignored: %#v", target)
	}
}

func TestBuilderLoadDirErrors(t *testing.T) {
	invalidPath := filepath.Join(t.TempDir(), "bad") + string([]byte{0})
	err := NewBuilder(".").loadDir(map[string]any{}, invalidPath)
	if err == nil {
		t.Fatal("expected invalid-path stat error")
	}
	if !strings.Contains(err.Error(), "stat config dir") {
		t.Fatalf("unexpected stat error: %v", err)
	}

	dir := filepath.Join(t.TempDir(), "restricted")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	writeFile(t, filepath.Join(dir, "app.yml"), "name: hidden\n")
	if err := os.Chmod(dir, 0o000); err != nil {
		t.Fatalf("Chmod: %v", err)
	}
	defer func() {
		if chmodErr := os.Chmod(dir, 0o755); chmodErr != nil {
			panic(fmt.Sprintf("restore directory mode: %v", chmodErr))
		}
	}()

	err = NewBuilder(".").loadDir(map[string]any{}, dir)
	if err == nil {
		t.Fatal("expected readDir error")
	}
	if !strings.Contains(err.Error(), "read config dir") {
		t.Fatalf("unexpected readDir error: %v", err)
	}
}

func writeFile(t *testing.T, path string, contents string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}
