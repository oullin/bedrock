package queue_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
)

// Ports of Framework\Tests\Queue\QueueManagerTest.
//
// Upstream's test uses Mockery on stdClass to assert that the Manager:
//   1. looks up config['queue.connections.{name}'] (or falls back for null),
//   2. instantiates the Connector factory registered via addConnector,
//   3. calls Connector::connect($config) to obtain the Queue,
//   4. calls Queue::setConnectionName($name),
//   5. calls Queue::setContainer($app),
//   6. caches the instance for subsequent connection() calls.
//
// The Go port wires the same flow through the real Manager, a fake
// connector, and a fake queue that implements ConnectionNameSetter and
// ContainerAware and records its calls.

// --- fixtures ---------------------------------------------------------

// recordedQueue captures the Manager's setConnectionName / setContainer
// calls while satisfying queue.Queue with no-op method implementations.
type recordedQueue struct {
	mu             sync.Mutex
	connectionName string
	container      any
}

// Minimal queue.Queue satisfying methods (all no-ops / empty returns).

// recordedConnector captures the config passed to Connect and returns
// the preconfigured queue exactly once.
type recordedConnector struct {
	mu           sync.Mutex
	queue        queue.Queue
	connectCalls int
	lastConfig   map[string]any
}

// fakeConnectionEnum is the Go stand-in for PHP's `enum
// QueueConnectionName: string { case Sync = 'sync'; }`. Any type with
// a String() string method is accepted by Manager.Connection as an
// enum-like reference.
type fakeConnectionEnum string

func (q *recordedQueue) SetConnectionName(name string) {
	q.mu.Lock()

	defer q.mu.Unlock()

	q.connectionName = name
}

func (q *recordedQueue) SetContainer(container any) {
	q.mu.Lock()

	defer q.mu.Unlock()

	q.container = container
}

func (q *recordedQueue) ConnectionName() string {
	q.mu.Lock()

	defer q.mu.Unlock()

	return q.connectionName
}

func (q *recordedQueue) Push(context.Context, string, []byte) (string, error) {
	return "", nil
}

func (q *recordedQueue) PushDelayed(context.Context, string, []byte, time.Duration) (string, error) {
	return "", nil
}

func (q *recordedQueue) PushMultiple(context.Context, string, [][]byte) ([]string, error) {
	return nil, nil
}

func (q *recordedQueue) Pop(context.Context, string) (queue.Job, error) {
	return nil, queue.ErrNoJob
}

func (q *recordedQueue) Size(context.Context, string) (int64, error) { return 0, nil }

func (q *recordedQueue) PendingSize(context.Context, string) (int64, error) { return 0, nil }

func (q *recordedQueue) DelayedSize(context.Context, string) (int64, error) { return 0, nil }

func (q *recordedQueue) ReservedSize(context.Context, string) (int64, error) { return 0, nil }

func (c *recordedConnector) Connect(config map[string]any) (queue.Queue, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.connectCalls++
	c.lastConfig = config

	return c.queue, nil
}

func (e fakeConnectionEnum) String() string { return string(e) }

const fakeConnectionSync fakeConnectionEnum = "sync"

// --- ports ------------------------------------------------------------

// Port of Framework\Tests\Queue\QueueManagerTest::testDefaultConnectionCanBeResolved
func TestDefaultConnectionCanBeResolved(t *testing.T) {
	t.Parallel()

	app := map[string]any{"name": "bedrock-app"}
	fakeQueue := &recordedQueue{}
	connector := &recordedConnector{queue: fakeQueue}

	manager := queue.NewManager()
	manager.SetContainer(app)
	manager.SetConfig("sync", map[string]any{"driver": "sync"})
	manager.SetDefaultConnection("sync")
	manager.AddConnector("sync", func() queue.Connector { return connector })

	got, err := manager.Connection("sync")

	if err != nil {
		t.Fatalf("Connection: %v", err)
	}

	if got != fakeQueue {
		t.Error("expected Connection to return the queue the connector produced")
	}

	if connector.connectCalls != 1 {
		t.Errorf("Connect calls: got %d, want 1", connector.connectCalls)
	}

	if driver, _ := connector.lastConfig["driver"].(string); driver != "sync" {
		t.Errorf("Connect config: got %v, want {driver:sync}", connector.lastConfig)
	}

	if fakeQueue.ConnectionName() != "sync" {
		t.Errorf("setConnectionName: got %q, want sync", fakeQueue.ConnectionName())
	}

	if fakeQueue.container == nil {
		t.Error("setContainer was not called")
	}
}

// Port of Framework\Tests\Queue\QueueManagerTest::testOtherConnectionCanBeResolved
func TestOtherConnectionCanBeResolved(t *testing.T) {
	t.Parallel()

	app := map[string]any{"name": "bedrock-app"}
	fakeQueue := &recordedQueue{}
	connector := &recordedConnector{queue: fakeQueue}

	manager := queue.NewManager()
	manager.SetContainer(app)
	manager.SetConfig("foo", map[string]any{"driver": "bar"})
	manager.SetDefaultConnection("sync")
	manager.AddConnector("bar", func() queue.Connector { return connector })

	got, err := manager.Connection("foo")

	if err != nil {
		t.Fatalf("Connection: %v", err)
	}

	if got != fakeQueue {
		t.Error("expected Connection to return the queue the connector produced")
	}

	if driver, _ := connector.lastConfig["driver"].(string); driver != "bar" {
		t.Errorf("Connect config: got %v, want {driver:bar}", connector.lastConfig)
	}

	// Upstream stamps the CONNECTION name (foo) not the driver name (bar).
	if fakeQueue.ConnectionName() != "foo" {
		t.Errorf("setConnectionName: got %q, want foo", fakeQueue.ConnectionName())
	}
}

// Port of Framework\Tests\Queue\QueueManagerTest::testNullConnectionCanBeResolved
func TestNullConnectionCanBeResolved(t *testing.T) {
	t.Parallel()

	app := map[string]any{"name": "bedrock-app"}
	fakeQueue := &recordedQueue{}
	connector := &recordedConnector{queue: fakeQueue}

	manager := queue.NewManager()
	manager.SetContainer(app)
	manager.SetDefaultConnection("null")
	manager.AddConnector("null", func() queue.Connector { return connector })

	got, err := manager.Connection("null")

	if err != nil {
		t.Fatalf("Connection: %v", err)
	}

	if got != fakeQueue {
		t.Error("expected Connection to return the queue the connector produced")
	}

	// Upstream synthesises ['driver' => 'null'] when no explicit config exists.
	if driver, _ := connector.lastConfig["driver"].(string); driver != "null" {
		t.Errorf("Connect config: got %v, want {driver:null}", connector.lastConfig)
	}

	if fakeQueue.ConnectionName() != "null" {
		t.Errorf("setConnectionName: got %q, want null", fakeQueue.ConnectionName())
	}
}

// Port of Framework\Tests\Queue\QueueManagerTest::testEnumConnectionCanBeResolved
func TestEnumConnectionCanBeResolved(t *testing.T) {
	t.Parallel()

	app := map[string]any{"name": "bedrock-app"}
	fakeQueue := &recordedQueue{}
	connector := &recordedConnector{queue: fakeQueue}

	manager := queue.NewManager()
	manager.SetContainer(app)
	manager.SetConfig("sync", map[string]any{"driver": "sync"})
	manager.SetDefaultConnection("sync")
	manager.AddConnector("sync", func() queue.Connector { return connector })

	got, err := manager.Connection(fakeConnectionSync)

	if err != nil {
		t.Fatalf("Connection: %v", err)
	}

	if got != fakeQueue {
		t.Error("expected Connection(enum) to return the queue the connector produced")
	}

	if fakeQueue.ConnectionName() != "sync" {
		t.Errorf("setConnectionName: got %q, want sync", fakeQueue.ConnectionName())
	}
}

// Port of Framework\Tests\Queue\QueueManagerTest::testEnumConnectionCanBeChecked
func TestEnumConnectionCanBeChecked(t *testing.T) {
	t.Parallel()

	fakeQueue := &recordedQueue{}
	connector := &recordedConnector{queue: fakeQueue}

	manager := queue.NewManager()
	manager.SetContainer(map[string]any{})
	manager.SetConfig("sync", map[string]any{"driver": "sync"})
	manager.SetDefaultConnection("sync")
	manager.AddConnector("sync", func() queue.Connector { return connector })

	if manager.Connected(fakeConnectionSync) {
		t.Error("Connected: expected false before first resolve")
	}

	if _, err := manager.Connection(fakeConnectionSync); err != nil {
		t.Fatalf("Connection: %v", err)
	}

	if !manager.Connected(fakeConnectionSync) {
		t.Error("Connected: expected true after first resolve")
	}
}
