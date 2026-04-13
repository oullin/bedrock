package drivers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

func TestBeanstalkdDriverPush(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	id, err := drv.Push(context.Background(), "default", []byte("payload"))

	if err != nil {
		t.Fatal(err)
	}

	if id != "1" {
		t.Errorf("expected id '1', got %q", id)
	}

	if client.lastPut.priority != 1024 {
		t.Errorf("expected priority 1024, got %d", client.lastPut.priority)
	}

	if client.lastPut.delay != 0 {
		t.Errorf("expected delay 0, got %v", client.lastPut.delay)
	}

	if client.lastPut.ttr != 60*time.Second {
		t.Errorf("expected ttr 60s, got %v", client.lastPut.ttr)
	}
}

func TestBeanstalkdDriverPushDelayed(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	_, err := drv.PushDelayed(context.Background(), "default", []byte("payload"), 30*time.Second)

	if err != nil {
		t.Fatal(err)
	}

	if client.lastPut.delay != 30*time.Second {
		t.Errorf("expected delay 30s, got %v", client.lastPut.delay)
	}
}

func TestBeanstalkdDriverPushMultiple(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	ids, err := drv.PushMultiple(context.Background(), "default", [][]byte{
		[]byte("a"), []byte("b"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 2 {
		t.Errorf("expected 2 ids, got %d", len(ids))
	}
}

func TestBeanstalkdDriverPop(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	_, _ = drv.Push(context.Background(), "default", []byte("body"))

	job, err := drv.Pop(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if string(job.Payload()) != "body" {
		t.Errorf("expected 'body', got %q", job.Payload())
	}

	if job.GetQueue() != "default" {
		t.Errorf("expected queue 'default', got %q", job.GetQueue())
	}

	if job.GetConnectionName() != "beanstalkd" {
		t.Errorf("expected connection 'beanstalkd', got %q", job.GetConnectionName())
	}
}

func TestBeanstalkdDriverPopEmpty(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	client.reserveErr = errors.New("timed out")
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	_, err := drv.Pop(context.Background(), "default")

	if !errors.Is(err, queue.ErrNoJob) {
		t.Fatalf("expected ErrNoJob, got %v", err)
	}
}

func TestBeanstalkdDriverJobDelete(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	_, _ = drv.Push(context.Background(), "default", []byte("body"))

	job, _ := drv.Pop(context.Background(), "default")

	err := job.Delete()

	if err != nil {
		t.Fatal(err)
	}
}

func TestBeanstalkdDriverJobRelease(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	_, _ = drv.Push(context.Background(), "default", []byte("body"))

	job, _ := drv.Pop(context.Background(), "default")

	err := job.Release(5 * time.Second)

	if err != nil {
		t.Fatal(err)
	}
}

func TestBeanstalkdDriverJobFail(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	_, _ = drv.Push(context.Background(), "default", []byte("body"))

	job, _ := drv.Pop(context.Background(), "default")

	err := job.Fail(errors.New("oops"))

	if err != nil {
		t.Fatal(err)
	}
}

func TestBeanstalkdDriverSize(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	client.stats["default"] = map[string]string{"current-jobs-ready": "7"}
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	n, err := drv.Size(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 7 {
		t.Errorf("expected 7, got %d", n)
	}
}

func TestBeanstalkdDriverDelayedSize(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	client.stats["default"] = map[string]string{"current-jobs-delayed": "3"}
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	n, err := drv.DelayedSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 3 {
		t.Errorf("expected 3, got %d", n)
	}
}

func TestBeanstalkdDriverReservedSize(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	client.stats["default"] = map[string]string{"current-jobs-reserved": "2"}
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	n, err := drv.ReservedSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}

func TestBeanstalkdDriverDefaultTTR(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 0) // 0 should default to 60s

	_, _ = drv.Push(context.Background(), "default", []byte("body"))

	if client.lastPut.ttr != 60*time.Second {
		t.Errorf("expected default ttr 60s, got %v", client.lastPut.ttr)
	}
}

func TestBeanstalkdDriverCustomTTR(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 120*time.Second)

	_, _ = drv.Push(context.Background(), "default", []byte("body"))

	if client.lastPut.ttr != 120*time.Second {
		t.Errorf("expected ttr 120s, got %v", client.lastPut.ttr)
	}
}

func TestBeanstalkdDriverConnectionName(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	drv := drivers.NewBeanstalkdDriver(client, "my-beanstalkd", 60*time.Second)

	if drv.ConnectionName() != "my-beanstalkd" {
		t.Errorf("expected 'my-beanstalkd', got %q", drv.ConnectionName())
	}
}

func TestBeanstalkdDriverPendingSize(t *testing.T) {
	t.Parallel()

	client := newMockBeanstalkdClient()
	client.stats["default"] = map[string]string{"current-jobs-ready": "4"}
	drv := drivers.NewBeanstalkdDriver(client, "beanstalkd", 60*time.Second)

	n, err := drv.PendingSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 4 {
		t.Errorf("expected 4, got %d", n)
	}
}
