package prompts

import "testing"

func TestTableDisplays(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Table(
		[]string{"Name", "Age"},
		[][]string{
			{"Alice", "30"},
			{"Bob", "25"},
		},
	)
	stripped := tp.StrippedContent()

	if !containsAll(stripped, "Name", "Age", "Alice", "30", "Bob", "25") {
		t.Fatalf("expected table content, got:\n%s", stripped)
	}
}

func TestTableEmpty(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Table(nil, nil)

	if tp.Content() != "" {
		t.Fatalf("expected empty output, got %q", tp.Content())
	}
}

func TestTableWithoutHeaders(t *testing.T) {
	tp := Fake(t, 80, 24)

	defer tp.Cleanup()

	Table(nil, [][]string{
		{"Alice", "30"},
		{"Bob", "25"},
	})
	stripped := tp.StrippedContent()

	if !containsAll(stripped, "Alice", "Bob") {
		t.Fatalf("expected rows, got:\n%s", stripped)
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !contains(s, sub) {
			return false
		}
	}

	return true
}

func contains(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && indexOf(s, sub) >= 0
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}

	return -1
}
