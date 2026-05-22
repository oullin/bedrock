package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// DriverCreator is the legacy one-shot creator: a function that takes a
// raw config map and returns a ready-to-use Queue. It is the simpler of
// the two registration paths — drivers that do not need a separate
// Connector step can register this way.
type DriverCreator func(config map[string]any) (Queue, error)

// ConnectorFactory is the upstream-faithful two-step path: the factory
// returns a Connector, and the Manager then calls Connector.Connect(config)
// to obtain a Queue.
// Ref: @bedrock/code-0231
type ConnectorFactory func() Connector

// ConnectionNameSetter is the optional contract a Queue implementation
// can satisfy if it wants the Manager to stamp the resolved connection
// name onto it immediately after creation.
// the upstream Queue::setConnectionName — exposed as an optional interface
// so the core Queue interface stays frozen.
type ConnectionNameSetter interface {
	SetConnectionName(name string)
}

// ContainerAware is the optional contract a Queue implementation can
// satisfy if it wants the Manager to hand it the application container
// (or any opaque value the caller chose as container) after creation.
type ContainerAware interface {
	SetContainer(container any)
}

// HookFunc is a generic event listener registered against the Manager.
// The underlying event is passed as any; listeners type-assert to the
// concrete events. * type they care about.
type HookFunc func(event any)

// Manager creates, caches, and coordinates named queue connections.
//
// Ref: @bedrock/code-0269
// two entry points for registering drivers (Register for the simple
// creator path, AddConnector for the two-step Connector path), enum-like
// connection references via Connection(any), optional queue hooks for
// setConnectionName and setContainer, pause/resume delegation through
// an embedded PauseResumer, and before/after/failing/starting/stopping
// hook registration for the worker pipeline.
type Manager struct {
	mu                sync.RWMutex
	queues            map[string]Queue
	creators          map[string]DriverCreator
	factories         map[string]ConnectorFactory
	configs           map[string]map[string]any
	defaultConnection string
	container         any

	// Worker-side lifecycle hooks.
	beforeHooks   []HookFunc
	afterHooks    []HookFunc
	failingHooks  []HookFunc
	startingHooks []HookFunc
	stoppingHooks []HookFunc

	// Pause/resume is delegated to a shared PauseResumer so the unit
	// tests written in Step 5 keep passing when Manager grows up.
	pauseResumer *PauseResumer
	pauseStore   PauseStore
	emitter      EventEmitter
}

// NewManager constructs an empty Manager with an in-memory pause store
// and no registered drivers or connectors.
func NewManager() *Manager {
	store := NewInMemoryPauseStore()
	m := &Manager{
		queues:       make(map[string]Queue),
		creators:     make(map[string]DriverCreator),
		factories:    make(map[string]ConnectorFactory),
		configs:      make(map[string]map[string]any),
		pauseStore:   store,
		pauseResumer: NewPauseResumer(store, nil),
	}

	return m
}

// --- driver registration ---------------------------------------------

// Register binds a legacy DriverCreator to the given driver name. This
// is the existing bedrock-native path; new code should prefer AddConnector.
func (m *Manager) Register(driver string, creator DriverCreator) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.creators[driver] = creator

	return m
}

// Extend is an alias for Register. Kept because the upstream QueueManager
// exposes extend() as well, and existing Go tests call it.
func (m *Manager) Extend(driver string, creator DriverCreator) *Manager {
	return m.Register(driver, creator)
}

// AddConnector registers a ConnectorFactory for the given driver name.
// The factory is invoked lazily the first time a connection using that
// Ref: @bedrock/code-0269
func (m *Manager) AddConnector(driver string, factory ConnectorFactory) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.factories[driver] = factory

	return m
}

// ForgetDriver removes both the creator and the connector factory
// registered for the given driver name. It is a no-op if neither is
// registered.
func (m *Manager) ForgetDriver(driver string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	delete(m.creators, driver)
	delete(m.factories, driver)
}

// --- config & container ----------------------------------------------

// SetConfig stores the configuration map for a named connection.
func (m *Manager) SetConfig(connection string, config map[string]any) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.configs[connection] = config

	return m
}

// SetContainer stores an opaque container value that the Manager will
// hand to every ContainerAware queue it creates.
func (m *Manager) SetContainer(container any) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.container = container

	return m
}

// --- connection resolution -------------------------------------------

// Driver returns (or lazily creates) the Queue for the given connection
// name. This is the string-typed, existing bedrock API.
func (m *Manager) Driver(connection string) (Queue, error) {
	return m.resolveConnection(connection)
}

// Connection returns (or lazily creates) the Queue for the given
// connection name. Accepts either a plain string or any value that
// implements fmt.Stringer — matching the upstream enum-or-string handling.
// Passing nil or an empty string resolves the default connection.
func (m *Manager) Connection(name any) (Queue, error) {
	key := connectionKey(name)

	if key == "" {
		key = m.GetDefaultConnection()
	}

	if key == "" {
		return nil, fmt.Errorf("queue: Connection: no name supplied and no default connection configured")
	}

	return m.resolveConnection(key)
}

// Connected reports whether the given connection has been resolved and
// cached. It does not trigger creation.
func (m *Manager) Connected(name any) bool {
	key := connectionKey(name)

	if key == "" {
		key = m.GetDefaultConnection()
	}

	m.mu.RLock()

	defer m.mu.RUnlock()

	_, ok := m.queues[key]

	return ok
}

// GetDefaultConnection returns the registered default connection name.
func (m *Manager) GetDefaultConnection() string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.defaultConnection
}

// SetDefaultConnection sets the default connection name.
func (m *Manager) SetDefaultConnection(connection string) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.defaultConnection = connection

	return m
}

// Purge removes the cached Queue for connection, forcing re-creation
// on the next resolve.
func (m *Manager) Purge(connection string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	delete(m.queues, connection)
}

// resolveConnection performs the unified creation flow shared by
// Driver and Connection: cache lookup → config lookup → connector or
// creator invocation → optional SetConnectionName / SetContainer hooks
// → cache and return.
func (m *Manager) resolveConnection(connection string) (Queue, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	if q, ok := m.queues[connection]; ok {
		return q, nil
	}

	cfg, ok := m.lookupConfigLocked(connection)

	if !ok {
		return nil, fmt.Errorf("queue: resolve %q: no config registered", connection)
	}

	driver, _ := cfg["driver"].(string)

	q, err := m.instantiateLocked(driver, cfg)

	if err != nil {
		return nil, fmt.Errorf("queue: create driver %q: %w", driver, err)
	}

	if s, ok := q.(ConnectionNameSetter); ok {
		s.SetConnectionName(connection)
	}

	if s, ok := q.(ContainerAware); ok {
		s.SetContainer(m.container)
	}

	m.queues[connection] = q

	return q, nil
}

// lookupConfigLocked returns the config for name, synthesising the
// Upstream-special "null" fallback if the caller did not SetConfig for
// it. Caller holds m.mu.
func (m *Manager) lookupConfigLocked(name string) (map[string]any, bool) {
	if cfg, ok := m.configs[name]; ok {
		return cfg, true
	}

	if name == "null" {
		return map[string]any{"driver": "null"}, true
	}

	return nil, false
}

// instantiateLocked runs the appropriate creation path for driver,
// preferring the upstream-faithful ConnectorFactory path over the legacy
// DriverCreator path. Caller holds m.mu.
func (m *Manager) instantiateLocked(driver string, cfg map[string]any) (Queue, error) {
	if factory, ok := m.factories[driver]; ok {
		connector := factory()

		if connector == nil {
			return nil, fmt.Errorf("%w: %q factory returned nil connector", ErrInvalidDriver, driver)
		}

		q, err := connector.Connect(cfg)

		if err != nil {
			return nil, err
		}

		if q == nil {
			return nil, fmt.Errorf("%w: %q connector returned nil queue", ErrInvalidDriver, driver)
		}

		return q, nil
	}

	if creator, ok := m.creators[driver]; ok {
		q, err := creator(cfg)

		if err != nil {
			return nil, err
		}

		if q == nil {
			return nil, fmt.Errorf("%w: %q creator returned nil queue", ErrInvalidDriver, driver)
		}

		return q, nil
	}

	return nil, fmt.Errorf("%w: %q", ErrInvalidDriver, driver)
}

// --- hook registration ------------------------------------------------

// Before registers a listener that runs before each job is processed.
func (m *Manager) Before(hook HookFunc) *Manager {
	return m.appendHook(&m.beforeHooks, hook)
}

// After registers a listener that runs after a job has been processed.
func (m *Manager) After(hook HookFunc) *Manager { return m.appendHook(&m.afterHooks, hook) }

// Failing registers a listener that runs when a job fails.
// the upstream Queue::failing.
func (m *Manager) Failing(hook HookFunc) *Manager { return m.appendHook(&m.failingHooks, hook) }

// Starting registers a listener that runs when a worker daemon boots.
func (m *Manager) Starting(hook HookFunc) *Manager { return m.appendHook(&m.startingHooks, hook) }

// Stopping registers a listener that runs when a worker daemon exits.
func (m *Manager) Stopping(hook HookFunc) *Manager { return m.appendHook(&m.stoppingHooks, hook) }

// BeforeHooks / AfterHooks / FailingHooks / StartingHooks / StoppingHooks
// return snapshots of the registered hook chains. They exist so the
// Worker (Step 8) can fan out without holding the Manager's lock during
// dispatch.

// BeforeHooks returns a snapshot of the before-hook chain.
func (m *Manager) BeforeHooks() []HookFunc { return m.snapshotHooks(m.beforeHooks) }

// AfterHooks returns a snapshot of the after-hook chain.
func (m *Manager) AfterHooks() []HookFunc { return m.snapshotHooks(m.afterHooks) }

// FailingHooks returns a snapshot of the failing-hook chain.
func (m *Manager) FailingHooks() []HookFunc { return m.snapshotHooks(m.failingHooks) }

// StartingHooks returns a snapshot of the starting-hook chain.
func (m *Manager) StartingHooks() []HookFunc { return m.snapshotHooks(m.startingHooks) }

// StoppingHooks returns a snapshot of the stopping-hook chain.
func (m *Manager) StoppingHooks() []HookFunc { return m.snapshotHooks(m.stoppingHooks) }

func (m *Manager) appendHook(slot *[]HookFunc, hook HookFunc) *Manager {
	if hook == nil {
		return m
	}

	m.mu.Lock()

	defer m.mu.Unlock()

	*slot = append(*slot, hook)

	return m
}

func (m *Manager) snapshotHooks(slot []HookFunc) []HookFunc {
	m.mu.RLock()

	defer m.mu.RUnlock()

	out := make([]HookFunc, len(slot))
	copy(out, slot)

	return out
}

// --- pause/resume delegation -----------------------------------------

// SetEmitter wires an EventEmitter into the Manager so pause/resume
// events (and eventually worker events) are dispatched through it. A
// nil emitter drops events on the floor.
func (m *Manager) SetEmitter(e EventEmitter) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.emitter = e
	m.pauseResumer = NewPauseResumer(m.pauseStore, e)

	return m
}

// WithPauseStore swaps the underlying PauseStore. Subsequent
// pause/resume calls use the new store; state in the old store is left
// untouched but becomes inaccessible through the Manager.
func (m *Manager) WithPauseStore(store PauseStore) *Manager {
	if store == nil {
		return m
	}

	m.mu.Lock()

	defer m.mu.Unlock()

	m.pauseStore = store
	m.pauseResumer = NewPauseResumer(store, m.emitter)

	return m
}

// Pause marks (connection, queue) as paused indefinitely.
func (m *Manager) Pause(connection, queue string) error {
	return m.pauseResumer.Pause(connection, queue)
}

// PauseFor marks (connection, queue) as paused until now+ttl.
func (m *Manager) PauseFor(connection, queue string, ttl time.Duration) error {
	return m.pauseResumer.PauseFor(connection, queue, ttl)
}

// Resume lifts any pause state for (connection, queue).
func (m *Manager) Resume(connection, queue string) error {
	return m.pauseResumer.Resume(connection, queue)
}

// IsPaused reports whether (connection, queue) is currently paused.
func (m *Manager) IsPaused(connection, queue string) bool {
	return m.pauseResumer.IsPaused(connection, queue)
}

// --- cross-queue inspection (upstream 13.8.0) -------------------------

// AllPendingJobs returns a snapshot of every pending job sitting on any
// queue belonging to connection. It resolves the connection, asks the
// driver for the set of queue names it currently knows about (via the
// optional QueueNamer contract), and concatenates the per-queue
// PendingJobs results in declared order.
// Queue::allPendingJobs.
func (m *Manager) AllPendingJobs(ctx context.Context, connection string) ([]InspectedJob, error) {
	return m.allJobs(ctx, connection, func(i JobInspector, name string) ([]InspectedJob, error) {
		return i.PendingJobs(ctx, name)
	})
}

// AllDelayedJobs returns every delayed (unreserved, not-yet-due) job
// across all queues on connection.
func (m *Manager) AllDelayedJobs(ctx context.Context, connection string) ([]InspectedJob, error) {
	return m.allJobs(ctx, connection, func(i JobInspector, name string) ([]InspectedJob, error) {
		return i.DelayedJobs(ctx, name)
	})
}

// AllReservedJobs returns every reserved (in-flight) job across all
// queues on connection.
func (m *Manager) AllReservedJobs(ctx context.Context, connection string) ([]InspectedJob, error) {
	return m.allJobs(ctx, connection, func(i JobInspector, name string) ([]InspectedJob, error) {
		return i.ReservedJobs(ctx, name)
	})
}

// allJobs implements the shared fan-out used by the three All*Jobs
// helpers. Drivers that do not implement QueueNamer or JobInspector
// surface as ErrNotSupported; per-queue calls that themselves return
// ErrNotSupported are skipped silently — the caller then receives only
// the upstream "best-effort across queues" semantics.
func (m *Manager) allJobs(ctx context.Context, connection string, fetch func(JobInspector, string) ([]InspectedJob, error)) ([]InspectedJob, error) {
	q, err := m.resolveConnection(connection)

	if err != nil {
		return nil, err
	}

	namer, ok := q.(QueueNamer)

	if !ok {
		return nil, fmt.Errorf("%w: connection %q has no QueueNamer", ErrNotSupported, connection)
	}

	inspector, ok := q.(JobInspector)

	if !ok {
		return nil, fmt.Errorf("%w: connection %q has no JobInspector", ErrNotSupported, connection)
	}

	names, err := namer.QueueNames(ctx)

	if err != nil {
		return nil, err
	}

	var out []InspectedJob

	var errs []error

	for _, name := range names {
		jobs, err := fetch(inspector, name)

		if err != nil {
			if errors.Is(err, ErrNotSupported) {
				continue
			}

			errs = append(errs, err)

			continue
		}

		out = append(out, jobs...)
	}

	if len(errs) > 0 {
		return out, errors.Join(errs...)
	}

	return out, nil
}

// --- helpers ----------------------------------------------------------

// connectionKey unwraps a connection reference to its string name. It
// accepts:
//
//   - nil       → ""
//   - string    → the string verbatim
//   - Stringer  → the result of String()
//   - anything  → ""
//
// The Stringer branch is the Go analogue of accepting a PHP BackedEnum.
func connectionKey(ref any) string {
	if ref == nil {
		return ""
	}

	switch v := ref.(type) {
	case string:
		return v
	case interface{ String() string }:
		return v.String()
	default:
		return ""
	}
}
