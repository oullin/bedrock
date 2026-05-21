package session

import "github.com/bedrock/packages/container"

// SessionServiceProvider registers the session manager into the container.
// It mirrors @bedrock\Session\SessionServiceProvider.
type SessionServiceProvider struct {
	app  *container.Container
	name string // default session/cookie name
}

// NewSessionServiceProvider constructs the provider.
// name is the session cookie name (e.g. "bedrock_session").
func NewSessionServiceProvider(app *container.Container, name string) *SessionServiceProvider {
	return &SessionServiceProvider{app: app, name: name}
}

// Register binds the session manager as a singleton under "session".
func (p *SessionServiceProvider) Register() {
	p.app.Singleton("session", func(_ *container.Container) (any, error) {
		return NewManager(p.name), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *SessionServiceProvider) Provides() []string {
	return []string{"session"}
}
