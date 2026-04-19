// Package app wires every standard bedrock service provider into a
// single Application. It is the quickstart for new applications:
//
//	application := app.Default()
//	app.SetApp(application)
//	// application.Make("cache"), facades, etc. all work
//
// Production apps are expected to compose providers manually so they can
// substitute drivers, choose specific defaults, and skip components they
// don't use. Default is a sensible starting point, not a binding contract.
package app

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

// Options configures the standard provider stack. Every field has a
// sensible default — pass an empty Options{} to accept all of them.
type Options struct {
	// AuthDefaultGuard is the name passed to AuthServiceProvider. Default: "web".
	AuthDefaultGuard string

	// CacheDefaultDriver is the name passed to CacheServiceProvider. Default: "array".
	CacheDefaultDriver string

	// QueueDefaultConnection is the name passed to QueueServiceProvider. Default: "sync".
	QueueDefaultConnection string

	// SessionName is the cookie name passed to SessionServiceProvider. Default: "bedrock_session".
	SessionName string

	// HashDefaultDriver is the default hashing algorithm. Default: hashing.DriverBcrypt.
	HashDefaultDriver hashing.Driver

	// EncryptionKey is required if you intend to resolve "encrypter".
	// 16/24/32 bytes for AES-128/192/256-GCM. If nil, the encryption
	// provider is skipped entirely (resolving "encrypter" returns ErrNotBound).
	EncryptionKey []byte

	// EncryptionCipher is the AES variant. Default: encryption.CipherAES256GCM.
	EncryptionCipher encryption.Cipher

	// CookieDefaults overrides the default cookie options. Default: cookie.DefaultOptions().
	CookieDefaults cookie.Options

	// LogConfig is forwarded to LogServiceProvider. Default: empty (single channel, debug level).
	LogConfig log.LogProviderConfig

	// ConcurrencyDefaultDriver is the name passed to ConcurrencyServiceProvider.
	// Default: "goroutine".
	ConcurrencyDefaultDriver string

	// TranslationLoader is the loader used by the translator. If nil, the
	// translation provider is skipped.
	TranslationLoader translation.Loader

	// TranslationLocale is the default locale. Default: "en".
	TranslationLocale string

	// AIDefaultProvider is the lab name of the default AI provider (e.g. "openai").
	// When empty, the AI service provider is not registered.
	AIDefaultProvider string

	// AIConfigs maps AI provider lab names to their configuration (api_key, etc.).
	// Only meaningful when AIDefaultProvider is non-empty.
	AIConfigs map[string]map[string]any
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

// Default returns an Application with every standard provider registered
// and booted, in the correct dependency order.
//
// The order matters because providers that resolve other providers in their
// factory closures (bus → queue, notifications → bus + events) need their
// prerequisites bound first. Lazy resolution makes order forgiving but not
// arbitrary — leaf providers come first.
func Default(opts ...Options) *container.Application {
	o := Options{}

	if len(opts) > 0 {
		o = opts[0]
	}

	o = o.withDefaults()
	application := container.NewApplication()

	providers := []provider.ServiceProvider{
		// Layer 0: foundational, no cross-provider deps
		events.NewEventsServiceProvider(application.Container),
		hashing.NewHashingServiceProviderWithDefaults(application.Container, o.HashDefaultDriver),
		filesystem.NewFilesystemServiceProvider(application.Container),
		cookie.NewCookieServiceProvider(application.Container, defaultCookieOptions(o.CookieDefaults)),
		validation.NewValidationServiceProvider(application.Container),
		concurrency.NewConcurrencyServiceProvider(application.Container, o.ConcurrencyDefaultDriver),

		// Layer 1: depend only on layer 0
		cache.NewCacheServiceProvider(application.Container, o.CacheDefaultDriver),
		session.NewSessionServiceProvider(application.Container, o.SessionName),
		queue.NewQueueServiceProvider(application.Container, o.QueueDefaultConnection),
		log.NewLogServiceProvider(application.Container, o.LogConfig),
		auth.NewAuthServiceProvider(application.Container, o.AuthDefaultGuard),

		// Layer 2: depend on layer 1 collaborators (resolved lazily)
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

	application.RegisterMany(providers)
	application.Boot()

	return application
}

// defaultCookieOptions returns the user-supplied options if set, otherwise
// cookie.DefaultOptions().
func defaultCookieOptions(o cookie.Options) cookie.Options {
	// "Set" detection is awkward because Options has zero-value bool pointers.
	// Use Path == "" as the signal for "caller did not supply options".
	if o.Path == "" && o.HTTPOnly == nil && o.SameSite == 0 {
		return cookie.DefaultOptions()
	}

	return o
}
