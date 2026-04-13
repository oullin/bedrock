package log

import (
	"fmt"
	"os"
	"sync"

	"github.com/bedrock/packages/config"
	cevents "github.com/bedrock/packages/contracts/events"
	clog "github.com/bedrock/packages/contracts/log"
)

// DriverFactory creates a Handler from a channel configuration.
type DriverFactory func(config ChannelConfig) (Handler, error)

// LogManager creates, caches, and delegates to channel-specific Logger instances.
type LogManager struct {
	mu             sync.RWMutex
	cfg            *config.Repository
	dispatcher     cevents.Dispatcher
	channels       map[string]*Logger
	customDrivers  map[string]DriverFactory
	sharedContext  map[string]any
	tapCallbacks   []func(*Logger)
	defaultChannel string
}

// ManagerOption configures a LogManager.
type ManagerOption func(*LogManager)

var _ clog.Logger = (*LogManager)(nil)

// WithEventDispatcher sets the event dispatcher on the manager.
func WithEventDispatcher(d cevents.Dispatcher) ManagerOption {
	return func(m *LogManager) {
		m.dispatcher = d
	}
}

// WithDefaultChannel sets the default channel name.
func WithDefaultChannel(name string) ManagerOption {
	return func(m *LogManager) {
		m.defaultChannel = name
	}
}

// NewManager creates a LogManager with the given configuration.
func NewManager(cfg *config.Repository, opts ...ManagerOption) *LogManager {
	m := &LogManager{
		cfg:           cfg,
		channels:      make(map[string]*Logger),
		customDrivers: make(map[string]DriverFactory),
		sharedContext: make(map[string]any),
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.defaultChannel == "" {
		m.defaultChannel = stringOr(cfg, "logging.default", "single")
	}

	return m
}

// Channel returns a cached logger for the given channel, creating it on first
// access. If no name is given, the default channel is used.
func (m *LogManager) Channel(name ...string) (*Logger, error) {
	n := m.defaultChannel

	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}

	return m.get(n)
}

// Driver is an alias for Channel.
func (m *LogManager) Driver(name ...string) (*Logger, error) {
	return m.Channel(name...)
}

// Stack creates a Logger that writes to all given channels simultaneously.
func (m *LogManager) Stack(channels []string, channel ...string) (*Logger, error) {
	name := "stack"

	if len(channel) > 0 && channel[0] != "" {
		name = channel[0]
	}

	var handlers []Handler

	for _, ch := range channels {
		logger, err := m.Channel(ch)

		if err != nil {
			return nil, fmt.Errorf("log: failed to resolve stack channel %q: %w", ch, err)
		}

		handlers = append(handlers, logger.GetHandler())
	}

	handler := NewStackHandler(handlers, LevelDebug)
	logger := NewLogger(handler, name)

	if m.dispatcher != nil {
		logger.SetEventDispatcher(m.dispatcher)
	}

	m.applySharedContext(logger)

	return logger, nil
}

// Build creates a Logger from the given config without caching it.
func (m *LogManager) Build(cc ChannelConfig) (*Logger, error) {
	handler, err := m.createDriver(cc)

	if err != nil {
		return nil, err
	}

	logger := NewLogger(handler, "ondemand")

	if m.dispatcher != nil {
		logger.SetEventDispatcher(m.dispatcher)
	}

	m.applySharedContext(logger)

	return logger, nil
}

// Extend registers a custom driver factory.
func (m *LogManager) Extend(driver string, factory DriverFactory) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.customDrivers[driver] = factory
}

// GetDefaultDriver returns the name of the default channel.
func (m *LogManager) GetDefaultDriver() string {
	return m.defaultChannel
}

// SetDefaultDriver sets the default channel name.
func (m *LogManager) SetDefaultDriver(name string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.defaultChannel = name
}

// ShareContext merges the given context into the shared context that is applied
// to all channels.
func (m *LogManager) ShareContext(ctx map[string]any) {
	m.mu.Lock()

	defer m.mu.Unlock()

	for k, v := range ctx {
		m.sharedContext[k] = v
	}

	for _, logger := range m.channels {
		m.applySharedContextLocked(logger)
	}
}

// SharedContext returns a copy of the current shared context.
func (m *LogManager) SharedContext() map[string]any {
	m.mu.RLock()

	defer m.mu.RUnlock()

	cp := make(map[string]any, len(m.sharedContext))

	for k, v := range m.sharedContext {
		cp[k] = v
	}

	return cp
}

// WithoutContext clears specific keys from the shared context. If no keys are
// given, all shared context is cleared.
func (m *LogManager) WithoutContext(keys ...string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	if len(keys) == 0 {
		m.sharedContext = make(map[string]any)
	} else {
		for _, key := range keys {
			delete(m.sharedContext, key)
		}
	}
}

// FlushSharedContext clears all shared context.
func (m *LogManager) FlushSharedContext() {
	m.WithoutContext()
}

// ForgetChannel removes a cached channel and closes its handler.
func (m *LogManager) ForgetChannel(names ...string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	if len(names) == 0 {
		names = []string{m.defaultChannel}
	}

	for _, name := range names {
		if logger, ok := m.channels[name]; ok {
			logger.GetHandler().Close()
			delete(m.channels, name)
		}
	}
}

// GetChannels returns a copy of the cached channel map.
func (m *LogManager) GetChannels() map[string]*Logger {
	m.mu.RLock()

	defer m.mu.RUnlock()

	cp := make(map[string]*Logger, len(m.channels))

	for k, v := range m.channels {
		cp[k] = v
	}

	return cp
}

// Tap registers callbacks that are applied to loggers after creation.
func (m *LogManager) Tap(callbacks ...func(*Logger)) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.tapCallbacks = append(m.tapCallbacks, callbacks...)
}

// Emergency logs at LevelEmergency on the default channel.
func (m *LogManager) Emergency(message string, context ...map[string]any) {
	m.log(LevelEmergency, message, context)
}

// Alert logs at LevelAlert on the default channel.
func (m *LogManager) Alert(message string, context ...map[string]any) {
	m.log(LevelAlert, message, context)
}

// Critical logs at LevelCritical on the default channel.
func (m *LogManager) Critical(message string, context ...map[string]any) {
	m.log(LevelCritical, message, context)
}

// Error logs at LevelError on the default channel.
func (m *LogManager) Error(message string, context ...map[string]any) {
	m.log(LevelError, message, context)
}

// Warning logs at LevelWarning on the default channel.
func (m *LogManager) Warning(message string, context ...map[string]any) {
	m.log(LevelWarning, message, context)
}

// Notice logs at LevelNotice on the default channel.
func (m *LogManager) Notice(message string, context ...map[string]any) {
	m.log(LevelNotice, message, context)
}

// Info logs at LevelInfo on the default channel.
func (m *LogManager) Info(message string, context ...map[string]any) {
	m.log(LevelInfo, message, context)
}

// Debug logs at LevelDebug on the default channel.
func (m *LogManager) Debug(message string, context ...map[string]any) {
	m.log(LevelDebug, message, context)
}

// Log logs at the given level on the default channel.
func (m *LogManager) Log(level Level, message string, context ...map[string]any) {
	m.log(level, message, context)
}

func (m *LogManager) log(level Level, message string, context []map[string]any) {
	logger, err := m.Channel()

	if err != nil {
		logger = m.emergencyLogger()
	}

	logger.writeLog(level, message, context)
}

func (m *LogManager) get(name string) (*Logger, error) {
	m.mu.RLock()

	if logger, ok := m.channels[name]; ok {
		m.mu.RUnlock()

		return logger, nil
	}

	m.mu.RUnlock()

	return m.resolve(name)
}

func (m *LogManager) resolve(name string) (*Logger, error) {
	cc, err := ParseChannelConfig(m.cfg, name)

	if err != nil {
		return nil, err
	}

	handler, err := m.createDriver(cc)

	if err != nil {
		return nil, fmt.Errorf("log: failed to create driver for channel %q: %w", name, err)
	}

	logger := NewLogger(handler, name)

	if m.dispatcher != nil {
		logger.SetEventDispatcher(m.dispatcher)
	}

	m.mu.Lock()
	m.applySharedContextLocked(logger)

	for _, tap := range m.tapCallbacks {
		tap(logger)
	}

	m.channels[name] = logger
	m.mu.Unlock()

	return logger, nil
}

func (m *LogManager) createDriver(cc ChannelConfig) (Handler, error) {
	m.mu.RLock()
	factory, hasCustom := m.customDrivers[string(cc.Driver)]
	m.mu.RUnlock()

	if hasCustom {
		return factory(cc)
	}

	switch cc.Driver {
	case DriverSingle:
		return m.createSingleDriver(cc)
	case DriverDaily:
		return m.createDailyDriver(cc)
	case DriverErrorlog:
		return m.createErrorlogDriver(cc)
	case DriverStack:
		return m.createStackDriver(cc)
	case DriverNull:
		return m.createNullDriver(cc)
	case DriverCustom:
		return m.createCustomDriver(cc)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedDriver, cc.Driver)
	}
}

func (m *LogManager) createSingleDriver(cc ChannelConfig) (Handler, error) {
	if cc.Path == "" {
		return nil, ErrMissingPath
	}

	handler, err := NewFileStreamHandler(cc.Path, cc.Level, 0644)

	if err != nil {
		return nil, err
	}

	if cc.Formatter != nil {
		handler.SetFormatter(cc.Formatter)
	}

	for _, p := range cc.Processors {
		handler.AddProcessor(p)
	}

	return handler, nil
}

func (m *LogManager) createDailyDriver(cc ChannelConfig) (Handler, error) {
	if cc.Path == "" {
		return nil, ErrMissingPath
	}

	handler := NewRotatingHandler(cc.Path, cc.Days, cc.Level)

	if cc.Formatter != nil {
		handler.SetFormatter(cc.Formatter)
	}

	for _, p := range cc.Processors {
		handler.AddProcessor(p)
	}

	return handler, nil
}

func (m *LogManager) createErrorlogDriver(cc ChannelConfig) (Handler, error) {
	handler := NewStderrHandler(cc.Level)

	if cc.Formatter != nil {
		handler.SetFormatter(cc.Formatter)
	}

	for _, p := range cc.Processors {
		handler.AddProcessor(p)
	}

	return handler, nil
}

func (m *LogManager) createStackDriver(cc ChannelConfig) (Handler, error) {
	var handlers []Handler

	for _, ch := range cc.Channels {
		logger, err := m.get(ch)

		if err != nil {
			return nil, fmt.Errorf("log: failed to resolve stack channel %q: %w", ch, err)
		}

		handlers = append(handlers, logger.GetHandler())
	}

	return NewStackHandler(handlers, cc.Level), nil
}

func (m *LogManager) createNullDriver(_ ChannelConfig) (Handler, error) {
	return NewNullHandler(), nil
}

func (m *LogManager) createCustomDriver(cc ChannelConfig) (Handler, error) {
	if cc.Via == nil {
		return nil, fmt.Errorf("log: custom driver requires a Via factory function")
	}

	return cc.Via(cc)
}

func (m *LogManager) emergencyLogger() *Logger {
	handler := NewStreamHandler(os.Stderr, LevelDebug)

	return NewLogger(handler, "emergency")
}

func (m *LogManager) applySharedContext(logger *Logger) {
	m.mu.RLock()

	defer m.mu.RUnlock()

	m.applySharedContextLocked(logger)
}

func (m *LogManager) applySharedContextLocked(logger *Logger) {
	for k, v := range m.sharedContext {
		logger.context[k] = v
	}
}
