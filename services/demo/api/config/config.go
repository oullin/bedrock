package config

// App mirrors the small amount of config the skeleton demo needs at runtime.
type App struct {
	Name string
	Env  string
	Key  string
	URL  string
}

// DefaultApp returns Upstream-skeleton-like defaults for local development.
func DefaultApp(env, key string) App {
	if env == "" {
		env = "local"
	}

	return App{
		Name: "Bedrock",
		Env:  env,
		Key:  key,
		URL:  "http://localhost:8080",
	}
}
