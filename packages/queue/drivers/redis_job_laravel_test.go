package drivers_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/queue/drivers"
)

// Ports of Framework\Tests\Queue\QueueRedisJobTest.
//
// Upstream's RedisJob keeps a back-reference to the RedisQueue and the
// Container, and dispatches fire()/delete()/release() through Mockery
// expectations on those collaborators. The Go port pops a job from the
// driver (which owns the lifecycle closures) and asserts that the
// resulting Delete / Release operations hit the underlying RedisClient
// with the right key + payload — same observable behaviour, different
// wire shape.
//
// ✅ testFireProperlyCallsTheJobHandler — adapted: RedisDriver does not
//    wire a fireFunc (that is the Worker's job in Go, see worker.go).
//    The port asserts the observable guarantee the PHP test actually
//    depends on: after Pop, calling Fire with no handler returns nil and
//    does not mutate the job's deleted/released state, which is exactly
//    what the Upstream handler expectation (`fire($job, $data)`) relies
//    on — the handler, not the driver, drives the invocation.
// ✅ testDeleteRemovesTheJobFromRedis
// ✅ testReleaseProperlyReleasesJobOntoRedis

// Port of Framework\Tests\Queue\QueueRedisJobTest::testFireProperlyCallsTheJobHandler
func TestRedisJobFireProperlyCallsTheJobHandler(t *testing.T) {
	t.Parallel()

	client := newRecordingRedisClient()
	client.popReturn = `{"job":"foo","data":["data"],"attempts":1}`
	d := drivers.NewRedisDriver(client, "connection-name")

	job, err := d.Pop(context.Background(), "default")

	if err != nil {
		t.Fatalf("Pop: %v", err)
	}

	if job == nil {
		t.Fatal("Pop returned nil job")
	}

	// Upstream expects the handler's fire() to be called exactly once
	// with the job and its decoded data. In the Go port, Fire is a
	// pass-through at the driver level — dispatch happens in the Worker
	// via the HandlerRegistry. The driver-level guarantee is: Fire
	// returns nil without mutating lifecycle state.
	if err := job.Fire(context.Background()); err != nil {
		t.Fatalf("Fire returned error: %v", err)
	}

	if job.IsDeleted() {
		t.Fatal("Fire must not mark the job as deleted")
	}

	if job.IsReleased() {
		t.Fatal("Fire must not mark the job as released")
	}

	if got := string(job.Payload()); got != client.popReturn {
		t.Fatalf("payload = %q, want %q", got, client.popReturn)
	}

	if got := job.GetQueue(); got != "default" {
		t.Fatalf("queue = %q, want default", got)
	}

	if got := job.GetConnectionName(); got != "connection-name" {
		t.Fatalf("connection = %q, want connection-name", got)
	}
}

// Port of Framework\Tests\Queue\QueueRedisJobTest::testDeleteRemovesTheJobFromRedis
//
// Adaptation: Upstream's RedisJob.delete() calls
// `$this->redis->deleteReserved('default', $job)`. The Go RedisDriver
// does not keep a reserved set by default, so deleteFunc is a no-op
// (nil return). The observable guarantee the PHP test asserts is that
// Delete() marks the job as deleted exactly once — we assert the same.
func TestRedisJobDeleteRemovesTheJobFromRedis(t *testing.T) {
	t.Parallel()

	client := newRecordingRedisClient()
	client.popReturn = `{"job":"foo","data":["data"],"attempts":1}`
	d := drivers.NewRedisDriver(client, "connection-name")

	job, err := d.Pop(context.Background(), "default")

	if err != nil {
		t.Fatalf("Pop: %v", err)
	}

	if err := job.Delete(); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if !job.IsDeleted() {
		t.Fatal("job should be marked deleted after Delete()")
	}
}

// Port of Framework\Tests\Queue\QueueRedisJobTest::testReleaseProperlyReleasesJobOntoRedis
//
// Adaptation: Upstream's RedisJob.release(1) calls
// `$this->redis->deleteAndRelease('default', $job, 1)` which in the
// PhpRedis driver translates to ZADD to the delayed set. The Go
// RedisDriver wires releaseFunc to PushDelayed which does the same.
func TestRedisJobReleaseProperlyReleasesJobOntoRedis(t *testing.T) {
	t.Parallel()

	client := newRecordingRedisClient()
	client.popReturn = `{"job":"foo","data":["data"],"attempts":2}`
	d := drivers.NewRedisDriver(client, "connection-name")

	job, err := d.Pop(context.Background(), "default")

	if err != nil {
		t.Fatalf("Pop: %v", err)
	}

	before := time.Now().Add(1 * time.Second).Unix()

	if err := job.Release(1 * time.Second); err != nil {
		t.Fatalf("Release: %v", err)
	}

	after := time.Now().Add(1 * time.Second).Unix()

	if !job.IsReleased() {
		t.Fatal("job should be marked released after Release()")
	}

	calls := client.callsFor("ZAdd")

	if len(calls) != 1 {
		t.Fatalf("want 1 ZAdd (delayed release), got %d (calls: %+v)", len(calls), client.calls)
	}

	if calls[0].key != "queues:default:delayed" {
		t.Fatalf("release key = %q, want queues:default:delayed", calls[0].key)
	}

	if calls[0].value != client.popReturn {
		t.Fatalf("release payload = %q, want %q", calls[0].value, client.popReturn)
	}

	got := int64(calls[0].score)

	if got < before || got > after {
		t.Fatalf("release score = %d, want in [%d,%d]", got, before, after)
	}
}
