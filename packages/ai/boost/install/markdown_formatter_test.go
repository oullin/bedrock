package install_test

import (
	"testing"

	"github.com/bedrock/packages/ai/boost/install"
)

func TestMarkdownFormatterStripFrontmatter(t *testing.T) {
	t.Parallel()

	f := &install.MarkdownFormatter{}

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no frontmatter",
			input: "# Hello\n",
			want:  "# Hello\n",
		},
		{
			name:  "with frontmatter",
			input: "---\nname: test\n---\n# Hello\n",
			want:  "# Hello\n",
		},
		{
			name:  "empty frontmatter",
			input: "---\n---\n# Body\n",
			want:  "# Body\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := f.StripFrontmatter(tc.input)

			if got != tc.want {
				t.Errorf("StripFrontmatter() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMarkdownFormatterAddFrontmatter(t *testing.T) {
	t.Parallel()

	f := &install.MarkdownFormatter{}
	content := "# Body\n"
	data := map[string]string{"name": "myskill"}

	result := f.AddFrontmatter(content, data)

	if result == content {
		t.Error("AddFrontmatter should modify the content")
	}

	// Must start with ---
	if len(result) < 4 || result[:4] != "---\n" {
		t.Errorf("AddFrontmatter result should start with ---\\n, got %q", result[:min(10, len(result))])
	}
}

func TestMarkdownFormatterNormalizeHeadingsNoPromotion(t *testing.T) {
	t.Parallel()

	f := &install.MarkdownFormatter{}
	content := "# H1\n## H2\n"
	got := f.NormalizeHeadings(content)

	if got != content {
		t.Errorf("NormalizeHeadings should be identity when min level is 1, got %q", got)
	}
}

func TestMarkdownFormatterNormalizeHeadingsPromotes(t *testing.T) {
	t.Parallel()

	f := &install.MarkdownFormatter{}
	content := "## H2\n### H3\n"
	got := f.NormalizeHeadings(content)

	if got == content {
		t.Error("NormalizeHeadings should promote H2 → H1 when no H1 exists")
	}

	if len(got) < 2 || got[:2] != "# " {
		t.Errorf("expected first heading to be H1, got %q", got[:min(5, len(got))])
	}
}

func TestMarkdownFormatterTrim(t *testing.T) {
	t.Parallel()

	f := &install.MarkdownFormatter{}

	if got := f.Trim("  hello  "); got != "hello" {
		t.Errorf("Trim() = %q, want \"hello\"", got)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
