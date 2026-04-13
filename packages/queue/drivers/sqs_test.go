package drivers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

func TestSQSDriverPush(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	drv := drivers.NewSQSDriver(client, map[string]string{"default": "https://sqs/default"}, "sqs")

	id, err := drv.Push(context.Background(), "default", []byte("payload"))

	if err != nil {
		t.Fatal(err)
	}

	if id == "" {
		t.Error("expected non-empty message id")
	}
}

func TestSQSDriverPushDelayed(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	drv := drivers.NewSQSDriver(client, map[string]string{"default": "https://sqs/default"}, "sqs")

	id, err := drv.PushDelayed(context.Background(), "default", []byte("payload"), 30*time.Second)

	if err != nil {
		t.Fatal(err)
	}

	if id == "" {
		t.Error("expected non-empty message id")
	}
}

func TestSQSDriverPushMultiple(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	drv := drivers.NewSQSDriver(client, map[string]string{"default": "https://sqs/default"}, "sqs")

	ids, err := drv.PushMultiple(context.Background(), "default", [][]byte{
		[]byte("a"), []byte("b"), []byte("c"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 3 {
		t.Errorf("expected 3 ids, got %d", len(ids))
	}
}

func TestSQSDriverPop(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	urls := map[string]string{"default": "https://sqs/default"}
	drv := drivers.NewSQSDriver(client, urls, "sqs")

	_, _ = drv.Push(context.Background(), "default", []byte("test-body"))

	job, err := drv.Pop(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if string(job.Payload()) != "test-body" {
		t.Errorf("expected 'test-body', got %q", job.Payload())
	}

	if job.GetQueue() != "default" {
		t.Errorf("expected queue 'default', got %q", job.GetQueue())
	}

	if job.GetConnectionName() != "sqs" {
		t.Errorf("expected connection 'sqs', got %q", job.GetConnectionName())
	}
}

func TestSQSDriverPopEmpty(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	drv := drivers.NewSQSDriver(client, map[string]string{"default": "https://sqs/default"}, "sqs")

	_, err := drv.Pop(context.Background(), "default")

	if !errors.Is(err, queue.ErrNoJob) {
		t.Fatalf("expected ErrNoJob, got %v", err)
	}
}

func TestSQSDriverPopError(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	client.receiveErr = errors.New("network error")
	drv := drivers.NewSQSDriver(client, map[string]string{"default": "https://sqs/default"}, "sqs")

	_, err := drv.Pop(context.Background(), "default")

	if !errors.Is(err, queue.ErrNoJob) {
		t.Fatalf("expected ErrNoJob, got %v", err)
	}
}

func TestSQSDriverJobDelete(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	urls := map[string]string{"default": "https://sqs/default"}
	drv := drivers.NewSQSDriver(client, urls, "sqs")

	_, _ = drv.Push(context.Background(), "default", []byte("body"))

	job, _ := drv.Pop(context.Background(), "default")

	err := job.Delete()

	if err != nil {
		t.Fatal(err)
	}
}

func TestSQSDriverJobRelease(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	urls := map[string]string{"default": "https://sqs/default"}
	drv := drivers.NewSQSDriver(client, urls, "sqs")

	_, _ = drv.Push(context.Background(), "default", []byte("body"))

	job, _ := drv.Pop(context.Background(), "default")

	err := job.Release(10 * time.Second)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSQSDriverJobFail(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	urls := map[string]string{"default": "https://sqs/default"}
	drv := drivers.NewSQSDriver(client, urls, "sqs")

	_, _ = drv.Push(context.Background(), "default", []byte("body"))

	job, _ := drv.Pop(context.Background(), "default")

	err := job.Fail(errors.New("oops"))

	if err != nil {
		t.Fatal(err)
	}
}

func TestSQSDriverSize(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	client.attrs["ApproximateNumberOfMessages"] = "42"
	drv := drivers.NewSQSDriver(client, map[string]string{"default": "https://sqs/default"}, "sqs")

	n, err := drv.Size(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 42 {
		t.Errorf("expected 42, got %d", n)
	}
}

func TestSQSDriverPendingSizeSameAsSize(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	client.attrs["ApproximateNumberOfMessages"] = "10"
	drv := drivers.NewSQSDriver(client, map[string]string{"default": "https://sqs/default"}, "sqs")

	n, err := drv.PendingSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 10 {
		t.Errorf("expected 10, got %d", n)
	}
}

func TestSQSDriverDelayedSizeAlwaysZero(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	drv := drivers.NewSQSDriver(client, map[string]string{"default": "https://sqs/default"}, "sqs")

	n, err := drv.DelayedSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestSQSDriverReservedSize(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	client.attrs["ApproximateNumberOfMessagesNotVisible"] = "5"
	drv := drivers.NewSQSDriver(client, map[string]string{"default": "https://sqs/default"}, "sqs")

	n, err := drv.ReservedSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 5 {
		t.Errorf("expected 5, got %d", n)
	}
}

func TestSQSDriverQueueURLMapping(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	urls := map[string]string{
		"emails": "https://sqs.us-east-1.amazonaws.com/123/email-queue",
	}
	drv := drivers.NewSQSDriver(client, urls, "sqs")

	_, _ = drv.Push(context.Background(), "emails", []byte("body"))

	// The mock stores messages under the URL key.
	if len(client.messages["https://sqs.us-east-1.amazonaws.com/123/email-queue"]) != 1 {
		t.Error("expected message stored under mapped URL")
	}
}

func TestSQSDriverQueueURLFallback(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	drv := drivers.NewSQSDriver(client, map[string]string{}, "sqs")

	_, _ = drv.Push(context.Background(), "https://sqs/direct-url", []byte("body"))

	// When no mapping exists, queueName is used as-is.
	if len(client.messages["https://sqs/direct-url"]) != 1 {
		t.Error("expected message stored under direct queue name as URL")
	}
}

func TestSQSDriverConnectionName(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	drv := drivers.NewSQSDriver(client, nil, "my-sqs")

	if drv.ConnectionName() != "my-sqs" {
		t.Errorf("expected 'my-sqs', got %q", drv.ConnectionName())
	}
}

func TestSQSDriverSizeError(t *testing.T) {
	t.Parallel()

	client := newMockSQSClient()
	client.attrErr = errors.New("attr error")
	drv := drivers.NewSQSDriver(client, map[string]string{"default": "https://sqs/default"}, "sqs")

	n, err := drv.Size(context.Background(), "default")

	if err == nil {
		t.Fatal("expected error")
	}

	if n != 0 {
		t.Errorf("expected 0 on error, got %d", n)
	}
}
