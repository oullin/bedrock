package scout

// Config holds the Scout configuration. It mirrors the scout.php config
// file from Scout.
type Config struct {
	// Driver is the default search engine driver name.
	Driver string `json:"driver"`
	// Prefix is prepended to all search index names.
	Prefix string `json:"prefix"`
	// Queue is the queue connection name for async indexing. Empty disables queueing.
	Queue string `json:"queue"`
	// AfterCommit defers index syncing until the active database transaction commits.
	AfterCommit bool `json:"after_commit"`
	// Chunk configures batch sizes for bulk operations.
	Chunk ChunkConfig `json:"chunk"`
	// SoftDelete keeps soft-deleted models in the search index when true.
	SoftDelete bool `json:"soft_delete"`
	// Identify sends the authenticated user ID with search requests (Algolia).
	Identify bool `json:"identify"`
	// Engines holds per-engine configuration keyed by driver name.
	Engines map[string]map[string]any `json:"engines"`
}

// ChunkConfig controls batch sizes for bulk index operations.
type ChunkConfig struct {
	// Searchable is the chunk size for making models searchable. Default 500.
	Searchable int `json:"searchable"`
	// Unsearchable is the chunk size for removing models from search. Default 500.
	Unsearchable int `json:"unsearchable"`
}

// DefaultConfig returns a Config with sensible defaults matching Scout.
func DefaultConfig() Config {
	return Config{
		Driver:      "database",
		AfterCommit: false,
		Chunk: ChunkConfig{
			Searchable:   500,
			Unsearchable: 500,
		},
		SoftDelete: false,
		Engines:    make(map[string]map[string]any),
	}
}

// GetChunkSearchable returns the searchable chunk size, defaulting to 500.
func (c Config) GetChunkSearchable() int {
	if c.Chunk.Searchable <= 0 {
		return 500
	}

	return c.Chunk.Searchable
}

// GetChunkUnsearchable returns the unsearchable chunk size, defaulting to 500.
func (c Config) GetChunkUnsearchable() int {
	if c.Chunk.Unsearchable <= 0 {
		return 500
	}

	return c.Chunk.Unsearchable
}

// EngineConfig returns the configuration map for the given engine driver.
func (c Config) EngineConfig(driver string) map[string]any {
	if c.Engines == nil {
		return nil
	}

	return c.Engines[driver]
}
