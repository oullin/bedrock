package queue

import "github.com/bedrock/packages/container"

// QueueServiceProvider registers the queue manager into the container.
// It mirrors @bedrock\Queue\QueueServiceProvider.
type QueueServiceProvider struct {
	app               *container.Container
	defaultConnection string
}

// NewQueueServiceProvider constructs the provider.
// defaultConnection is the name of the default queue connection (e.g. "sync").
func NewQueueServiceProvider(app *container.Container, defaultConnection string) *QueueServiceProvider {
	return &QueueServiceProvider{app: app, defaultConnection: defaultConnection}
}

// Register binds the queue manager as a singleton under "queue".
func (p *QueueServiceProvider) Register() {
	p.app.Singleton("queue", func(_ *container.Container) (any, error) {
		m := NewManager()
		m.SetDefaultConnection(p.defaultConnection)

		return m, nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *QueueServiceProvider) Provides() []string {
	return []string{"queue"}
}
