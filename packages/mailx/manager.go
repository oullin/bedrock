package mailx

import (
	"fmt"
	"sync"

	"github.com/bedrock/packages/config"
	cevents "github.com/bedrock/packages/contracts/events"
	clog "github.com/bedrock/packages/contracts/log"
)

// TransportFactory creates a Transport from a mailer configuration.
type TransportFactory func(config map[string]any) (Transport, error)

// MailManager creates, caches, and delegates to named Mailer instances.
type MailManager struct {
	mu            sync.RWMutex
	cfg           *config.Repository
	dispatcher    cevents.Dispatcher
	logger        clog.Logger
	mailers       map[string]*Mailer
	customDrivers map[string]TransportFactory
	defaultMailer string
}

// ManagerOption configures a MailManager.
type ManagerOption func(*MailManager)

// WithManagerDispatcher sets the event dispatcher on the manager.
func WithManagerDispatcher(d cevents.Dispatcher) ManagerOption {
	return func(m *MailManager) {
		m.dispatcher = d
	}
}

// WithLogger sets the logger used by LogTransport.
func WithLogger(l clog.Logger) ManagerOption {
	return func(m *MailManager) {
		m.logger = l
	}
}

// WithDefaultMailer sets the default mailer name.
func WithDefaultMailer(name string) ManagerOption {
	return func(m *MailManager) {
		m.defaultMailer = name
	}
}

// NewManager creates a MailManager with the given configuration.
func NewManager(cfg *config.Repository, opts ...ManagerOption) *MailManager {
	m := &MailManager{
		cfg:           cfg,
		mailers:       make(map[string]*Mailer),
		customDrivers: make(map[string]TransportFactory),
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.defaultMailer == "" {
		v := cfg.Get("mail.default", "smtp")

		if s, ok := v.(string); ok {
			m.defaultMailer = s
		} else {
			m.defaultMailer = "smtp"
		}
	}

	return m
}

// Mailer returns a cached mailer by name, creating it on first access.
// If no name is given, the default mailer is used.
func (m *MailManager) Mailer(name ...string) (*Mailer, error) {
	n := m.defaultMailer

	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}

	return m.get(n)
}

// Driver is an alias for Mailer.
func (m *MailManager) Driver(name ...string) (*Mailer, error) {
	return m.Mailer(name...)
}

// Build creates a Mailer from the given config without caching.
func (m *MailManager) Build(cfg map[string]any) (*Mailer, error) {
	transport, err := m.createTransport(cfg)

	if err != nil {
		return nil, err
	}

	mailer := NewMailer("ondemand", transport)

	if m.dispatcher != nil {
		mailer.dispatcher = m.dispatcher
	}

	return mailer, nil
}

// Extend registers a custom transport factory.
func (m *MailManager) Extend(driver string, factory TransportFactory) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.customDrivers[driver] = factory
}

// Purge removes a cached mailer instance.
func (m *MailManager) Purge(name ...string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	n := m.defaultMailer

	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}

	delete(m.mailers, n)
}

// ForgetMailers removes all cached mailer instances.
func (m *MailManager) ForgetMailers() {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.mailers = make(map[string]*Mailer)
}

// GetDefaultDriver returns the default mailer name.
func (m *MailManager) GetDefaultDriver() string {
	return m.defaultMailer
}

// SetDefaultDriver sets the default mailer name.
func (m *MailManager) SetDefaultDriver(name string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.defaultMailer = name
}

// GetMailers returns a copy of the cached mailer map.
func (m *MailManager) GetMailers() map[string]*Mailer {
	m.mu.RLock()

	defer m.mu.RUnlock()

	cp := make(map[string]*Mailer, len(m.mailers))

	for k, v := range m.mailers {
		cp[k] = v
	}

	return cp
}

func (m *MailManager) get(name string) (*Mailer, error) {
	m.mu.RLock()

	if mailer, ok := m.mailers[name]; ok {
		m.mu.RUnlock()

		return mailer, nil
	}

	m.mu.RUnlock()

	return m.resolve(name)
}

func (m *MailManager) resolve(name string) (*Mailer, error) {
	cfg := m.getMailerConfig(name)

	if cfg == nil {
		return nil, fmt.Errorf("%w: %s", ErrMailerNotFound, name)
	}

	transport, err := m.createTransport(cfg)

	if err != nil {
		return nil, fmt.Errorf("mail: failed to create transport for mailer %q: %w", name, err)
	}

	mailer := NewMailer(name, transport)

	if m.dispatcher != nil {
		mailer.dispatcher = m.dispatcher
	}

	// Apply global from address if configured.
	fromAddr := m.cfg.Get("mail.from.address", "")
	fromName := m.cfg.Get("mail.from.name", "")

	if addr, ok := fromAddr.(string); ok && addr != "" {
		n := ""

		if nm, ok := fromName.(string); ok {
			n = nm
		}

		mailer.AlwaysFrom(addr, n)
	}

	m.mu.Lock()
	m.mailers[name] = mailer
	m.mu.Unlock()

	return mailer, nil
}

func (m *MailManager) getMailerConfig(name string) map[string]any {
	key := "mail.mailers." + name
	v := m.cfg.Get(key, nil)

	if v == nil {
		return nil
	}

	if cfg, ok := v.(map[string]any); ok {
		return cfg
	}

	return nil
}

func (m *MailManager) createTransport(cfg map[string]any) (Transport, error) {
	driver, _ := cfg["transport"].(string)

	if driver == "" {
		driver, _ = cfg["driver"].(string)
	}

	// Check custom drivers first.
	m.mu.RLock()
	factory, hasCustom := m.customDrivers[driver]
	m.mu.RUnlock()

	if hasCustom {
		return factory(cfg)
	}

	switch driver {
	case "smtp":
		return m.createSMTPTransport(cfg), nil
	case "log":
		return m.createLogTransport(), nil
	case "array":
		return NewArrayTransport(), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidDriver, driver)
	}
}

func (m *MailManager) createSMTPTransport(cfg map[string]any) Transport {
	host, _ := cfg["host"].(string)
	port := 587

	if p, ok := cfg["port"].(float64); ok {
		port = int(p)
	} else if p, ok := cfg["port"].(int); ok {
		port = p
	}

	username, _ := cfg["username"].(string)
	password, _ := cfg["password"].(string)

	var opts []SMTPOption

	if enc, ok := cfg["encryption"].(string); ok && enc != "" {
		opts = append(opts, WithEncryption(enc))
	}

	return NewSMTPTransport(host, port, username, password, opts...)
}

func (m *MailManager) createLogTransport() Transport {
	if m.logger != nil {
		return NewLogTransport(m.logger)
	}

	return NewArrayTransport()
}
