package routes

import (
	"fmt"
	"net/http"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/lottery"
	"github.com/bedrock/packages/routing"
	"github.com/bedrock/packages/support"
	"github.com/bedrock/packages/validation"
	democonfig "github.com/bedrock/services/demo/api/config"
)

// RegisterWeb mounts the Upstream skeleton web and health routes.
func RegisterWeb(router *routing.Router, application *container.Application) {
	router.Get("/", func() any {
		app := democonfig.DefaultApp("", "")

		if raw, err := application.Make("demo.config.app"); err == nil {
			if cfg, ok := raw.(democonfig.App); ok {
				app = cfg
			}
		}

		return &routing.HTTPResponse{
			Body:    fmt.Sprintf("<!doctype html><html><head><title>%s</title></head><body><main><h1>%s</h1><p>Upstream skeleton port running on Bedrock.</p></main></body></html>\n", app.Name, app.Name),
			Status:  http.StatusOK,
			Headers: map[string][]string{"Content-Type": {"text/html; charset=utf-8"}},
		}
	})

	router.Get("/up", func() any {
		return &routing.HTTPResponse{
			Body:    "OK\n",
			Status:  http.StatusOK,
			Headers: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
		}
	})

	router.Get("/lottery", func() any {
		result := lottery.NewLottery(1, 1).
			Winner(func(...any) any { return "winner" }).
			Choose()

		return &routing.HTTPResponse{
			Body:    fmt.Sprintf("%v\n", result),
			Status:  http.StatusOK,
			Headers: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
		}
	})

	router.Get("/features/validation", func() any {
		input := map[string]any{
			"email": "taylor@example.com",
			"name":  "Taylor Otwell",
		}

		validated, err := validation.NewFactory().Validate(input, map[string]any{
			"email": "required|email",
			"name":  "required|string",
		}, nil, nil)

		if err != nil {
			return &routing.HTTPResponse{
				Body:    "validation failed\n",
				Status:  http.StatusUnprocessableEntity,
				Headers: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
			}
		}

		return &routing.HTTPResponse{
			Body:    fmt.Sprintf("Hello %s <%s>\n", support.ArrString(validated, "name"), support.ArrString(validated, "email")),
			Status:  http.StatusOK,
			Headers: map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}},
		}
	})
}
