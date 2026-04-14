package validation_test

import (
	"testing"

	"github.com/bedrock/packages/validation"
)

func TestStudlyCase(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want string
	}{
		{"required", "Required"},
		{"required_if", "RequiredIf"},
		{"required_without_all", "RequiredWithoutAll"},
		{"max", "Max"},
		{"not_in", "NotIn"},
		{"alpha_dash", "AlphaDash"},
		{"Min", "Min"},
	}

	for _, tc := range cases {
		got := validation.StudlyCase(tc.in)
		if got != tc.want {
			t.Errorf("StudlyCase(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		rule       string
		wantName   string
		wantParams []string
	}{
		{"required", "Required", nil},
		{"max:255", "Max", []string{"255"}},
		{"required_if:field,value", "RequiredIf", []string{"field", "value"}},
		{"between:1,100", "Between", []string{"1", "100"}},
		{"regex:/^[a-z]+$/", "Regex", []string{"/^[a-z]+$/"}},
		{"not_regex:/[0-9]/", "NotRegex", []string{"/[0-9]/"}},
	}

	for _, tc := range cases {
		got := validation.Parse(tc.rule)

		if got.Name != tc.wantName {
			t.Errorf("Parse(%q).Name = %q, want %q", tc.rule, got.Name, tc.wantName)
		}

		if len(got.Parameters) != len(tc.wantParams) {
			t.Errorf("Parse(%q).Parameters = %v, want %v", tc.rule, got.Parameters, tc.wantParams)
			continue
		}

		for i, p := range tc.wantParams {
			if got.Parameters[i] != p {
				t.Errorf("Parse(%q).Parameters[%d] = %q, want %q", tc.rule, i, got.Parameters[i], p)
			}
		}
	}
}

func TestExplode_String(t *testing.T) {
	t.Parallel()

	rules := validation.Explode("required|email|max:255")

	if len(rules) != 3 {
		t.Fatalf("Explode: got %d rules, want 3", len(rules))
	}

	if rules[0].Name != "Required" {
		t.Errorf("rules[0].Name = %q", rules[0].Name)
	}

	if rules[1].Name != "Email" {
		t.Errorf("rules[1].Name = %q", rules[1].Name)
	}

	if rules[2].Name != "Max" || rules[2].Parameters[0] != "255" {
		t.Errorf("rules[2] = %v", rules[2])
	}
}

func TestExplode_StringSlice(t *testing.T) {
	t.Parallel()

	rules := validation.Explode([]string{"required", "email"})

	if len(rules) != 2 {
		t.Fatalf("Explode([]string): got %d rules, want 2", len(rules))
	}
}

func TestExpandWildcards(t *testing.T) {
	t.Parallel()

	flat := map[string]any{
		"items.0.name": "Alice",
		"items.1.name": "Bob",
		"items.0.age":  30,
		"other":        "value",
	}

	matches := validation.ExpandWildcards("items.*.name", flat)

	if len(matches) != 2 {
		t.Fatalf("ExpandWildcards: got %d matches, want 2: %v", len(matches), matches)
	}
}

func TestExpandWildcards_NoWildcard(t *testing.T) {
	t.Parallel()

	flat := map[string]any{"name": "Alice"}
	matches := validation.ExpandWildcards("name", flat)

	if len(matches) != 1 || matches[0] != "name" {
		t.Errorf("ExpandWildcards(no wildcard): got %v", matches)
	}
}

func TestFlattenData(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
		},
		"items": []any{
			map[string]any{"id": 1},
			map[string]any{"id": 2},
		},
	}

	flat := validation.FlattenData(data)

	if flat["user.name"] != "Alice" {
		t.Errorf("user.name: got %v", flat["user.name"])
	}

	if flat["items.0.id"] != 1 {
		t.Errorf("items.0.id: got %v", flat["items.0.id"])
	}

	if flat["items.1.id"] != 2 {
		t.Errorf("items.1.id: got %v", flat["items.1.id"])
	}
}
