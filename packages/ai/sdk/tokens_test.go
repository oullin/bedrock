package ai_test

import (
	"strings"
	"testing"

	ai "github.com/bedrock/packages/ai/sdk"
)

func TestEstimateTokens(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"a", 1},
		{"abcd", 1},
		{"abcde", 2},
		{strings.Repeat("x", 100), 25},
	}

	for _, c := range cases {
		if got := ai.EstimateTokens(c.in); got != c.want {
			t.Errorf("EstimateTokens(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}
