package auth

import "fmt"

// App bundles the auth HTTP handlers, guards, and session helpers.
type App struct {
	container Container
	service   service
}

// NewApp builds the auth runtime with the provided host integrations.
func NewApp(container Container) (App, error) {
	if err := container.Validate(); err != nil {
		return App{}, fmt.Errorf("auth: %w", err)
	}

	svc, err := newService(container.DB)

	if err != nil {
		return App{}, fmt.Errorf("auth: %w", err)
	}

	return App{
		container: container,
		service:   svc,
	}, nil
}
