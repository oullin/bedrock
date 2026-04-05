package routes

import (
	"context"

	"github.com/gollin/packages/illuminate/console"
	"github.com/gollin/packages/illuminate/foundation"
)

// RegisterConsole installs the demo console routes.
func RegisterConsole(kernel *console.Kernel) {
	kernel.Command("inspire", "Display an inspiring quote", func(_ context.Context, invocation *console.Invocation) error {
		invocation.Comment(foundation.InspiringQuote())
		return nil
	})
}
