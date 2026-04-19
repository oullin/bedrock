// Package demo is a standalone example of how to compose a Bedrock
// application entirely from packages/ primitives.
package demo

import (
	"github.com/bedrock/packages/ai"
	"github.com/bedrock/packages/auth"
	"github.com/bedrock/packages/bus"
	"github.com/bedrock/packages/cache"
	"github.com/bedrock/packages/concurrency"
	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/contracts/provider"
	"github.com/bedrock/packages/cookie"
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
)

// Options configures the standard provider stack used by the demo application.
type Options struct {
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
	application := container.NewApplication()

	application.RegisterMany(StandardProviders(application, opts...))
	application.Boot()

	return application
}

func defaultCookieOptions(o cookie.Options) cookie.Options {
	if o.Path == "" && o.HTTPOnly == nil && o.SameSite == 0 {
		return cookie.DefaultOptions()
	}

	return o
}
