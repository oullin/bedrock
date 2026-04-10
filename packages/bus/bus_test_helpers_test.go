package bus_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bedrock/packages/bus"
	"github.com/bedrock/packages/queue"
)

// mockQueue implements queue.Queue for testing.
type mockQueue struct {
	mu       sync.Mutex
	pushes   []mockPush
	pushErr  error
	connName string
}

type mockPush struct {
	Queue   string
	Payload []byte
}

// mockBatchRepository implements bus.BatchRepository for testing.
type mockBatchRepository struct {
	mu     sync.Mutex
	calls  []string
	batch  *bus.Batch
	getErr error

	storeErr              error
	cancelErr             error
	deleteErr             error
	decrementResult       *bus.UpdatedBatchJobCounts
	decrementErr          error
	incrementFailedResult *bus.UpdatedBatchJobCounts
	incrementFailedErr    error
}

// mockCacheStore implements bus.CacheStore for testing.
type mockCacheStore struct {
	mu    sync.Mutex
	data  map[string]string
	calls []mockCacheCall
}

type mockCacheCall struct {
	Method string
	Key    string
	Value  string
	TTL    int
}

// mockQueueingDispatcher implements bus.QueueingDispatcher for testing PendingBatch.
type mockQueueingDispatcher struct {
	mu              sync.Mutex
	dispatchedQueue []any
	dispatchErr     error
	dispatchErrAt   int // fail at this index (0-based)
	batchRepo       bus.BatchRepository
}

var errTestFailure = fmt.Errorf("test failure")

func newMockQueue() *mockQueue {
	return &mockQueue{connName: "mock"}
}

func (q *mockQueue) Push(_ context.Context, queueName string, payload []byte) (string, error) {
	q.mu.Lock()

	defer q.mu.Unlock()

	if q.pushErr != nil {
		return "", q.pushErr
	}

	q.pushes = append(q.pushes, mockPush{Queue: queueName, Payload: payload})

	return fmt.Sprintf("job-%d", len(q.pushes)), nil
}

func (q *mockQueue) PushDelayed(_ context.Context, queueName string, payload []byte, _ time.Duration) (string, error) {
	return q.Push(context.Background(), queueName, payload)
}

func (q *mockQueue) PushMultiple(_ context.Context, queueName string, payloads [][]byte) ([]string, error) {
	ids := make([]string, 0, len(payloads))

	for _, p := range payloads {
		id, err := q.Push(context.Background(), queueName, p)

		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (q *mockQueue) Pop(_ context.Context, _ string) (queue.Job, error) {
	return nil, fmt.Errorf("not implemented")
}

func (q *mockQueue) Size(_ context.Context, _ string) (int64, error)         { return 0, nil }
func (q *mockQueue) PendingSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (q *mockQueue) DelayedSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (q *mockQueue) ReservedSize(_ context.Context, _ string) (int64, error) { return 0, nil }
func (q *mockQueue) ConnectionName() string                                  { return q.connName }

func newMockBatchRepo() *mockBatchRepository {
	return &mockBatchRepository{}
}

func (r *mockBatchRepository) Get(_ context.Context, id string) (*bus.Batch, error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.calls = append(r.calls, "Get:"+id)

	return r.batch, r.getErr
}

func (r *mockBatchRepository) Store(_ context.Context, b *bus.Batch) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.calls = append(r.calls, "Store:"+b.ID)
	r.batch = b

	return r.storeErr
}

func (r *mockBatchRepository) IncrementTotalJobs(_ context.Context, id string, amount int) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.calls = append(r.calls, fmt.Sprintf("IncrementTotalJobs:%s:%d", id, amount))

	return nil
}

func (r *mockBatchRepository) DecrementPendingJobs(_ context.Context, id string) (*bus.UpdatedBatchJobCounts, error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.calls = append(r.calls, "DecrementPendingJobs:"+id)

	if r.decrementResult != nil {
		return r.decrementResult, r.decrementErr
	}

	return &bus.UpdatedBatchJobCounts{PendingJobs: 0, FailedJobs: 0}, r.decrementErr
}

func (r *mockBatchRepository) IncrementFailedJobs(_ context.Context, id string, failedJobID string) (*bus.UpdatedBatchJobCounts, error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.calls = append(r.calls, "IncrementFailedJobs:"+id+":"+failedJobID)

	if r.incrementFailedResult != nil {
		return r.incrementFailedResult, r.incrementFailedErr
	}

	return &bus.UpdatedBatchJobCounts{PendingJobs: 0, FailedJobs: 1}, r.incrementFailedErr
}

func (r *mockBatchRepository) MarkAsFinished(_ context.Context, id string) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.calls = append(r.calls, "MarkAsFinished:"+id)

	return nil
}

func (r *mockBatchRepository) Cancel(_ context.Context, id string) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.calls = append(r.calls, "Cancel:"+id)

	return r.cancelErr
}

func (r *mockBatchRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.calls = append(r.calls, "Delete:"+id)

	return r.deleteErr
}

func (r *mockBatchRepository) Transaction(_ context.Context, fn func(bus.BatchRepository) error) error {
	return fn(r)
}

func (r *mockBatchRepository) hasCalled(method string) bool {
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, c := range r.calls {
		if c == method || len(c) > len(method) && c[:len(method)] == method {
			return true
		}
	}

	return false
}

func newMockCacheStore() *mockCacheStore {
	return &mockCacheStore{data: make(map[string]string)}
}

func (c *mockCacheStore) Get(_ context.Context, key string) (string, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.calls = append(c.calls, mockCacheCall{Method: "Get", Key: key})

	return c.data[key], nil
}

func (c *mockCacheStore) Put(_ context.Context, key, value string, ttlSeconds int) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.calls = append(c.calls, mockCacheCall{Method: "Put", Key: key, Value: value, TTL: ttlSeconds})
	c.data[key] = value

	return nil
}

func (c *mockCacheStore) Forget(_ context.Context, key string) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.calls = append(c.calls, mockCacheCall{Method: "Forget", Key: key})
	delete(c.data, key)

	return nil
}

func newMockQueueingDispatcher() *mockQueueingDispatcher {
	return &mockQueueingDispatcher{dispatchErrAt: -1}
}

func (d *mockQueueingDispatcher) Dispatch(_ context.Context, command any) (any, error) {
	return nil, nil
}

func (d *mockQueueingDispatcher) DispatchSync(_ context.Context, command any) (any, error) {
	return nil, nil
}

func (d *mockQueueingDispatcher) DispatchNow(_ context.Context, command any) (any, error) {
	return nil, nil
}

func (d *mockQueueingDispatcher) DispatchAfterResponse(_ context.Context, command any) error {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.dispatchedQueue = append(d.dispatchedQueue, command)

	return nil
}

func (d *mockQueueingDispatcher) PipeThrough(pipes ...bus.Pipe) bus.Dispatcher {
	return d
}

func (d *mockQueueingDispatcher) Map(command any, handler bus.Handler) bus.Dispatcher {
	return d
}

func (d *mockQueueingDispatcher) HasCommandHandler(_ any) bool {
	return false
}

func (d *mockQueueingDispatcher) GetCommandHandler(_ any) (bus.Handler, bool) {
	return nil, false
}

func (d *mockQueueingDispatcher) Chain(jobs []any) *bus.PendingChain {
	return bus.NewPendingChain(d, jobs)
}

func (d *mockQueueingDispatcher) DispatchToQueue(_ context.Context, command any) error {
	d.mu.Lock()

	defer d.mu.Unlock()

	idx := len(d.dispatchedQueue)
	d.dispatchedQueue = append(d.dispatchedQueue, command)

	if d.dispatchErrAt >= 0 && idx == d.dispatchErrAt {
		return d.dispatchErr
	}

	return nil
}

func (d *mockQueueingDispatcher) FindBatch(_ context.Context, id string) (*bus.Batch, error) {
	return nil, nil
}

func (d *mockQueueingDispatcher) Batch(jobs []any) *bus.PendingBatch {
	return bus.NewPendingBatch(d, jobs)
}
