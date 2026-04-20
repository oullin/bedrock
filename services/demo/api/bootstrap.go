// Package api is a standalone example of how to compose a Bedrock
// application entirely from packages/ primitives.
package api

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ai "github.com/bedrock/packages/ai/sdk"
	"github.com/bedrock/packages/auth"
	"github.com/bedrock/packages/bus"
	"github.com/bedrock/packages/cache"
	"github.com/bedrock/packages/concurrency"
	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/contracts/provider"
	"github.com/bedrock/packages/cookie"
	"github.com/bedrock/packages/database"
	"github.com/bedrock/packages/encryption"
	"github.com/bedrock/packages/events"
	"github.com/bedrock/packages/filesystem"
	"github.com/bedrock/packages/hashing"
	"github.com/bedrock/packages/log"
	"github.com/bedrock/packages/notifications"
	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/routing"
	"github.com/bedrock/packages/session"
	"github.com/bedrock/packages/translation"
	"github.com/bedrock/packages/validation"
	demomigrations "github.com/bedrock/services/demo/api/database/migrations"
	demoseeders "github.com/bedrock/services/demo/api/database/seeders"

	_ "modernc.org/sqlite"
)

// Options configures the standard provider stack used by the demo application.
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

func (o Options) withDefaults() Options {
	if o.BasePath == "" {
		if cwd, err := os.Getwd(); err == nil {
			o.BasePath = cwd
		}
	}

	if o.Env == "" {
		o.Env = "local"
	}

	if o.StoragePath == "" {
		o.StoragePath = filepath.Join(o.BasePath, "storage")
	}

	if o.DatabaseURL == "" {
		o.DatabaseURL = "sqlite:///" + filepath.ToSlash(filepath.Join(o.StoragePath, "database.sqlite"))
	}

	if o.RunMigrations == nil {
		yes := true
		o.RunMigrations = &yes
	}

	if o.Seed == nil {
		yes := true
		o.Seed = &yes
	}

	if o.AuthDefaultGuard == "" {
		o.AuthDefaultGuard = "web"
	}

	if o.CacheDefaultDriver == "" {
		o.CacheDefaultDriver = "array"
	}

	if o.QueueDefaultConnection == "" {
		o.QueueDefaultConnection = "sync"
	}

	if o.SessionName == "" {
		o.SessionName = "bedrock_session"
	}

	if o.HashDefaultDriver == "" {
		o.HashDefaultDriver = hashing.DriverBcrypt
	}

	if o.ConcurrencyDefaultDriver == "" {
		o.ConcurrencyDefaultDriver = "goroutine"
	}

	if o.TranslationLocale == "" {
		o.TranslationLocale = "en"
	}

	return o
}

func resolveOptions(opts []Options) Options {
	o := Options{}

	if len(opts) > 0 {
		o = opts[0]
	}

	return o.withDefaults()
}

// StandardProviders returns the canonical provider stack for the demo app.
func StandardProviders(application *container.Application, opts ...Options) []provider.ServiceProvider {
	o := resolveOptions(opts)

	providers := []provider.ServiceProvider{
		events.NewEventsServiceProvider(application.Container),
		hashing.NewHashingServiceProviderWithDefaults(application.Container, o.HashDefaultDriver),
		filesystem.NewFilesystemServiceProvider(application.Container),
		cookie.NewCookieServiceProvider(application.Container, defaultCookieOptions(o.CookieDefaults)),
		validation.NewValidationServiceProvider(application.Container),
		concurrency.NewConcurrencyServiceProvider(application.Container, o.ConcurrencyDefaultDriver),
		cache.NewCacheServiceProvider(application.Container, o.CacheDefaultDriver),
		session.NewSessionServiceProvider(application.Container, o.SessionName),
		queue.NewQueueServiceProvider(application.Container, o.QueueDefaultConnection),
		log.NewLogServiceProvider(application.Container, o.LogConfig),
		auth.NewAuthServiceProvider(application.Container, o.AuthDefaultGuard),
		bus.NewBusServiceProvider(application.Container),
		notifications.NewNotificationsServiceProvider(application.Container),
		database.NewDatabaseServiceProvider(application.Container, "sqlite"),
		routing.NewRoutingServiceProvider(application.Container),
	}

	if o.EncryptionKey != nil {
		cipher := o.EncryptionCipher

		if cipher == "" {
			cipher = encryption.AES256GCM
		}

		providers = append(providers, encryption.NewEncryptionServiceProvider(application.Container, o.EncryptionKey, cipher))
	}

	if o.TranslationLoader != nil {
		providers = append(providers, translation.NewTranslationServiceProvider(application.Container, o.TranslationLoader, o.TranslationLocale))
	}

	if o.AIDefaultProvider != "" {
		providers = append(providers, ai.NewAiServiceProvider(application.Container, o.AIDefaultProvider, o.AIConfigs))
	}

	return providers
}

// NewApplication builds and boots a demo application from the standard stack.
func NewApplication(opts ...Options) *container.Application {
	application, err := newApplication(opts...)

	if err != nil {
		panic(err)
	}

	return application
}

func newApplication(opts ...Options) (*container.Application, error) {
	o := resolveOptions(opts)
	application := container.NewApplication()

	application.RegisterMany(StandardProviders(application, o))
	application.Boot()

	if err := configureSkeleton(application, o); err != nil {
		return nil, err
	}

	return application, nil
}

func defaultCookieOptions(o cookie.Options) cookie.Options {
	if o.Path == "" && o.HTTPOnly == nil && o.SameSite == 0 {
		return cookie.DefaultOptions()
	}

	return o
}

func configureSkeleton(application *container.Application, o Options) error {
	if err := os.MkdirAll(o.StoragePath, 0o755); err != nil {
		return fmt.Errorf("demo: create storage path: %w", err)
	}

	dbPath, err := sqlitePath(o.DatabaseURL)

	if err != nil {
		return err
	}

	db, err := openSQLite(dbPath)

	if err != nil {
		return err
	}

	if o.RunMigrations != nil && *o.RunMigrations {
		if err := demomigrations.Run(db); err != nil {
			db.Close()

			return err
		}
	}

	if o.Seed != nil && *o.Seed {
		if err := demoseeders.Run(db); err != nil {
			db.Close()

			return err
		}
	}

	application.Container.Instance("demo.options", o)
	application.Container.Instance("demo.sql", db)

	if raw, err := application.Make("db"); err == nil {
		if manager, ok := raw.(*database.Manager); ok {
			manager.Extend("sqlite", func(config database.ConnectionConfig) (*database.Connection, error) {
				connDB, err := openSQLite(config.Database)

				if err != nil {
					return nil, err
				}

				conn := database.NewConnection(connDB, "sqlite", config.Database, "", map[string]any{
					"driver":   "sqlite",
					"database": config.Database,
				})
				conn.SetDriverName("sqlite")

				return conn, nil
			})
			manager.AddConnection("sqlite", database.ConnectionConfig{
				Driver:   "sqlite",
				Database: dbPath,
			})
		}
	}

	return nil
}

func sqlitePath(rawURL string) (string, error) {
	cfg, err := database.ParseDatabaseURL(rawURL)

	if err != nil {
		return "", fmt.Errorf("demo: parse database url: %w", err)
	}

	if cfg.Driver != "sqlite" {
		return "", fmt.Errorf("demo: unsupported database driver %q", cfg.Driver)
	}

	if cfg.Database == "" {
		return ":memory:", nil
	}

	return cfg.Database, nil
}

func openSQLite(path string) (*sql.DB, error) {
	dsn := path

	if strings.Contains(dsn, "?") {
		dsn += "&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
	} else {
		dsn += "?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
	}

	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, fmt.Errorf("demo: open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		db.Close()

		return nil, fmt.Errorf("demo: ping sqlite: %w", err)
	}

	return db, nil
}
