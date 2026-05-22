package console

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed stubs/middleware.stub
var middlewareStub string

// MiddlewareMakeCommand renders the middleware stub for a new file.
//
// Ref: @bedrock/code-0291
type MiddlewareMakeCommand struct{}

// Render returns the file contents for a middleware named name in the given
// Go package.
func (MiddlewareMakeCommand) Render(pkg, name string) string {
	out := strings.ReplaceAll(middlewareStub, "{{ namespace }}", pkg)
	out = strings.ReplaceAll(out, "{{ name }}", name)

	return out
}

// Path returns the conventional file path for a middleware.
func (MiddlewareMakeCommand) Path(name string) string {
	return fmt.Sprintf("middleware/%s.go", strings.ToLower(name))
}
