package database

import "github.com/bedrock/packages/container"

// DatabaseServiceProvider registers the database manager into the container.
// Ref: @bedrock/code-0205
type DatabaseServiceProvider struct {
	app               *container.Container
	defaultConnection string
}

// NewDatabaseServiceProvider constructs the provider.
// defaultConnection is the name of the default database connection (e.g. "mysql", "sqlite").
func NewDatabaseServiceProvider(app *container.Container, defaultConnection string) *DatabaseServiceProvider {
	return &DatabaseServiceProvider{app: app, defaultConnection: defaultConnection}
}

// Register binds the database manager as a singleton under "db".
func (p *DatabaseServiceProvider) Register() {
	p.app.Singleton("db", func(_ *container.Container) (any, error) {
		m := NewManager()
		m.SetDefaultConnection(p.defaultConnection)

		return m, nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *DatabaseServiceProvider) Provides() []string {
	return []string{"db"}
}
