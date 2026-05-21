package failed_test

// Ported from @bedrock\Tests\Queue\DynamoDbFailedJobProviderTest (5/5).
//
// ✅ testCanProperlyLogFailedJob
// ✅ testCanRetrieveAllFailedJobs
// ✅ testASingleJobCanBeFound
// ✅ testNullIsReturnedIfJobNotFound
// ✅ testJobsCanBeDeleted

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/bedrock/packages/queue/failed"
)

// fakeDynamoClient records the last arguments to each call and returns
// canned responses. It is the Go analogue of the Mockery mocks used in
// DynamoDbFailedJobProviderTest.
type fakeDynamoClient struct {
	putArgs    map[string]any
	getArgs    map[string]any
	deleteArgs map[string]any
	queryArgs  map[string]any

	getResp   map[string]any
	queryResp map[string]any
}

func (f *fakeDynamoClient) PutItem(_ context.Context, params map[string]any) (map[string]any, error) {
	f.putArgs = params

	return map[string]any{}, nil
}
func (f *fakeDynamoClient) Query(_ context.Context, params map[string]any) (map[string]any, error) {
	f.queryArgs = params

	return f.queryResp, nil
}
func (f *fakeDynamoClient) GetItem(_ context.Context, params map[string]any) (map[string]any, error) {
	f.getArgs = params

	return f.getResp, nil
}
func (f *fakeDynamoClient) DeleteItem(_ context.Context, params map[string]any) (map[string]any, error) {
	f.deleteArgs = params

	return map[string]any{}, nil
}

// Port of @bedrock\Tests\Queue\DynamoDbFailedJobProviderTest::testCanProperlyLogFailedJob
func TestDynamoCanProperlyLogFailedJob(t *testing.T) {
	t.Parallel()
	fake := &fakeDynamoClient{}
	p := failed.NewDynamoDbFailedJobProvider(fake, "application", "table")

	now := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	failed.SetNow(p, func() time.Time { return now })

	uuid := "the-uuid"
	payload, _ := json.Marshal(map[string]string{"uuid": uuid})
	exception := errors.New("Something went wrong.")

	if _, err := p.Log("connection", "queue", string(payload), exception); err != nil {
		t.Fatal(err)
	}

	wantItem := map[string]any{
		"application": map[string]any{"S": "application"},
		"uuid":        map[string]any{"S": uuid},
		"connection":  map[string]any{"S": "connection"},
		"queue":       map[string]any{"S": "queue"},
		"payload":     map[string]any{"S": string(payload)},
		"exception":   map[string]any{"S": exception.Error()},
		"failed_at":   map[string]any{"N": strconv.FormatInt(now.Unix(), 10)},
		"expires_at":  map[string]any{"N": strconv.FormatInt(now.Add(7*24*time.Hour).Unix(), 10)},
	}
	want := map[string]any{"TableName": "table", "Item": wantItem}

	if !reflect.DeepEqual(fake.putArgs, want) {
		t.Fatalf("putArgs mismatch\n got: %#v\nwant: %#v", fake.putArgs, want)
	}
}

// Port of @bedrock\Tests\Queue\DynamoDbFailedJobProviderTest::testCanRetrieveAllFailedJobs
func TestDynamoCanRetrieveAllFailedJobs(t *testing.T) {
	t.Parallel()
	timeNow := time.Now().Unix()
	fake := &fakeDynamoClient{
		queryResp: map[string]any{
			"Items": []map[string]any{
				{
					"application": map[string]any{"S": "application"},
					"uuid":        map[string]any{"S": "uuid"},
					"connection":  map[string]any{"S": "connection"},
					"queue":       map[string]any{"S": "queue"},
					"payload":     map[string]any{"S": "payload"},
					"exception":   map[string]any{"S": "exception"},
					"failed_at":   map[string]any{"N": strconv.FormatInt(timeNow, 10)},
					"expires_at":  map[string]any{"N": strconv.FormatInt(timeNow, 10)},
				},
			},
		},
	}
	p := failed.NewDynamoDbFailedJobProvider(fake, "application", "table")

	all, err := p.All()

	if err != nil {
		t.Fatal(err)
	}

	if len(all) != 1 {
		t.Fatalf("count=%d", len(all))
	}

	wantQueryArgs := map[string]any{
		"TableName":              "table",
		"Select":                 "ALL_ATTRIBUTES",
		"KeyConditionExpression": "application = :application",
		"ExpressionAttributeValues": map[string]any{
			":application": map[string]any{"S": "application"},
		},
		"ScanIndexForward": false,
	}

	if !reflect.DeepEqual(fake.queryArgs, wantQueryArgs) {
		t.Fatalf("queryArgs mismatch\n got: %#v\nwant: %#v", fake.queryArgs, wantQueryArgs)
	}

	got := all[0]

	if got.ID != "uuid" || got.Connection != "connection" || got.Queue != "queue" ||
		got.Payload != "payload" || got.Exception != "exception" {
		t.Fatalf("unexpected: %+v", got)
	}

	if got.FailedAt.Unix() != timeNow {
		t.Fatalf("failed_at=%d want %d", got.FailedAt.Unix(), timeNow)
	}
}

// Port of @bedrock\Tests\Queue\DynamoDbFailedJobProviderTest::testASingleJobCanBeFound
func TestDynamoASingleJobCanBeFound(t *testing.T) {
	t.Parallel()
	timeNow := time.Now().Unix()
	fake := &fakeDynamoClient{
		getResp: map[string]any{
			"Item": map[string]any{
				"application": map[string]any{"S": "application"},
				"uuid":        map[string]any{"S": "uuid"},
				"connection":  map[string]any{"S": "connection"},
				"queue":       map[string]any{"S": "queue"},
				"payload":     map[string]any{"S": "payload"},
				"exception":   map[string]any{"S": "exception"},
				"failed_at":   map[string]any{"N": strconv.FormatInt(timeNow, 10)},
				"expires_at":  map[string]any{"N": strconv.FormatInt(timeNow, 10)},
			},
		},
	}
	p := failed.NewDynamoDbFailedJobProvider(fake, "application", "table")

	got, _ := p.Find("id")

	if got == nil {
		t.Fatal("expected job")
	}

	wantGetArgs := map[string]any{
		"TableName": "table",
		"Key": map[string]any{
			"application": map[string]any{"S": "application"},
			"uuid":        map[string]any{"S": "id"},
		},
	}

	if !reflect.DeepEqual(fake.getArgs, wantGetArgs) {
		t.Fatalf("getArgs mismatch\n got: %#v\nwant: %#v", fake.getArgs, wantGetArgs)
	}

	if got.ID != "uuid" || got.FailedAt.Unix() != timeNow {
		t.Fatalf("unexpected: %+v", got)
	}
}

// Port of @bedrock\Tests\Queue\DynamoDbFailedJobProviderTest::testNullIsReturnedIfJobNotFound
func TestDynamoNullIsReturnedIfJobNotFound(t *testing.T) {
	t.Parallel()
	fake := &fakeDynamoClient{getResp: map[string]any{}}
	p := failed.NewDynamoDbFailedJobProvider(fake, "application", "table")

	got, _ := p.Find("id")

	if got != nil {
		t.Fatalf("expected nil, got %+v", got)
	}
}

// Port of @bedrock\Tests\Queue\DynamoDbFailedJobProviderTest::testJobsCanBeDeleted
func TestDynamoJobsCanBeDeleted(t *testing.T) {
	t.Parallel()
	fake := &fakeDynamoClient{}
	p := failed.NewDynamoDbFailedJobProvider(fake, "application", "table")

	if _, err := p.Forget("id"); err != nil {
		t.Fatal(err)
	}

	wantDeleteArgs := map[string]any{
		"TableName": "table",
		"Key": map[string]any{
			"application": map[string]any{"S": "application"},
			"uuid":        map[string]any{"S": "id"},
		},
	}

	if !reflect.DeepEqual(fake.deleteArgs, wantDeleteArgs) {
		t.Fatalf("deleteArgs mismatch\n got: %#v\nwant: %#v", fake.deleteArgs, wantDeleteArgs)
	}
}
