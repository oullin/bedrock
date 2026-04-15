package redis

import "github.com/bedrock/packages/container"

// RedisServiceProvider registers the Redis manager into the container.
// It mirrors Illuminate\Redis\RedisServiceProvider.
type RedisServiceProvider struct {
	app        *container.Container
	defaultConn string
	configs    map[string]ConnectionConfig
}

// NewRedisServiceProvider constructs the provider.
// defaultConn is the name of the default connection (e.g. "default").
// configs maps connection names to their ConnectionConfig.
func NewRedisServiceProvider(app *container.Container, defaultConn string, configs map[string]ConnectionConfig) *RedisServiceProvider {
	return &RedisServiceProvider{
		app:        app,
		defaultConn: defaultConn,
		configs:    configs,
	}
}

// Register binds the Redis manager as a singleton under "redis".
func (p *RedisServiceProvider) Register() {
	p.app.Singleton("redis", func(_ *container.Container) (any, error) {
		return NewManager(p.defaultConn, p.configs), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *RedisServiceProvider) Provides() []string {
	return []string{"redis"}
}
