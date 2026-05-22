package credentials

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awscreds "github.com/aws/aws-sdk-go-v2/credentials"
)

// ErrUnknownProvider is returned by Registry.Resolve when the caller
// asks for a provider name that has not been registered.

// Provider aliases aws.CredentialsProvider so callers can depend on
// this package without importing the AWS SDK directly when they only
// need the registry surface. The alias means a *Provider returned here
// is interchangeable with anything the SDK's config package accepts.
type Provider = aws.CredentialsProvider

// Factory constructs (or resolves) a Provider lazily. Resolution is
// deferred until the first Resolve call so that processes that boot
// long before they ever touch a queue do not pay the cost of, say,
// reaching the EC2 metadata service or loading a shared config file.
type Factory func(ctx context.Context) (Provider, error)

// Registry maps human-readable provider names to credential factories.
// It is safe for concurrent use. The intended lifecycle is:
//
//   - construct one Registry per process
//   - Register every named provider at bootstrap
//   - look up by name when a connection is resolved
//
// providers introduced in 13.8.0.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
	cache     map[string]Provider
}

var ErrUnknownProvider = errors.New("credentials: unknown provider")

// NewRegistry constructs an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]Factory),
		cache:     make(map[string]Provider),
	}
}

// Register binds factory to name. A subsequent call with the same name
// replaces the previous factory and invalidates any cached resolution
// so the next Resolve picks up the new wiring.
func (r *Registry) Register(name string, factory Factory) *Registry {
	if name == "" || factory == nil {
		return r
	}

	r.mu.Lock()

	defer r.mu.Unlock()

	r.factories[name] = factory

	delete(r.cache, name)

	return r
}

// Names returns the registered provider names in sorted order. The
// ordering is deterministic so callers can render them stably in CLI
// help output or logs.
func (r *Registry) Names() []string {
	r.mu.RLock()

	defer r.mu.RUnlock()

	out := make([]string, 0, len(r.factories))

	for name := range r.factories {
		out = append(out, name)
	}

	sort.Strings(out)

	return out
}

// Resolve returns the Provider for name, invoking the registered
// factory the first time and caching the result for subsequent calls.
// Returns ErrUnknownProvider wrapped with the requested name when no
// factory is registered.
func (r *Registry) Resolve(ctx context.Context, name string) (Provider, error) {
	r.mu.RLock()

	if cached, ok := r.cache[name]; ok {
		r.mu.RUnlock()

		return cached, nil
	}

	factory, ok := r.factories[name]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownProvider, name)
	}

	provider, err := factory(ctx)

	if err != nil {
		return nil, fmt.Errorf("credentials: resolve %q: %w", name, err)
	}

	r.mu.Lock()

	r.cache[name] = provider

	r.mu.Unlock()

	return provider, nil
}

// --- Built-in factories ----------------------------------------------

// Profile returns a Factory backed by a named AWS shared-config
// profile. Equivalent to running with AWS_PROFILE=<profileName>.
func Profile(profileName string) Factory {
	return func(ctx context.Context) (Provider, error) {
		cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithSharedConfigProfile(profileName))

		if err != nil {
			return nil, err
		}

		return cfg.Credentials, nil
	}
}

// Static returns a Factory that always yields the literal access key,
// secret, and (optional) session token. Useful for tests, CI runners,
// and one-off scripts. Avoid in production — rotating static
// credentials is a known footgun.
func Static(accessKeyID, secretAccessKey, sessionToken string) Factory {
	return func(_ context.Context) (Provider, error) {
		return awscreds.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, sessionToken), nil
	}
}

// EC2InstanceRole returns a Factory backed by the IAM Role attached to
// the EC2 instance running the process — i.e. credentials fetched from
// the instance metadata service. Resolves to the SDK's default config,
// which already prefers IMDS when no other source is present.
func EC2InstanceRole() Factory {
	return func(ctx context.Context) (Provider, error) {
		cfg, err := awsconfig.LoadDefaultConfig(ctx)

		if err != nil {
			return nil, err
		}

		return cfg.Credentials, nil
	}
}

// SSO returns a Factory backed by an AWS SSO profile. The profile must
// already be configured (via `aws configure sso`) and the caller must
// have run `aws sso login --profile <profileName>` recently enough
// that the cached SSO token is still valid.
func SSO(profileName string) Factory {
	return Profile(profileName)
}

// Chain returns a Factory that tries each child Factory in order and
// returns the first that successfully resolves. Failures from earlier
// factories are not surfaced; only the last error is returned if every
// child fails. Use this to express a preference order ("prefer SSO,
// fall back to the static key").
func Chain(factories ...Factory) Factory {
	return func(ctx context.Context) (Provider, error) {
		var lastErr error

		for _, f := range factories {
			if f == nil {
				continue
			}

			provider, err := f(ctx)

			if err == nil {
				return provider, nil
			}

			lastErr = err
		}

		if lastErr != nil {
			return nil, lastErr
		}

		return nil, errors.New("credentials: empty factory chain")
	}
}
