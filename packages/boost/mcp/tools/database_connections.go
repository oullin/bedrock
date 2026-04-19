package tools

// DatabaseConnections returns the configured database connections.
// Mirrors Laravel\Boost\Mcp\Tools\DatabaseConnections.
// Tagged IsReadOnly.
type DatabaseConnections struct {
	// Connections is a map of connection-name → DSN/config. Injected by the caller.
	Connections    map[string]string
	DefaultConnect string
}

func (t *DatabaseConnections) Name() string     { return "database_connections" }
func (t *DatabaseConnections) IsReadOnly() bool { return true }

func (t *DatabaseConnections) Description() string {
	return "List all configured database connections and the default connection name."
}

func (t *DatabaseConnections) Schema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{},
	}
}

// Handle returns all connections and the default.
func (t *DatabaseConnections) Handle(_ McpRequest) (McpResponse, error) {
	names := make([]string, 0, len(t.Connections))

	for name := range t.Connections {
		names = append(names, name)
	}

	return OkResponse(map[string]any{
		"default":     t.DefaultConnect,
		"connections": names,
	}), nil
}
