package routes

import "github.com/bedrock/packages/anvil/routing"

// RegisterWeb installs the demo HTTP routes.
func RegisterWeb(router *routing.Router) {
	router.Get("/", func(ctx *routing.Context) error {
		return ctx.Render("welcome", map[string]any{
			"title": "Bedrock Demo",
		})
	})
}
