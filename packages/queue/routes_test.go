package queue_test

import (
	"reflect"
	"testing"

	"github.com/bedrock/packages/queue"
)

// Ports of @bedrock\Tests\Queue\QueueRoutesTest. the upstream test walks
// PHP's class_parents / class_implements / class_uses chain; Go has no
// runtime class hierarchy, so fake values here implement RouteLineage
// to report the same logical lookup chain. Behaviour under assertion
// (queue vs connection resolution, override semantics, All() shape)
// matches the PHP test line-for-line.

// --- fixtures ---------------------------------------------------------

// Route-key constants mirror the ::class literals in the PHP test.

// lineageFixture is a generic routable test fixture that reports the
// lineage it was constructed with.
type lineageFixture struct{ names []string }

const (
	keyQueueRoutes      = "@bedrock\\Queue\\QueueRoutes"
	keyBaseNotification = "@bedrock\\Tests\\Queue\\BaseNotification"
	keyCustomTrait      = "@bedrock\\Tests\\Queue\\CustomTrait"
	keyPaymentContract  = "@bedrock\\Tests\\Queue\\PaymentContract"
	keySomeJob          = "@bedrock\\Tests\\Queue\\SomeJob"
)

func (f lineageFixture) RouteLineage() []string { return f.names }

// newSomeJob = the upstream `new SomeJob` (uses Queueable + CustomTrait).
func newSomeJob() lineageFixture {
	return lineageFixture{names: []string{keySomeJob, keyCustomTrait}}
}

// newFinanceNotification = `new FinanceNotification` (extends BaseNotification).
func newFinanceNotification() lineageFixture {
	return lineageFixture{names: []string{
		"@bedrock\\Tests\\Queue\\FinanceNotification",
		keyBaseNotification,
	}}
}

// newPayment = `new Payment` (implements PaymentContract).
func newPayment() lineageFixture {
	return lineageFixture{names: []string{
		"@bedrock\\Tests\\Queue\\Payment",
		keyPaymentContract,
	}}
}

// --- ports ------------------------------------------------------------

// Port of @bedrock\Tests\Queue\QueueRoutesTest::testSet
func TestSet(t *testing.T) {
	t.Parallel()

	routes := queue.NewRoutes()

	routes.Set(keyQueueRoutes, "some-queue", "")
	routes.Set(keyBaseNotification, "some-queue", "some-connection")

	want := map[string]any{
		keyQueueRoutes:      [2]string{"", "some-queue"},
		keyBaseNotification: [2]string{"some-connection", "some-queue"},
	}

	if got := routes.All(); !reflect.DeepEqual(got, want) {
		t.Errorf("after initial Set: got %v, want %v", got, want)
	}

	// Ensure same class overrides — bulk form with a plain-string value
	// replaces the keyQueueRoutes entry, BaseNotification is untouched,
	// SomeJob is added.
	if err := routes.SetMany(map[string]any{
		keyQueueRoutes: "queue-many",
		keySomeJob:     "important",
	}); err != nil {
		t.Fatalf("SetMany: %v", err)
	}

	want = map[string]any{
		keyQueueRoutes:      "queue-many",
		keyBaseNotification: [2]string{"some-connection", "some-queue"},
		keySomeJob:          "important",
	}

	if got := routes.All(); !reflect.DeepEqual(got, want) {
		t.Errorf("after override: got %v, want %v", got, want)
	}
}

// Port of @bedrock\Tests\Queue\QueueRoutesTest::testGetQueue
func TestGetQueue(t *testing.T) {
	t.Parallel()

	routes := queue.NewRoutes()

	if err := routes.SetMany(map[string]any{
		keyBaseNotification: "notifications",
		keyCustomTrait:      "jobs",
		keyPaymentContract:  "payments",
	}); err != nil {
		t.Fatalf("SetMany: %v", err)
	}

	// Override PaymentContract with array form, connection-only.
	// Upstream: set(PaymentContract::class, connection: 'payment-connection')
	routes.Set(keyPaymentContract, "", "payment-connection")

	if got := routes.GetQueue(newFinanceNotification()); got != "notifications" {
		t.Errorf("FinanceNotification queue: got %q, want notifications", got)
	}

	if got := routes.GetQueue(newSomeJob()); got != "jobs" {
		t.Errorf("SomeJob queue: got %q, want jobs", got)
	}

	if got := routes.GetQueue(newPayment()); got != "" {
		t.Errorf("Payment queue: got %q, want empty", got)
	}
}

// Port of @bedrock\Tests\Queue\QueueRoutesTest::testGetConnection
func TestGetConnection(t *testing.T) {
	t.Parallel()

	routes := queue.NewRoutes()

	if err := routes.SetMany(map[string]any{
		keyBaseNotification: [2]string{"notification-connection", "notifications"},
		keyCustomTrait:      [2]string{"job-connection", "jobs"},
	}); err != nil {
		t.Fatalf("SetMany: %v", err)
	}

	// Override PaymentContract with queue-only (no connection).
	// Upstream: set(PaymentContract::class, 'payments')
	routes.Set(keyPaymentContract, "payments", "")

	if got := routes.GetConnection(newFinanceNotification()); got != "notification-connection" {
		t.Errorf("FinanceNotification connection: got %q, want notification-connection", got)
	}

	if got := routes.GetConnection(newSomeJob()); got != "job-connection" {
		t.Errorf("SomeJob connection: got %q, want job-connection", got)
	}

	if got := routes.GetConnection(newPayment()); got != "" {
		t.Errorf("Payment connection: got %q, want empty", got)
	}
}

// Port of @bedrock\Tests\Queue\QueueRoutesTest::testStringRouteDefaultsToQueueNotConnection
func TestStringRouteDefaultsToQueueNotConnection(t *testing.T) {
	t.Parallel()

	routes := queue.NewRoutes()

	if err := routes.SetMany(map[string]any{
		keySomeJob: "jobs",
	}); err != nil {
		t.Fatalf("SetMany: %v", err)
	}

	job := newSomeJob()

	if got := routes.GetQueue(job); got != "jobs" {
		t.Errorf("GetQueue: got %q, want jobs", got)
	}

	if got := routes.GetConnection(job); got != "" {
		t.Errorf("GetConnection: got %q, want empty", got)
	}
}
