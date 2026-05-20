# Configuration

<!-- upstream-docs: configuration.md#introduction -->
<!-- upstream-docs: configuration.md#environment-configuration -->
<!-- upstream-docs: configuration.md#accessing-configuration-values -->

Bedrock has two layers of configuration. The first is the **bootstrap
options** you pass when building the application — typed, package-shaped,
and read by service providers at registration time. The second is the
**`config.Repository`** — a runtime, dot-notation key-value store that
works the same way as the upstream `config(...)`.

You will use both. The bootstrap options decide which providers exist and
how they're constructed; the repository holds the values your application
code reads at runtime.

## Bootstrap Options

Every Bedrock application defines an `Options` struct that's the union of
every config field every standard provider needs. The demo's is a good
example to copy:

```go
// services/demo/api/bootstrap.go:38
type Options struct {
    BasePath                 string
    Env                      string
    AppKey                   string
    DatabaseURL              string
    StoragePath              string
    RunMigrations            *bool
    Seed                     *bool
    AuthDefaultGuard         string
    CacheDefaultDriver       string
    QueueDefaultConnection   string
    SessionName              string
    HashDefaultDriver        hashing.Driver
    EncryptionKey            []byte
    EncryptionCipher         encryption.Cipher
    CookieDefaults           cookie.Options
    LogConfig                log.LogProviderConfig
    ConcurrencyDefaultDriver string
    TranslationLoader        translation.Loader
    TranslationLocale        string
    AIDefaultProvider        string
    AIConfigs                map[string]map[string]any
}
```

Two patterns make this work well:

### 1. Defaults live alongside the struct

```go
// services/demo/api/bootstrap.go:63
func (o Options) withDefaults() Options {
    if o.Env == "" {
        o.Env = "local"
    }
    if o.CacheDefaultDriver == "" {
        o.CacheDefaultDriver = "array"
    }
    if o.QueueDefaultConnection == "" {
        o.QueueDefaultConnection = "sync"
    }
    /* ... */
    return o
}
```

This is the place to encode "what does an unconfigured Bedrock app do?"
The demo defaults to in-memory cache, sync queue, sqlite database, bcrypt
hashing — the cheapest, batteries-included answer for each driver.

### 2. Options are passed to provider constructors verbatim

```go
// services/demo/api/bootstrap.go:144
cache.NewCacheServiceProvider(application.Container, o.CacheDefaultDriver),
queue.NewQueueServiceProvider(application.Container, o.QueueDefaultConnection),
log.NewLogServiceProvider(application.Container, o.LogConfig),
```

The provider stack is just a function of your `Options`. To run with a
different config — production-strength Redis cache, real SMTP, S3
filesystem — you build a different `Options` value and the same provider
stack picks up the new defaults.

## Loading Options From the Environment

Whether the values come from environment variables, a YAML file, a CLI
flag set, or a Vault lookup is up to your `main`. The `Options` struct is
just a Go value. A common shape:

```go
func loadOptions() (api.Options, error) {
    return api.Options{
        Env:                    os.Getenv("APP_ENV"),
        AppKey:                 os.Getenv("APP_KEY"),
        DatabaseURL:            os.Getenv("DATABASE_URL"),
        CacheDefaultDriver:     getenvDefault("CACHE_DRIVER", "array"),
        QueueDefaultConnection: getenvDefault("QUEUE_CONNECTION", "sync"),
        EncryptionKey:          decodeKey(os.Getenv("APP_KEY")),
    }, nil
}

func main() {
    opts, err := loadOptions()
    if err != nil {
        log.Fatal(err)
    }

    if err := api.Run(ctx, opts); err != nil {
        log.Fatal(err)
    }
}
```

Bedrock doesn't prescribe a config-file format. Pick whatever fits — `os.Getenv`
for 12-factor deployments, [`viper`](https://pkg.go.dev/github.com/spf13/viper)
for layered TOML/YAML, or your own loader.

## The `config.Repository` — Runtime Lookups

Some providers read their per-feature configuration from a
`*config.Repository` rather than from individual constructor arguments —
the log manager is the prominent example
([`packages/log/manager.go:47`](https://github.com/gocanto/bedrock/blob/main/packages/log/manager.go#L47)).
The repository gives you:

- Dot-notation keys (`logging.channels.stack.driver`)
- Type-safe getters (`GetString`, `GetInt`, `GetDuration`, …)
- Layered defaults via `viper`

Build one from a plain map and pass it to whatever provider needs it:

```go
import "github.com/bedrock/packages/config"

repo := config.New(map[string]any{
    "logging": map[string]any{
        "default": "stack",
        "channels": map[string]any{
            "stack":  map[string]any{"driver": "stack", "channels": []string{"stderr"}},
            "stderr": map[string]any{"driver": "stderr", "level": "info"},
        },
    },
})

value := repo.Get("logging.default")            // "stack"
level := repo.GetString("logging.channels.stderr.level") // "info"
```

See [`packages/config/repository.go:34`](https://github.com/gocanto/bedrock/blob/main/packages/config/repository.go#L34)
for the constructor and [`Has`/`Get`/`Set`](https://github.com/gocanto/bedrock/blob/main/packages/config/repository.go#L64)
for the lookup API.

### Pre-shaped configs from typed structs

The `LogServiceProvider` ships a tiny `LogProviderConfig` struct that
collapses to a `*config.Repository` on registration. This is the bridge
between the typed-options world and the runtime-repository world:

```go
// packages/log/log_service_provider.go:12
type LogProviderConfig struct {
    Default  string
    Channels map[string]map[string]any
}

func (c LogProviderConfig) toRepository() *config.Repository { /* ... */ }
```

If you write your own provider that needs a repository internally, this is
the pattern to copy: take a typed config in the constructor, convert it to
a repository in `Register()`.

## Storing Per-App Values in the Container

When you have a value that doesn't belong to any package — your
application's name, a feature flag, a request signing key — bind it as a
container instance:

```go
// services/demo/api/server.go:26
application.Container.Instance("demo.config.app", democonfig.DefaultApp(o.Env, o.AppKey))
```

Then read it back wherever you need it:

```go
// services/demo/api/routes/web.go:20
if raw, err := application.Make("demo.config.app"); err == nil {
    if cfg, ok := raw.(democonfig.App); ok {
        // use cfg
    }
}
```

This is how the demo passes its app-level config through to handlers
without a global. See
[Service Container](/architecture/service-container#instancename-value--register-a-pre-built-object).

## Choosing Between Layers

| You want to                                                | Use                                                               |
| ---------------------------------------------------------- | ----------------------------------------------------------------- |
| Decide which providers exist                               | Conditional `append` to your `StandardProviders` list             |
| Pick the default driver for a manager                      | A field on your `Options` struct                                  |
| Configure a manager that reads dot-notation keys (logging) | Build a `*config.Repository` and pass it via `LogProviderConfig`  |
| Share an app-specific value across handlers                | `Container.Instance("app.something", value)`                      |
| Look up environment-derived values at request time         | Read from `Container.Instance` or a service that wraps env access |

## See Also

- [Application Bootstrap](/architecture/application) — where the `Options`
  struct is consumed.
- [Service Providers](/architecture/service-providers) — what reads
  `Options` and what reads `*config.Repository`.
- [Drivers](/architecture/drivers) — the link between
  `o.CacheDefaultDriver = "redis"` and `cache.Manager.Build(...)`.
- [`packages/config/repository.go`](https://github.com/gocanto/bedrock/blob/main/packages/config/repository.go)
  — the full repository API.
