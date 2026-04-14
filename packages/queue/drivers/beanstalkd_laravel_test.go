package drivers_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

// Ports of Illuminate\Tests\Queue\QueueBeanstalkdQueueTest.
//
// Laravel's tests use Mockery on a Pheanstalk transport; they assert
// that push/later/pop/delete reach the client with exact priorities,
// delays, TTRs, and reserve timeouts. The Go port uses an inline
// recording BeanstalkdClient that captures every method call so the
// same observable behaviours can be asserted.
//
// Adaptation rules applied (see PARITY.md §2):
//
//   - Laravel's BeanstalkdQueue fuses useTube()+put() as two wire calls.
//     The Go BeanstalkdClient.Put already takes the tube as an argument,
//     so the "useTube called twice" assertion becomes "Put called with
//     the expected tube twice".
//   - Laravel's test pins the payload body byte-for-byte via Str::uuid
//     and Carbon::setTestNow shims. The Go port asserts only the
//     driver-level routing parameters (tube, priority, delay, ttr) —
//     payload construction has its own dedicated port in Step 4's
//     abstract_queue_test.go.

// --- recording mock ---------------------------------------------------

type recordingBeanstalkdPut struct {
	tube     string
	body     []byte
	priority uint32
	delay    time.Duration
	ttr      time.Duration
}

type recordingBeanstalkdReserve struct {
	tube    string
	timeout time.Duration
}

// recordingBeanstalkdClient captures every call so the ports can
// inspect a full history rather than just the last one.
type recordingBeanstalkdClient struct {
	mu       sync.Mutex
	puts     []recordingBeanstalkdPut
	reserves []recordingBeanstalkdReserve
	deletes  []uint64
	nextID   uint64

	// Canned reserve response.
	reserveID   uint64
	reserveBody []byte
	reserveErr  error
}

func (c *recordingBeanstalkdClient) Put(_ context.Context, tube string, body []byte, priority uint32, delay, ttr time.Duration) (uint64, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.nextID++
	c.puts = append(c.puts, recordingBeanstalkdPut{
		tube:     tube,
		body:     body,
		priority: priority,
		delay:    delay,
		ttr:      ttr,
	})

	return c.nextID, nil
}

func (c *recordingBeanstalkdClient) ReserveWithTimeout(_ context.Context, tube string, timeout time.Duration) (uint64, []byte, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.reserves = append(c.reserves, recordingBeanstalkdReserve{tube: tube, timeout: timeout})

	if c.reserveErr != nil {
		return 0, nil, c.reserveErr
	}

	return c.reserveID, c.reserveBody, nil
}

func (c *recordingBeanstalkdClient) Delete(_ context.Context, id uint64) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.deletes = append(c.deletes, id)

	return nil
}

func (c *recordingBeanstalkdClient) Release(_ context.Context, _ uint64, _ uint32, _ time.Duration) error {
	return nil
}

func (c *recordingBeanstalkdClient) Bury(_ context.Context, _ uint64, _ uint32) error {
	return nil
}

func (c *recordingBeanstalkdClient) StatsTube(_ context.Context, _ string) (map[string]string, error) {
	return map[string]string{}, nil
}

// --- test helpers -----------------------------------------------------

func newBeanstalkdDriverForPort(blockFor time.Duration) (*drivers.BeanstalkdDriver, *recordingBeanstalkdClient) {
	client := &recordingBeanstalkdClient{}
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second).
		SetDefaultTube("default").
		SetBlockFor(blockFor)

	return drv, client
}

// Port of Illuminate\Tests\Queue\QueueBeanstalkdQueueTest::testPushProperlyPushesJobOntoBeanstalkd
func TestPushProperlyPushesJobOntoBeanstalkd(t *testing.T) {
	t.Parallel()

	drv, client := newBeanstalkdDriverForPort(0)

	payload := []byte(`{"job":"foo","data":["data"]}`)

	if _, err := drv.Push(context.Background(), "stack", payload); err != nil {
		t.Fatalf("Push stack: %v", err)
	}

	if _, err := drv.Push(context.Background(), "", payload); err != nil {
		t.Fatalf("Push default: %v", err)
	}

	if len(client.puts) != 2 {
		t.Fatalf("put calls: got %d, want 2", len(client.puts))
	}

	if client.puts[0].tube != "stack" {
		t.Errorf("puts[0].tube: got %q, want stack", client.puts[0].tube)
	}

	if client.puts[1].tube != "default" {
		t.Errorf("puts[1].tube: got %q, want default (via defaultTube fallback)", client.puts[1].tube)
	}

	for i, p := range client.puts {
		if p.priority != 1024 {
			t.Errorf("puts[%d].priority: got %d, want 1024", i, p.priority)
		}

		if p.delay != 0 {
			t.Errorf("puts[%d].delay: got %s, want 0", i, p.delay)
		}

		if p.ttr != 60*time.Second {
			t.Errorf("puts[%d].ttr: got %s, want 60s", i, p.ttr)
		}
	}
}

// Port of Illuminate\Tests\Queue\QueueBeanstalkdQueueTest::testDelayedPushProperlyPushesJobOntoBeanstalkd
func TestDelayedPushProperlyPushesJobOntoBeanstalkd(t *testing.T) {
	t.Parallel()

	drv, client := newBeanstalkdDriverForPort(0)

	payload := []byte(`{"job":"foo","data":["data"]}`)

	if _, err := drv.PushDelayed(context.Background(), "stack", payload, 5*time.Second); err != nil {
		t.Fatalf("PushDelayed stack: %v", err)
	}

	if _, err := drv.PushDelayed(context.Background(), "", payload, 5*time.Second); err != nil {
		t.Fatalf("PushDelayed default: %v", err)
	}

	if len(client.puts) != 2 {
		t.Fatalf("put calls: got %d, want 2", len(client.puts))
	}

	if client.puts[0].tube != "stack" {
		t.Errorf("puts[0].tube: got %q, want stack", client.puts[0].tube)
	}

	if client.puts[1].tube != "default" {
		t.Errorf("puts[1].tube: got %q, want default", client.puts[1].tube)
	}

	for i, p := range client.puts {
		if p.priority != 1024 {
			t.Errorf("puts[%d].priority: got %d, want 1024", i, p.priority)
		}

		if p.delay != 5*time.Second {
			t.Errorf("puts[%d].delay: got %s, want 5s", i, p.delay)
		}

		// Delayed pushes use the driver's default TTR, matching
		// Laravel's Pheanstalk::DEFAULT_TTR constant.
		if p.ttr != 60*time.Second {
			t.Errorf("puts[%d].ttr: got %s, want 60s", i, p.ttr)
		}
	}
}

// Port of Illuminate\Tests\Queue\QueueBeanstalkdQueueTest::testPopProperlyPopsJobOffOfBeanstalkd
func TestPopProperlyPopsJobOffOfBeanstalkd(t *testing.T) {
	t.Parallel()

	drv, client := newBeanstalkdDriverForPort(0)
	client.reserveID = 42
	client.reserveBody = []byte(`{"job":"foo"}`)

	job, err := drv.Pop(context.Background(), "default")

	if err != nil {
		t.Fatalf("Pop: %v", err)
	}

	if job == nil {
		t.Fatal("Pop returned nil job")
	}

	if len(client.reserves) != 1 {
		t.Fatalf("reserves: got %d, want 1", len(client.reserves))
	}

	// Default blockFor is 0 → non-blocking reserve.
	if client.reserves[0].timeout != 0 {
		t.Errorf("reserve timeout: got %s, want 0", client.reserves[0].timeout)
	}

	if client.reserves[0].tube != "default" {
		t.Errorf("reserve tube: got %q, want default", client.reserves[0].tube)
	}
}

// Port of Illuminate\Tests\Queue\QueueBeanstalkdQueueTest::testBlockingPopProperlyPopsJobOffOfBeanstalkd
func TestBlockingPopProperlyPopsJobOffOfBeanstalkd(t *testing.T) {
	t.Parallel()

	drv, client := newBeanstalkdDriverForPort(60 * time.Second)
	client.reserveID = 7
	client.reserveBody = []byte(`{"job":"foo"}`)

	job, err := drv.Pop(context.Background(), "default")

	if err != nil {
		t.Fatalf("Pop: %v", err)
	}

	if job == nil {
		t.Fatal("Pop returned nil job")
	}

	if len(client.reserves) != 1 {
		t.Fatalf("reserves: got %d, want 1", len(client.reserves))
	}

	if client.reserves[0].timeout != 60*time.Second {
		t.Errorf("reserve timeout: got %s, want 60s", client.reserves[0].timeout)
	}
}

// Port of Illuminate\Tests\Queue\QueueBeanstalkdQueueTest::testDeleteProperlyRemoveJobsOffBeanstalkd
func TestDeleteProperlyRemoveJobsOffBeanstalkd(t *testing.T) {
	t.Parallel()

	drv, client := newBeanstalkdDriverForPort(0)
	client.reserveID = 1
	client.reserveBody = []byte(`{"job":"foo"}`)

	job, err := drv.Pop(context.Background(), "default")

	if err != nil {
		t.Fatalf("Pop: %v", err)
	}

	if job == nil {
		t.Fatal("Pop returned nil job")
	}

	if err := job.Delete(); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if len(client.deletes) != 1 {
		t.Fatalf("deletes: got %d, want 1", len(client.deletes))
	}

	if client.deletes[0] != 1 {
		t.Errorf("deletes[0]: got %d, want 1", client.deletes[0])
	}

	// Guard against accidental ErrNoJob being returned when the
	// mock has primed a reserve response.
	if errors.Is(err, queue.ErrNoJob) {
		t.Error("unexpected ErrNoJob from primed reserve")
	}
}
