package bus_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/bus"
)

// mockDynamoClient implements bus.DynamoClient for testing.
type mockDynamoClient struct {
	mu      sync.Mutex
	items   map[string]map[string]any // keyed by "app:id"
	putErr  error
	getErr  error
	updErr  error
	delErr  error
}

func newMockDynamoClient() *mockDynamoClient {
	return &mockDynamoClient{items: make(map[string]map[string]any)}
}

func (c *mockDynamoClient) itemKey(key map[string]any) string {
	return fmt.Sprintf("%v:%v", key["application"], key["id"])
}

func (c *mockDynamoClient) PutItem(_ context.Context, _ string, item map[string]any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.putErr != nil {
		return c.putErr
	}

	k := fmt.Sprintf("%v:%v", item["application"], item["id"])
	c.items[k] = item

	return nil
}

func (c *mockDynamoClient) GetItem(_ context.Context, _ string, key map[string]any) (map[string]any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.getErr != nil {
		return nil, c.getErr
	}

	return c.items[c.itemKey(key)], nil
}

func (c *mockDynamoClient) UpdateItem(_ context.Context, _ string, key map[string]any, _ string, values map[string]any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.updErr != nil {
		return c.updErr
	}

	k := c.itemKey(key)
	item, ok := c.items[k]
	if !ok {
		return nil
	}

	// Apply simple value updates for testing.
	for vk, vv := range values {
		item[vk] = vv
	}

	return nil
}

func (c *mockDynamoClient) DeleteItem(_ context.Context, _ string, key map[string]any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.delErr != nil {
		return c.delErr
	}

	delete(c.items, c.itemKey(key))

	return nil
}

func (c *mockDynamoClient) Query(_ context.Context, _ string, _ string, _ map[string]any, _ int) ([]map[string]any, error) {
	return nil, nil
}

func TestDynamoRepoStore(t *testing.T) {
	client := newMockDynamoClient()
	repo := bus.NewDynamoBatchRepository(client, "myapp", "batches", nil, "")

	batch := &bus.Batch{
		ID:           "batch-1",
		Name:         "dynamo-test",
		TotalJobs:    5,
		PendingJobs:  5,
		FailedJobIDs: []string{"j1"},
		Options:      map[string]any{"allowFailures": true},
		CreatedAt:    time.Now(),
	}

	err := repo.Store(context.Background(), batch)
	if err != nil {
		t.Fatal(err)
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	item, ok := client.items["myapp:batch-1"]
	if !ok {
		t.Fatal("expected item to be stored")
	}

	if item["name"] != "dynamo-test" {
		t.Errorf("expected name 'dynamo-test', got %v", item["name"])
	}
}

func TestDynamoRepoGet(t *testing.T) {
	client := newMockDynamoClient()
	repo := bus.NewDynamoBatchRepository(client, "myapp", "batches", nil, "")

	failedJSON, _ := json.Marshal([]string{})
	optsJSON, _ := json.Marshal(map[string]any{})

	client.items["myapp:batch-1"] = map[string]any{
		"id":             "batch-1",
		"name":           "test",
		"total_jobs":     float64(10),
		"pending_jobs":   float64(3),
		"failed_jobs":    float64(1),
		"failed_job_ids": string(failedJSON),
		"options":        string(optsJSON),
		"created_at":     float64(time.Now().Unix()),
	}

	batch, err := repo.Get(context.Background(), "batch-1")
	if err != nil {
		t.Fatal(err)
	}

	if batch == nil {
		t.Fatal("expected non-nil batch")
	}

	if batch.ID != "batch-1" {
		t.Errorf("expected ID 'batch-1', got %q", batch.ID)
	}

	if batch.TotalJobs != 10 {
		t.Errorf("expected TotalJobs 10, got %d", batch.TotalJobs)
	}
}

func TestDynamoRepoGetNotFound(t *testing.T) {
	client := newMockDynamoClient()
	repo := bus.NewDynamoBatchRepository(client, "myapp", "batches", nil, "")

	batch, err := repo.Get(context.Background(), "nonexistent")
	if err != nil {
		t.Fatal(err)
	}

	if batch != nil {
		t.Error("expected nil batch for nonexistent ID")
	}
}

func TestDynamoRepoMarkAsFinished(t *testing.T) {
	client := newMockDynamoClient()
	client.items["myapp:batch-1"] = map[string]any{
		"application": "myapp",
		"id":          "batch-1",
	}

	repo := bus.NewDynamoBatchRepository(client, "myapp", "batches", nil, "")

	err := repo.MarkAsFinished(context.Background(), "batch-1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDynamoRepoCancel(t *testing.T) {
	client := newMockDynamoClient()
	client.items["myapp:batch-1"] = map[string]any{
		"application": "myapp",
		"id":          "batch-1",
	}

	repo := bus.NewDynamoBatchRepository(client, "myapp", "batches", nil, "")

	err := repo.Cancel(context.Background(), "batch-1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDynamoRepoDelete(t *testing.T) {
	client := newMockDynamoClient()
	client.items["myapp:batch-1"] = map[string]any{
		"application": "myapp",
		"id":          "batch-1",
	}

	repo := bus.NewDynamoBatchRepository(client, "myapp", "batches", nil, "")

	err := repo.Delete(context.Background(), "batch-1")
	if err != nil {
		t.Fatal(err)
	}

	client.mu.Lock()
	_, exists := client.items["myapp:batch-1"]
	client.mu.Unlock()

	if exists {
		t.Error("expected item to be deleted")
	}
}

func TestDynamoRepoStoreWithTTL(t *testing.T) {
	ttl := 86400
	client := newMockDynamoClient()
	repo := bus.NewDynamoBatchRepository(client, "myapp", "batches", &ttl, "expires_at")

	batch := &bus.Batch{
		ID:        "batch-ttl",
		Options:   map[string]any{},
		CreatedAt: time.Now(),
	}

	err := repo.Store(context.Background(), batch)
	if err != nil {
		t.Fatal(err)
	}

	client.mu.Lock()
	item := client.items["myapp:batch-ttl"]
	client.mu.Unlock()

	if item == nil {
		t.Fatal("expected item to be stored")
	}

	if _, ok := item["expires_at"]; !ok {
		t.Error("expected expires_at field to be set with TTL")
	}
}

func TestDynamoRepoTransaction(t *testing.T) {
	client := newMockDynamoClient()
	repo := bus.NewDynamoBatchRepository(client, "myapp", "batches", nil, "")

	called := false
	err := repo.Transaction(context.Background(), func(r bus.BatchRepository) error {
		called = true

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error("expected fn to be called")
	}
}

func TestDynamoRepoStoreError(t *testing.T) {
	client := newMockDynamoClient()
	client.putErr = fmt.Errorf("dynamo error")

	repo := bus.NewDynamoBatchRepository(client, "myapp", "batches", nil, "")

	batch := &bus.Batch{ID: "err", Options: map[string]any{}, CreatedAt: time.Now()}
	err := repo.Store(context.Background(), batch)
	if err == nil {
		t.Error("expected error")
	}
}

func TestDynamoRepoIncrementTotalJobs(t *testing.T) {
	client := newMockDynamoClient()
	client.items["myapp:batch-1"] = map[string]any{
		"application": "myapp",
		"id":          "batch-1",
	}

	repo := bus.NewDynamoBatchRepository(client, "myapp", "batches", nil, "")

	err := repo.IncrementTotalJobs(context.Background(), "batch-1", 5)
	if err != nil {
		t.Fatal(err)
	}
}

// Ensure DynamoBatchRepository satisfies BatchRepository.
var _ bus.BatchRepository = (*bus.DynamoBatchRepository)(nil)
