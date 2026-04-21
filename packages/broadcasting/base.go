package broadcasting

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	contractsauth "github.com/bedrock/packages/contracts/auth"
)

// BaseBroadcaster contains framework-level channel registration and
// authorization behavior shared by concrete broadcaster backends.
type BaseBroadcaster struct {
	mu                        sync.RWMutex
	patterns                  []string
	channels                  map[string]any
	channelOptions            map[string]ChannelOptions
	bindings                  map[string]BindingFunc
	authenticatedUserCallback func(AuthRequest) any
}

// NewBaseBroadcaster creates a BaseBroadcaster.
func NewBaseBroadcaster() *BaseBroadcaster {
	return &BaseBroadcaster{
		channels:       make(map[string]any),
		channelOptions: make(map[string]ChannelOptions),
		bindings:       make(map[string]BindingFunc),
	}
}

// Channel registers an authorization handler for a channel pattern.
func (b *BaseBroadcaster) Channel(pattern string, handler any, options ...ChannelOptions) *BaseBroadcaster {
	b.mu.Lock()

	defer b.mu.Unlock()

	if _, exists := b.channels[pattern]; !exists {
		b.patterns = append(b.patterns, pattern)
	}

	b.channels[pattern] = handler

	if len(options) > 0 {
		b.channelOptions[pattern] = options[0]
	} else {
		b.channelOptions[pattern] = ChannelOptions{}
	}

	return b
}

// Bind registers an explicit resolver for a named channel parameter.
func (b *BaseBroadcaster) Bind(key string, binding BindingFunc) *BaseBroadcaster {
	b.mu.Lock()

	defer b.mu.Unlock()

	b.bindings[key] = binding

	return b
}

// ResolveAuthenticatedUserUsing registers the user authentication callback used
// by Pusher-compatible user authentication.
func (b *BaseBroadcaster) ResolveAuthenticatedUserUsing(callback func(AuthRequest) any) *BaseBroadcaster {
	b.mu.Lock()

	defer b.mu.Unlock()

	b.authenticatedUserCallback = callback

	return b
}

// ResolveAuthenticatedUser resolves the authenticated user payload.
func (b *BaseBroadcaster) ResolveAuthenticatedUser(request AuthRequest) any {
	b.mu.RLock()
	callback := b.authenticatedUserCallback
	b.mu.RUnlock()

	if callback == nil {
		return nil
	}

	return callback(request)
}

// VerifyUserCanAccessChannel invokes the registered channel handler for a
// normalized channel name and returns the handler result.
func (b *BaseBroadcaster) VerifyUserCanAccessChannel(request AuthRequest, channel string) (any, error) {
	pattern, handler, ok := b.matchingChannel(channel)

	if !ok {
		return nil, ErrAccessDenied
	}

	params, err := b.ExtractAuthParameters(pattern, channel)

	if err != nil {
		return nil, err
	}

	call, err := normalizeChannelHandler(handler)

	if err != nil {
		return nil, err
	}

	result, err := call(b.RetrieveUser(request, channel), params...)

	if err != nil {
		return nil, err
	}

	if !truthy(result) {
		return nil, ErrAccessDenied
	}

	return result, nil
}

// ExtractAuthParameters resolves named parameters from a pattern/channel pair.
func (b *BaseBroadcaster) ExtractAuthParameters(pattern, channel string) ([]any, error) {
	keys, values, ok := extractChannelKeys(pattern, channel)

	if !ok {
		return nil, ErrAccessDenied
	}

	params := make([]any, 0, len(keys))

	b.mu.RLock()

	defer b.mu.RUnlock()

	for i, key := range keys {
		value := values[i]

		if binding, exists := b.bindings[key]; exists {
			resolved, err := binding(value)

			if err != nil {
				return nil, fmt.Errorf("%w: %s", ErrAccessDenied, err)
			}

			params = append(params, resolved)

			continue
		}

		params = append(params, value)
	}

	return params, nil
}

// RetrieveChannelOptions returns options for the first matching channel pattern.
func (b *BaseBroadcaster) RetrieveChannelOptions(channel string) ChannelOptions {
	b.mu.RLock()

	defer b.mu.RUnlock()

	for _, pattern := range b.patterns {
		if channelNameMatchesPattern(channel, pattern) {
			return cloneOptions(b.channelOptions[pattern])
		}
	}

	return ChannelOptions{}
}

// RetrieveUser returns the authenticated user for a channel, respecting guard
// options and guard order.
func (b *BaseBroadcaster) RetrieveUser(request AuthRequest, channel string) contractsauth.Authenticatable {
	if request.UserResolver == nil {
		return nil
	}

	options := b.RetrieveChannelOptions(channel)

	if len(options.Guards) == 0 {
		return request.UserResolver.User("")
	}

	for _, guard := range options.Guards {
		if user := request.UserResolver.User(guard); user != nil {
			return user
		}
	}

	return nil
}

// ChannelNameMatchesPattern reports whether channel matches pattern.
func (b *BaseBroadcaster) ChannelNameMatchesPattern(channel, pattern string) bool {
	return channelNameMatchesPattern(channel, pattern)
}

func (b *BaseBroadcaster) matchingChannel(channel string) (string, any, bool) {
	b.mu.RLock()

	defer b.mu.RUnlock()

	for _, pattern := range b.patterns {
		if channelNameMatchesPattern(channel, pattern) {
			return pattern, b.channels[pattern], true
		}
	}

	return "", nil, false
}

func normalizeChannelHandler(handler any) (ChannelHandler, error) {
	switch h := handler.(type) {
	case ChannelHandler:
		return h, nil
	case func(contractsauth.Authenticatable, ...any) (any, error):
		return ChannelHandler(h), nil
	case func(contractsauth.Authenticatable, ...any) any:
		return func(user contractsauth.Authenticatable, params ...any) (any, error) {
			return h(user, params...), nil
		}, nil
	case func(contractsauth.Authenticatable, ...any) bool:
		return func(user contractsauth.Authenticatable, params ...any) (any, error) {
			return h(user, params...), nil
		}, nil
	case ChannelJoiner:
		return h.Join, nil
	default:
		return nil, ErrUnknownChannelHandler
	}
}

func truthy(value any) bool {
	if value == nil {
		return false
	}

	if b, ok := value.(bool); ok {
		return b
	}

	return true
}

func cloneOptions(options ChannelOptions) ChannelOptions {
	if len(options.Guards) == 0 {
		return ChannelOptions{}
	}

	return ChannelOptions{Guards: append([]string(nil), options.Guards...)}
}

func channelNameMatchesPattern(channel, pattern string) bool {
	_, _, ok := extractChannelKeys(pattern, channel)

	return ok
}

func extractChannelKeys(pattern, channel string) ([]string, []string, bool) {
	keys, expr := compileChannelPattern(pattern)
	matches := expr.FindStringSubmatch(channel)

	if matches == nil {
		return nil, nil, false
	}

	return keys, matches[1:], true
}

func compileChannelPattern(pattern string) ([]string, *regexp.Regexp) {
	keys := make([]string, 0)

	var out strings.Builder

	out.WriteString("^")

	for i := 0; i < len(pattern); {
		if pattern[i] == '{' {
			end := strings.IndexByte(pattern[i:], '}')

			if end >= 0 {
				key := pattern[i+1 : i+end]
				keys = append(keys, key)
				out.WriteString("([^\\.]+)")
				i += end + 1

				continue
			}
		}

		out.WriteString(regexp.QuoteMeta(pattern[i : i+1]))
		i++
	}

	out.WriteString("$")

	return keys, regexp.MustCompile(out.String())
}
