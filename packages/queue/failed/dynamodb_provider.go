package failed

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"
)

// DynamoDBClient is the subset of the AWS DynamoDB API the failed-job
// provider needs. Tests inject fakes; production wires the real AWS
// client through a thin adapter.
//
// PutItem for Log, Query for All, GetItem for Find, DeleteItem for
// Forget. Each method receives the request parameters upstream emits
// verbatim so tests can assert the on-wire request shape.
type DynamoDBClient interface {
	PutItem(ctx context.Context, params map[string]any) (map[string]any, error)
	Query(ctx context.Context, params map[string]any) (map[string]any, error)
	GetItem(ctx context.Context, params map[string]any) (map[string]any, error)
	DeleteItem(ctx context.Context, params map[string]any) (map[string]any, error)
}

// Ref: @bedrock/code-0257
type DynamoDbFailedJobProvider struct {
	client          DynamoDBClient
	applicationName string
	table           string
	now             func() time.Time
}

// NewDynamoDbFailedJobProvider constructs a provider that writes to the
// given DynamoDB table.
func NewDynamoDbFailedJobProvider(client DynamoDBClient, applicationName, table string) *DynamoDbFailedJobProvider {
	return &DynamoDbFailedJobProvider{
		client:          client,
		applicationName: applicationName,
		table:           table,
		now:             time.Now,
	}
}

// Log implements FailedJobProvider.
func (p *DynamoDbFailedJobProvider) Log(connection, queue, payload string, exception error) (string, error) {
	id := extractUUID(payload)
	failedAt := p.now()
	_, err := p.client.PutItem(context.Background(), map[string]any{
		"TableName": p.table,
		"Item": map[string]any{
			"application": dynAttr("S", p.applicationName),
			"uuid":        dynAttr("S", id),
			"connection":  dynAttr("S", connection),
			"queue":       dynAttr("S", queue),
			"payload":     dynAttr("S", payload),
			"exception":   dynAttr("S", errString(exception)),
			"failed_at":   dynAttr("N", strconv.FormatInt(failedAt.Unix(), 10)),
			"expires_at":  dynAttr("N", strconv.FormatInt(failedAt.Add(7*24*time.Hour).Unix(), 10)),
		},
	})

	if err != nil {
		return "", err
	}

	return id, nil
}

// IDs implements FailedJobProvider.
func (p *DynamoDbFailedJobProvider) IDs(queueFilter string) ([]string, error) {
	all, err := p.All()

	if err != nil {
		return nil, err
	}

	out := []string{}

	for _, j := range all {
		if queueFilter != "" && j.Queue != queueFilter {
			continue
		}

		out = append(out, j.ID)
	}

	return out, nil
}

// All implements FailedJobProvider.
func (p *DynamoDbFailedJobProvider) All() ([]FailedJob, error) {
	resp, err := p.client.Query(context.Background(), map[string]any{
		"TableName":              p.table,
		"Select":                 "ALL_ATTRIBUTES",
		"KeyConditionExpression": "application = :application",
		"ExpressionAttributeValues": map[string]any{
			":application": dynAttr("S", p.applicationName),
		},
		"ScanIndexForward": false,
	})

	if err != nil {
		return nil, err
	}

	itemsAny, ok := resp["Items"]

	if !ok {
		return nil, nil
	}

	items, ok := itemsAny.([]map[string]any)

	if !ok {
		return nil, nil
	}
	// Sort by failed_at descending, matching the upstream sortByDesc.
	sort.SliceStable(items, func(i, j int) bool {
		return dynNumber(items[i], "failed_at") > dynNumber(items[j], "failed_at")
	})

	out := make([]FailedJob, 0, len(items))

	for _, it := range items {
		out = append(out, itemToFailedJob(it))
	}

	return out, nil
}

// Find implements FailedJobProvider.
func (p *DynamoDbFailedJobProvider) Find(id string) (*FailedJob, error) {
	resp, err := p.client.GetItem(context.Background(), map[string]any{
		"TableName": p.table,
		"Key": map[string]any{
			"application": dynAttr("S", p.applicationName),
			"uuid":        dynAttr("S", id),
		},
	})

	if err != nil {
		return nil, err
	}

	itemAny, ok := resp["Item"]

	if !ok {
		return nil, nil
	}

	item, ok := itemAny.(map[string]any)

	if !ok {
		return nil, nil
	}

	fj := itemToFailedJob(item)

	return &fj, nil
}

// Forget implements FailedJobProvider.
func (p *DynamoDbFailedJobProvider) Forget(id string) (bool, error) {
	_, err := p.client.DeleteItem(context.Background(), map[string]any{
		"TableName": p.table,
		"Key": map[string]any{
			"application": dynAttr("S", p.applicationName),
			"uuid":        dynAttr("S", id),
		},
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

// Flush implements FailedJobProvider. DynamoDB storage relies on the
// table's TTL feature; upstream throws an exception here and this port
// returns a matching sentinel error.
func (p *DynamoDbFailedJobProvider) Flush(_ int) error {
	return ErrDynamoDbFlushUnsupported
}

// from DynamoDbFailedJobProvider::flush.
var ErrDynamoDbFlushUnsupported = errors.New("dynamodb failed job storage may not be flushed. use the table's TTL feature on expires_at")

// dynAttr renders a DynamoDB attribute map `{"S":"value"}`.
func dynAttr(kind, value string) map[string]any {
	return map[string]any{kind: value}
}

// dynNumber extracts a numeric attribute from an item map.
func dynNumber(item map[string]any, field string) int64 {
	v, ok := item[field].(map[string]any)

	if !ok {
		return 0
	}

	s, ok := v["N"].(string)

	if !ok {
		return 0
	}

	n, _ := strconv.ParseInt(s, 10, 64)

	return n
}

func dynString(item map[string]any, field string) string {
	v, ok := item[field].(map[string]any)

	if !ok {
		return ""
	}

	s, _ := v["S"].(string)

	return s
}

func itemToFailedJob(it map[string]any) FailedJob {
	uuid := dynString(it, "uuid")

	return FailedJob{
		ID:         uuid,
		UUID:       uuid,
		Connection: dynString(it, "connection"),
		Queue:      dynString(it, "queue"),
		Payload:    dynString(it, "payload"),
		Exception:  dynString(it, "exception"),
		FailedAt:   time.Unix(dynNumber(it, "failed_at"), 0),
	}
}
