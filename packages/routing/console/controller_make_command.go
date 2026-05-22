// Ref: @bedrock/code-0289
// The PHP package supplies "make:controller" and "make:middleware" CLI
// commands that scaffold class files from stubs. The Go port exposes the
// stub-rendering primitives so a bedrock CLI can wire the same UX without
// duplicating the templates.
package console

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed stubs/controller.stub
var controllerStub string

// ControllerMakeCommand renders the controller stub for a new file.
//
// Ref: @bedrock/code-0290
type ControllerMakeCommand struct{}

// Render returns the file contents for a controller named name in the given
// Go package (namespace).
func (ControllerMakeCommand) Render(pkg, name string) string {
	out := strings.ReplaceAll(controllerStub, "{{ namespace }}", pkg)
	out = strings.ReplaceAll(out, "{{ name }}", name)

	return out
}

// Path returns the conventional file path for a controller.
func (ControllerMakeCommand) Path(name string) string {
	return fmt.Sprintf("controllers/%s.go", strings.ToLower(name))
}
