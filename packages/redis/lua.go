package redis

import (
	"context"
	"strconv"
)

// Eval runs a Lua script via EVAL. Parity with PhpRedisConnection::eval.
func (c *Connection) Eval(ctx context.Context, script string, keys []string, args ...any) (any, error) {
	argv := make([]any, 0, 2+len(keys)+len(args))
	argv = append(argv, script, strconv.Itoa(len(keys)))
	for _, k := range keys {
		argv = append(argv, k)
	}
	argv = append(argv, args...)
	return c.Command(ctx, "EVAL", argv...)
}

// EvalSha runs a cached Lua script by SHA.
func (c *Connection) EvalSha(ctx context.Context, sha string, keys []string, args ...any) (any, error) {
	argv := make([]any, 0, 2+len(keys)+len(args))
	argv = append(argv, sha, strconv.Itoa(len(keys)))
	for _, k := range keys {
		argv = append(argv, k)
	}
	argv = append(argv, args...)
	return c.Command(ctx, "EVALSHA", argv...)
}

// ScriptLoad uploads a script and returns its SHA.
func (c *Connection) ScriptLoad(ctx context.Context, script string) (string, error) {
	v, err := c.Command(ctx, "SCRIPT", "LOAD", script)
	if err != nil {
		return "", err
	}
	return toString(v)
}

// ScriptExists reports for each sha whether the script is cached.
func (c *Connection) ScriptExists(ctx context.Context, shas ...string) ([]bool, error) {
	args := append([]any{"EXISTS"}, toAnySlice(shas)...)
	v, err := c.Command(ctx, "SCRIPT", args...)
	if err != nil {
		return nil, err
	}
	s, err := toSlice(v)
	if err != nil {
		return nil, err
	}
	out := make([]bool, len(s))
	for i, e := range s {
		n, _ := toInt64(e)
		out[i] = n == 1
	}
	return out, nil
}
