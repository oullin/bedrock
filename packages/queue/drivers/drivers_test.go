package drivers_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bedrock/packages/queue/drivers"
)

// --- Mock Redis Client ---

type mockRedisClient struct {
	mu      sync.Mutex
	lists   map[string][]string
	sorted  map[string][]sortedEntry
	pushErr error
}

type sortedEntry struct {
	score  float64
	member string
}

// --- Mock DB Execer ---

type mockDBRow struct {
	values []any
	err    error
}

type mockDBExecer struct {
	mu         sync.Mutex
	rows       []*mockDBRow
	rowIndex   int
	execCalls  []mockExecCall
	execErr    error
	queryRows  []*mockDBRow // rows returned by the next Query call
	queryCalls []mockExecCall
	queryErr   error
}

type mockExecCall struct {
	Query string
	Args  []any
}

// --- Mock SQS Client ---

type mockSQSClient struct {
	mu         sync.Mutex
	messages   map[string][]mockSQSMsg
	sendErr    error
	receiveErr error
	deleteErr  error
	visErr     error
	attrErr    error
	attrs      map[string]string
	nextMsgID  int
}

type mockSQSMsg struct {
	id      string
	receipt string
	body    string
}

// --- Mock Beanstalkd Client ---

type mockBeanstalkdClient struct {
	mu         sync.Mutex
	tubes      map[string][]mockBeanstalkdJob
	nextID     uint64
	reserveErr error
	deleteErr  error
	releaseErr error
	buryErr    error
	putErr     error
	stats      map[string]map[string]string
	lastPut    struct {
		tube     string
		priority uint32
		delay    time.Duration
		ttr      time.Duration
	}
}

type mockBeanstalkdJob struct {
	id   uint64
	body []byte
}

// addQueryRow stages a single row to be returned by the next Query call.
// Call multiple times to stage a multi-row result set.

// Hand the staged rows to a fresh iterator and clear the slot.

// mockDBRows iterates a pre-staged slice of rows. It is the minimal
// drivers.DBRows implementation tests need.
type mockDBRows struct {
	rows  []*mockDBRow
	index int
	err   error
}

func newMockRedisClient() *mockRedisClient {
	return &mockRedisClient{
		lists:  make(map[string][]string),
		sorted: make(map[string][]sortedEntry),
	}
}

func (c *mockRedisClient) LPush(_ context.Context, key string, values ...any) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	if c.pushErr != nil {
		return c.pushErr
	}

	for _, v := range values {
		c.lists[key] = append([]string{fmt.Sprint(v)}, c.lists[key]...)
	}

	return nil
}

func (c *mockRedisClient) RPop(_ context.Context, key string) (string, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	list := c.lists[key]

	if len(list) == 0 {
		return "", errors.New("empty")
	}

	val := list[len(list)-1]
	c.lists[key] = list[:len(list)-1]

	return val, nil
}

func (c *mockRedisClient) ZAdd(_ context.Context, key string, score float64, member string) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.sorted[key] = append(c.sorted[key], sortedEntry{score: score, member: member})

	return nil
}

func (c *mockRedisClient) ZRangeByScore(_ context.Context, key string, min, max float64) ([]string, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	var result []string

	for _, e := range c.sorted[key] {
		if e.score >= min && e.score <= max {
			result = append(result, e.member)
		}
	}

	return result, nil
}

func (c *mockRedisClient) ZRem(_ context.Context, key string, members ...any) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	memberSet := make(map[string]bool)

	for _, m := range members {
		memberSet[fmt.Sprint(m)] = true
	}

	var remaining []sortedEntry

	for _, e := range c.sorted[key] {
		if !memberSet[e.member] {
			remaining = append(remaining, e)
		}
	}

	c.sorted[key] = remaining

	return nil
}

func (c *mockRedisClient) LLen(_ context.Context, key string) (int64, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	return int64(len(c.lists[key])), nil
}

func (c *mockRedisClient) ZCard(_ context.Context, key string) (int64, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	return int64(len(c.sorted[key])), nil
}

func (r *mockDBRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}

	for i, d := range dest {
		if i < len(r.values) {
			switch ptr := d.(type) {
			case *int64:
				switch v := r.values[i].(type) {
				case int64:
					*ptr = v
				case int:
					*ptr = int64(v)
				}
			case *int:
				switch v := r.values[i].(type) {
				case int:
					*ptr = v
				case int64:
					*ptr = int(v)
				}
			case *string:
				if v, ok := r.values[i].(string); ok {
					*ptr = v
				}
			}
		}
	}

	return nil
}

func newMockDBExecer() *mockDBExecer {
	return &mockDBExecer{}
}

func (db *mockDBExecer) addRow(values ...any) {
	db.mu.Lock()

	defer db.mu.Unlock()

	db.rows = append(db.rows, &mockDBRow{values: values})
}

func (db *mockDBExecer) addErrorRow(err error) {
	db.mu.Lock()

	defer db.mu.Unlock()

	db.rows = append(db.rows, &mockDBRow{err: err})
}

func (db *mockDBExecer) QueryRow(_ context.Context, _ string, _ ...any) drivers.DBRow {
	db.mu.Lock()

	defer db.mu.Unlock()

	if db.rowIndex >= len(db.rows) {
		return &mockDBRow{err: errors.New("no rows")}
	}

	row := db.rows[db.rowIndex]
	db.rowIndex++

	return row
}

func (db *mockDBExecer) Exec(_ context.Context, query string, args ...any) error {
	db.mu.Lock()

	defer db.mu.Unlock()

	db.execCalls = append(db.execCalls, mockExecCall{Query: query, Args: args})

	return db.execErr
}

func (db *mockDBExecer) addQueryRow(values ...any) {
	db.mu.Lock()

	defer db.mu.Unlock()

	db.queryRows = append(db.queryRows, &mockDBRow{values: values})
}

func (db *mockDBExecer) Query(_ context.Context, query string, args ...any) (drivers.DBRows, error) {
	db.mu.Lock()

	defer db.mu.Unlock()

	db.queryCalls = append(db.queryCalls, mockExecCall{Query: query, Args: args})

	if db.queryErr != nil {
		return nil, db.queryErr
	}

	iter := &mockDBRows{rows: db.queryRows}
	db.queryRows = nil

	return iter, nil
}

func (r *mockDBRows) Next() bool {
	if r.err != nil || r.index >= len(r.rows) {
		return false
	}

	return true
}

func (r *mockDBRows) Scan(dest ...any) error {
	if r.index >= len(r.rows) {
		return errors.New("no more rows")
	}

	row := r.rows[r.index]
	r.index++

	if row.err != nil {
		r.err = row.err

		return row.err
	}

	if len(dest) != len(row.values) {
		return fmt.Errorf("mockDBRows.Scan: dest len %d, values len %d", len(dest), len(row.values))
	}

	for i, v := range row.values {
		switch d := dest[i].(type) {
		case *int64:
			*d = v.(int64)
		case *int:
			*d = v.(int)
		case *string:
			*d = v.(string)
		case **int64:
			// Accept either an untyped nil or a typed *int64 value.
			if v == nil {
				*d = nil
			} else if p, ok := v.(*int64); ok {
				*d = p
			} else if x, ok := v.(int64); ok {
				*d = &x
			} else {
				return fmt.Errorf("mockDBRows.Scan: dest **int64 got unsupported value type %T", v)
			}
		default:
			return fmt.Errorf("mockDBRows.Scan: unsupported dest type %T for value %v", dest[i], v)
		}
	}

	return nil
}

func (r *mockDBRows) Close() error { return nil }

func (r *mockDBRows) Err() error { return r.err }

func newMockSQSClient() *mockSQSClient {
	return &mockSQSClient{
		messages: make(map[string][]mockSQSMsg),
		attrs:    make(map[string]string),
	}
}

func (c *mockSQSClient) SendMessage(_ context.Context, queueURL string, body string, _ time.Duration) (string, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if c.sendErr != nil {
		return "", c.sendErr
	}

	c.nextMsgID++
	id := fmt.Sprintf("msg-%d", c.nextMsgID)
	c.messages[queueURL] = append(c.messages[queueURL], mockSQSMsg{
		id:      id,
		receipt: fmt.Sprintf("receipt-%d", c.nextMsgID),
		body:    body,
	})

	return id, nil
}

func (c *mockSQSClient) SendMessageBatch(_ context.Context, queueURL string, bodies []string) ([]string, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if c.sendErr != nil {
		return nil, c.sendErr
	}

	ids := make([]string, len(bodies))

	for i, body := range bodies {
		c.nextMsgID++
		id := fmt.Sprintf("msg-%d", c.nextMsgID)
		ids[i] = id
		c.messages[queueURL] = append(c.messages[queueURL], mockSQSMsg{
			id:      id,
			receipt: fmt.Sprintf("receipt-%d", c.nextMsgID),
			body:    body,
		})
	}

	return ids, nil
}

func (c *mockSQSClient) ReceiveMessage(_ context.Context, queueURL string, maxMessages int, _ int) ([]drivers.SQSMessage, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if c.receiveErr != nil {
		return nil, c.receiveErr
	}

	msgs := c.messages[queueURL]

	if len(msgs) == 0 {
		return nil, nil
	}

	n := maxMessages

	if n > len(msgs) {
		n = len(msgs)
	}

	result := make([]drivers.SQSMessage, n)

	for i := 0; i < n; i++ {
		result[i] = drivers.SQSMessage{
			MessageID:     msgs[i].id,
			ReceiptHandle: msgs[i].receipt,
			Body:          msgs[i].body,
		}
	}

	c.messages[queueURL] = msgs[n:]

	return result, nil
}

func (c *mockSQSClient) DeleteMessage(_ context.Context, _ string, _ string) error {
	return c.deleteErr
}

func (c *mockSQSClient) ChangeMessageVisibility(_ context.Context, _ string, _ string, _ time.Duration) error {
	return c.visErr
}

func (c *mockSQSClient) GetQueueAttributes(_ context.Context, _ string, _ []string) (map[string]string, error) {
	if c.attrErr != nil {
		return nil, c.attrErr
	}

	return c.attrs, nil
}

func newMockBeanstalkdClient() *mockBeanstalkdClient {
	return &mockBeanstalkdClient{
		tubes: make(map[string][]mockBeanstalkdJob),
		stats: make(map[string]map[string]string),
	}
}

func (c *mockBeanstalkdClient) Put(_ context.Context, tube string, body []byte, priority uint32, delay, ttr time.Duration) (uint64, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if c.putErr != nil {
		return 0, c.putErr
	}

	c.nextID++
	c.tubes[tube] = append(c.tubes[tube], mockBeanstalkdJob{id: c.nextID, body: body})
	c.lastPut.tube = tube
	c.lastPut.priority = priority
	c.lastPut.delay = delay
	c.lastPut.ttr = ttr

	return c.nextID, nil
}

func (c *mockBeanstalkdClient) ReserveWithTimeout(_ context.Context, tube string, _ time.Duration) (uint64, []byte, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if c.reserveErr != nil {
		return 0, nil, c.reserveErr
	}

	jobs := c.tubes[tube]

	if len(jobs) == 0 {
		return 0, nil, errors.New("empty")
	}

	job := jobs[0]
	c.tubes[tube] = jobs[1:]

	return job.id, job.body, nil
}

func (c *mockBeanstalkdClient) Delete(_ context.Context, _ uint64) error {
	return c.deleteErr
}

func (c *mockBeanstalkdClient) Release(_ context.Context, _ uint64, _ uint32, _ time.Duration) error {
	return c.releaseErr
}

func (c *mockBeanstalkdClient) Bury(_ context.Context, _ uint64, _ uint32) error {
	return c.buryErr
}

func (c *mockBeanstalkdClient) StatsTube(_ context.Context, tube string) (map[string]string, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if s, ok := c.stats[tube]; ok {
		return s, nil
	}

	return map[string]string{}, nil
}
