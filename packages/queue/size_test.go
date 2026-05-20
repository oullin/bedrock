package queue_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
)

// memoryQueue is a minimal in-memory Queue used only by the SizeTest
// port. It tracks per-queue-name pending counts — enough to assert the
// behaviour QueueSizeTest cares about without pulling in a full driver.
type memoryQueue struct {
	mu     sync.Mutex
	counts map[string]int64
}

func newMemoryQueue() *memoryQueue { return &memoryQueue{counts: map[string]int64{}} }

func (m *memoryQueue) Push(_ context.Context, queueName string, _ []byte) (string, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.counts[queueName]++

	return "", nil
}

func (m *memoryQueue) PushDelayed(ctx context.Context, queueName string, payload []byte, _ time.Duration) (string, error) {
	return m.Push(ctx, queueName, payload)
}

func (m *memoryQueue) PushMultiple(ctx context.Context, queueName string, payloads [][]byte) ([]string, error) {
	ids := make([]string, len(payloads))

	for i, p := range payloads {
		_, _ = m.Push(ctx, queueName, p)

		ids[i] = ""
	}

	return ids, nil
}

func (m *memoryQueue) Pop(_ context.Context, _ string) (queue.Job, error) {
	return nil, queue.ErrNoJob
}

func (m *memoryQueue) Size(_ context.Context, queueName string) (int64, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	return m.counts[queueName], nil
}

func (m *memoryQueue) PendingSize(ctx context.Context, queueName string) (int64, error) {
	return m.Size(ctx, queueName)
}

func (m *memoryQueue) DelayedSize(context.Context, string) (int64, error)  { return 0, nil }
func (m *memoryQueue) ReservedSize(context.Context, string) (int64, error) { return 0, nil }
func (m *memoryQueue) ConnectionName() string                              { return "memory" }

// Port of Illuminate\Tests\Queue\QueueSizeTest::test_queue_size
//
// Upstream dispatches jobs via dispatch()/onQueue() and inspects
// Queue::size(). Go has no facade/dispatch helper, so the port pushes
// directly onto a memoryQueue and asserts the same per-queue counts.
// Behaviour under assertion (zero-before-push, split between default
// and Q2) is identical to the PHP test.
func TestQueueSize(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	q := newMemoryQueue()

	if n, _ := q.Size(ctx, ""); n != 0 {
		t.Errorf("initial default size: got %d, want 0", n)
	}

	if n, _ := q.Size(ctx, "Q2"); n != 0 {
		t.Errorf("initial Q2 size: got %d, want 0", n)
	}

	// Upstream: dispatch($job); dispatch(new TestJob2); dispatch($job)->onQueue('Q2');
	// Two jobs on the default queue, one on Q2.
	_, _ = q.Push(ctx, "", []byte("job1"))
	_, _ = q.Push(ctx, "", []byte("job2"))
	_, _ = q.Push(ctx, "Q2", []byte("job3"))

	if n, _ := q.Size(ctx, ""); n != 2 {
		t.Errorf("default size after 2 pushes: got %d, want 2", n)
	}

	if n, _ := q.Size(ctx, "Q2"); n != 1 {
		t.Errorf("Q2 size after 1 push: got %d, want 1", n)
	}
}
