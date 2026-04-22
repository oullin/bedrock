package support

import "testing"

// Exact inventory markers covered by the executable tests in this file:
// SupportHtmlStringTest::testToHtml
// SupportHtmlStringTest::testToString
// SupportHtmlStringTest::testIsEmpty
// SupportHtmlStringTest::testIsNotEmpty

func TestHtmlString(t *testing.T) {
	t.Parallel()

	html := NewHtmlString("<strong>Upstream</strong>")

	if got := html.ToHTML(); got != "<strong>Upstream</strong>" {
		t.Fatalf("ToHTML = %q", got)
	}

	if got := html.String(); got != "<strong>Upstream</strong>" {
		t.Fatalf("String = %q", got)
	}

	if html.IsEmpty() {
		t.Fatal("expected non-empty HtmlString")
	}

	if !html.IsNotEmpty() {
		t.Fatal("expected IsNotEmpty")
	}

	if !NewHtmlString("").IsEmpty() {
		t.Fatal("expected empty HtmlString")
	}
}
