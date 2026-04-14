package redis_test

import (
	"context"
	"testing"
)

func TestScriptLoadAndExists(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	sha, err := c.ScriptLoad(ctx, "return 1")
	if err != nil {
		t.Fatal(err)
	}
	if sha == "" {
		t.Fatal("empty sha")
	}
	exists, err := c.ScriptExists(ctx, sha)
	if err != nil {
		t.Fatal(err)
	}
	if len(exists) != 1 || !exists[0] {
		t.Fatalf("ScriptExists=%+v", exists)
	}
}
