package mcp_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/mcp"
)

// Port of Upstream\Mcp\Tests\CompletionTest

// completablePrompt is a Prompt that also implements Completable.
type completablePrompt struct {
	mcp.Prompt
}

func (p *completablePrompt) Complete(_ context.Context, argument, value string) *mcp.CompletionResult {
	switch argument {
	case "language":
		return mcp.MatchCompletion([]string{"Go", "Rust", "Python", "Ruby", "JavaScript"}, value)
	case "style":
		return mcp.EnumCompletion([]string{"formal", "casual", "technical"})
	}
	return mcp.EmptyCompletion()
}

// completableResource is a Resource that also implements Completable.
type completableResource struct {
	mcp.Resource
}

func (r *completableResource) Complete(_ context.Context, argument, value string) *mcp.CompletionResult {
	if argument == "lang" {
		return mcp.MatchCompletion([]string{"Go", "Goat", "Gorilla"}, value)
	}
	return mcp.EmptyCompletion()
}

func TestCompletionCompleteForPromptReturnsMatchingValues(t *testing.T) {
	t.Parallel()

	base := mcp.NewPrompt("code-review", "Review code", []*mcp.Argument{
		mcp.NewArgument("language", "Programming language", true),
	}, func(_ context.Context, _ *mcp.Request) ([]*mcp.Message, error) {
		return nil, nil
	})
	prompt := &completablePrompt{Prompt: base}

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddPrompt(prompt)

	result := srv.Test(t).Complete("ref/prompt", "code-review", "language", "Go")
	result.AssertOK().AssertHasCompletions()
	result.AssertSee("Go")
}

func TestCompletionCompleteForPromptPrefixFilters(t *testing.T) {
	t.Parallel()

	base := mcp.NewPrompt("p", "P", nil, func(_ context.Context, _ *mcp.Request) ([]*mcp.Message, error) {
		return nil, nil
	})
	prompt := &completablePrompt{Prompt: base}

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddPrompt(prompt)

	result := srv.Test(t).Complete("ref/prompt", "p", "language", "Py")
	result.AssertOK().AssertCompletionValues("Python")
}

func TestCompletionCompleteForResourceReturnsValues(t *testing.T) {
	t.Parallel()

	base := mcp.NewResource("langs", "Languages", "file://langs", "text/plain",
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("ok"), nil
		})
	resource := &completableResource{Resource: base}

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(resource)

	result := srv.Test(t).Complete("ref/resource", "file://langs", "lang", "Go")
	result.AssertOK().AssertHasCompletions()
}

func TestCompletionCompleteNonCompletablePrimitiveReturnsEmpty(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddPrompt(mcp.NewPrompt("plain", "Plain prompt", nil,
		func(_ context.Context, _ *mcp.Request) ([]*mcp.Message, error) {
			return nil, nil
		}))

	result := srv.Test(t).Complete("ref/prompt", "plain", "anything", "val")
	result.AssertOK().AssertCompletionCount(0)
}

func TestCompletionCompleteUnknownRefReturnsEmpty(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	result := srv.Test(t).Complete("ref/prompt", "nonexistent", "arg", "val")
	result.AssertOK().AssertCompletionCount(0)
}

func TestMatchCompletionCaseInsensitive(t *testing.T) {
	t.Parallel()

	result := mcp.MatchCompletion([]string{"Go", "Goat", "Python"}, "go")
	if len(result.Values) != 2 {
		t.Fatalf("expected 2 matches for 'go' (case-insensitive), got %d: %v", len(result.Values), result.Values)
	}
}

func TestEnumCompletionReturnsAll(t *testing.T) {
	t.Parallel()

	values := []string{"a", "b", "c"}
	result := mcp.EnumCompletion(values)
	if len(result.Values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(result.Values))
	}
	if result.HasMore {
		t.Fatal("expected HasMore=false for small enum")
	}
}

func TestEnumCompletionTruncatesAt100(t *testing.T) {
	t.Parallel()

	values := make([]string, 150)
	for i := range values {
		values[i] = "item"
	}
	result := mcp.EnumCompletion(values)
	if len(result.Values) != 100 {
		t.Fatalf("expected 100 values (truncated), got %d", len(result.Values))
	}
	if !result.HasMore {
		t.Fatal("expected HasMore=true when list exceeds 100")
	}
	if result.Total != 150 {
		t.Fatalf("expected Total=150, got %d", result.Total)
	}
}

func TestEmptyCompletionReturnsZeroValues(t *testing.T) {
	t.Parallel()

	result := mcp.EmptyCompletion()
	if len(result.Values) != 0 {
		t.Fatalf("expected 0 values, got %d", len(result.Values))
	}
	if result.HasMore {
		t.Fatal("expected HasMore=false")
	}
}
