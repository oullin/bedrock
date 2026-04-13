package bus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// DynamoClient abstracts AWS DynamoDB operations for testability.
type DynamoClient interface {
	PutItem(ctx context.Context, table string, item map[string]any) error
	GetItem(ctx context.Context, table string, key map[string]any) (map[string]any, error)
	UpdateItem(ctx context.Context, table string, key map[string]any, update string, values map[string]any) error
	DeleteItem(ctx context.Context, table string, key map[string]any) error
	Query(ctx context.Context, table string, keyCondition string, values map[string]any, limit int) ([]map[string]any, error)
}

// DynamoBatchRepository persists batch state in AWS DynamoDB.
type DynamoBatchRepository struct {
	client          DynamoClient
	applicationName string
	table           string
	ttl             *int
	ttlAttribute    string
}

// NewDynamoBatchRepository creates a DynamoBatchRepository.
func NewDynamoBatchRepository(client DynamoClient, applicationName, table string, ttl *int, ttlAttribute string) *DynamoBatchRepository {
	if table == "" {
		table = "job_batches"
	}

	if ttlAttribute == "" {
		ttlAttribute = "expires_at"
	}

	return &DynamoBatchRepository{
		client:          client,
		applicationName: applicationName,
		table:           table,
		ttl:             ttl,
		ttlAttribute:    ttlAttribute,
	}
}

// GetList retrieves batches using a query, limited by count and optionally before a given ID.
func (r *DynamoBatchRepository) GetList(ctx context.Context, limit int, before string) ([]*Batch, error) {
	keyCondition := "application = :app"
	values := map[string]any{":app": r.applicationName}

	if before != "" {
		keyCondition += " AND id < :before"
		values[":before"] = before
	}

	items, err := r.client.Query(ctx, r.table, keyCondition, values, limit)

	if err != nil {
		return nil, fmt.Errorf("bus: dynamo get list: %w", err)
	}

	batches := make([]*Batch, 0, len(items))

	for _, item := range items {
		batch, err := r.itemToBatch(item)

		if err != nil {
			return nil, err
		}

		batches = append(batches, batch)
	}

	return batches, nil
}

// Get retrieves a batch by ID.
func (r *DynamoBatchRepository) Get(ctx context.Context, id string) (*Batch, error) {
	key := map[string]any{
		"application": r.applicationName,
		"id":          id,
	}

	item, err := r.client.GetItem(ctx, r.table, key)

	if err != nil {
		return nil, fmt.Errorf("bus: dynamo get batch: %w", err)
	}

	if item == nil {
		return nil, nil
	}

	return r.itemToBatch(item)
}

// Store persists a new batch.
func (r *DynamoBatchRepository) Store(ctx context.Context, batch *Batch) error {
	failedIDs, err := json.Marshal(batch.FailedJobIDs)

	if err != nil {
		return fmt.Errorf("bus: marshal failed_job_ids: %w", err)
	}

	opts, err := json.Marshal(batch.Options)

	if err != nil {
		return fmt.Errorf("bus: marshal options: %w", err)
	}

	item := map[string]any{
		"application":    r.applicationName,
		"id":             batch.ID,
		"name":           batch.Name,
		"total_jobs":     batch.TotalJobs,
		"pending_jobs":   batch.PendingJobs,
		"failed_jobs":    batch.FailedJobs,
		"failed_job_ids": string(failedIDs),
		"options":        string(opts),
		"created_at":     batch.CreatedAt.Unix(),
	}

	if batch.CancelledAt != nil {
		item["cancelled_at"] = batch.CancelledAt.Unix()
	}

	if batch.FinishedAt != nil {
		item["finished_at"] = batch.FinishedAt.Unix()
	}

	if r.ttl != nil {
		item[r.ttlAttribute] = time.Now().Add(time.Duration(*r.ttl) * time.Second).Unix()
	}

	if err = r.client.PutItem(ctx, r.table, item); err != nil {
		return fmt.Errorf("bus: dynamo store batch: %w", err)
	}

	return nil
}

// IncrementTotalJobs increments the total job count.
func (r *DynamoBatchRepository) IncrementTotalJobs(ctx context.Context, id string, amount int) error {
	key := map[string]any{
		"application": r.applicationName,
		"id":          id,
	}

	err := r.client.UpdateItem(ctx, r.table, key,
		"SET total_jobs = total_jobs + :amount, pending_jobs = pending_jobs + :amount",
		map[string]any{":amount": amount},
	)

	if err != nil {
		return fmt.Errorf("bus: dynamo increment total jobs: %w", err)
	}

	return nil
}

// DecrementPendingJobs decrements the pending job count and returns updated counts.
func (r *DynamoBatchRepository) DecrementPendingJobs(ctx context.Context, id string) (*UpdatedBatchJobCounts, error) {
	key := map[string]any{
		"application": r.applicationName,
		"id":          id,
	}

	err := r.client.UpdateItem(ctx, r.table, key,
		"SET pending_jobs = pending_jobs - :one",
		map[string]any{":one": 1},
	)

	if err != nil {
		return nil, fmt.Errorf("bus: dynamo decrement pending jobs: %w", err)
	}

	return r.fetchCounts(ctx, id)
}

// IncrementFailedJobs increments the failed job count and appends the job ID.
func (r *DynamoBatchRepository) IncrementFailedJobs(ctx context.Context, id string, failedJobID string) (*UpdatedBatchJobCounts, error) {
	// First fetch current failed_job_ids.
	batch, err := r.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	if batch == nil {
		return nil, fmt.Errorf("bus: batch %s not found", id)
	}

	failedIDs := append(batch.FailedJobIDs, failedJobID)
	failedJSON, err := json.Marshal(failedIDs)

	if err != nil {
		return nil, fmt.Errorf("bus: marshal failed_job_ids: %w", err)
	}

	key := map[string]any{
		"application": r.applicationName,
		"id":          id,
	}

	err = r.client.UpdateItem(ctx, r.table, key,
		"SET failed_jobs = failed_jobs + :one, pending_jobs = pending_jobs - :one, failed_job_ids = :ids",
		map[string]any{":one": 1, ":ids": string(failedJSON)},
	)

	if err != nil {
		return nil, fmt.Errorf("bus: dynamo increment failed jobs: %w", err)
	}

	return r.fetchCounts(ctx, id)
}

// MarkAsFinished sets the finished_at timestamp.
func (r *DynamoBatchRepository) MarkAsFinished(ctx context.Context, id string) error {
	key := map[string]any{
		"application": r.applicationName,
		"id":          id,
	}

	err := r.client.UpdateItem(ctx, r.table, key,
		"SET finished_at = :ts",
		map[string]any{":ts": time.Now().Unix()},
	)

	if err != nil {
		return fmt.Errorf("bus: dynamo mark as finished: %w", err)
	}

	return nil
}

// Cancel sets the cancelled_at timestamp.
func (r *DynamoBatchRepository) Cancel(ctx context.Context, id string) error {
	key := map[string]any{
		"application": r.applicationName,
		"id":          id,
	}

	err := r.client.UpdateItem(ctx, r.table, key,
		"SET cancelled_at = :ts",
		map[string]any{":ts": time.Now().Unix()},
	)

	if err != nil {
		return fmt.Errorf("bus: dynamo cancel batch: %w", err)
	}

	return nil
}

// Delete removes a batch.
func (r *DynamoBatchRepository) Delete(ctx context.Context, id string) error {
	key := map[string]any{
		"application": r.applicationName,
		"id":          id,
	}

	if err := r.client.DeleteItem(ctx, r.table, key); err != nil {
		return fmt.Errorf("bus: dynamo delete batch: %w", err)
	}

	return nil
}

// Transaction executes fn. DynamoDB does not support transactions in the same
// way as SQL, so this simply invokes fn directly.
func (r *DynamoBatchRepository) Transaction(_ context.Context, fn func(BatchRepository) error) error {
	return fn(r)
}

// RollBack is a no-op for DynamoDB as it does not support transactions.
func (r *DynamoBatchRepository) RollBack(_ context.Context) error {
	return nil
}

func (r *DynamoBatchRepository) fetchCounts(ctx context.Context, id string) (*UpdatedBatchJobCounts, error) {
	batch, err := r.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	if batch == nil {
		return nil, fmt.Errorf("bus: batch %s not found", id)
	}

	return &UpdatedBatchJobCounts{
		PendingJobs: batch.PendingJobs,
		FailedJobs:  batch.FailedJobs,
	}, nil
}

func (r *DynamoBatchRepository) itemToBatch(item map[string]any) (*Batch, error) {
	b := &Batch{
		Options: make(map[string]any),
	}

	if v, ok := item["id"].(string); ok {
		b.ID = v
	}

	if v, ok := item["name"].(string); ok {
		b.Name = v
	}

	if v, ok := item["total_jobs"].(float64); ok {
		b.TotalJobs = int(v)
	}

	if v, ok := item["pending_jobs"].(float64); ok {
		b.PendingJobs = int(v)
	}

	if v, ok := item["failed_jobs"].(float64); ok {
		b.FailedJobs = int(v)
	}

	if v, ok := item["failed_job_ids"].(string); ok {
		if err := json.Unmarshal([]byte(v), &b.FailedJobIDs); err != nil {
			b.FailedJobIDs = []string{}
		}
	}

	if v, ok := item["options"].(string); ok {
		if err := json.Unmarshal([]byte(v), &b.Options); err != nil {
			b.Options = make(map[string]any)
		}
	}

	if v, ok := item["created_at"].(float64); ok {
		b.CreatedAt = time.Unix(int64(v), 0)
	}

	if v, ok := item["cancelled_at"].(float64); ok {
		t := time.Unix(int64(v), 0)
		b.CancelledAt = &t
	}

	if v, ok := item["finished_at"].(float64); ok {
		t := time.Unix(int64(v), 0)
		b.FinishedAt = &t
	}

	return b, nil
}
