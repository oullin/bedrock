package bus_test

import (
	"context"
	"sync"
	"testing"

	"github.com/bedrock/packages/bus"
)

type testCommand struct{ Value string }

// selfHandlingCommand implements the Handle(ctx) interface.
type selfHandlingCommand struct {
	Value  string
	Called bool
}

// DispatchSync should bypass pipeline.

// Dispatch should go through pipeline.

type customQueueCommand struct {
	Value string
}

type anotherCommand struct{ Name string }

func TestDispatcherSyncDispatch(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)
	d.Map(testCommand{}, func(ctx context.Context, cmd any) (any, error) {
		c := cmd.(testCommand)

		return c.Value + "_handled", nil
	})

	result, err := d.Dispatch(context.Background(), testCommand{Value: "hello"})

	if err != nil {
		t.Fatal(err)
	}

	if result != "hello_handled" {
		t.Errorf("got %v, want hello_handled", result)
	}
}

func TestDispatcherPipelineMiddleware(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	order := []string{}
	d.PipeThrough(
		func(ctx context.Context, cmd any, next bus.Handler) (any, error) {
			order = append(order, "pipe1_before")
			result, err := next(ctx, cmd)
			order = append(order, "pipe1_after")

			return result, err
		},
		func(ctx context.Context, cmd any, next bus.Handler) (any, error) {
			order = append(order, "pipe2_before")
			result, err := next(ctx, cmd)
			order = append(order, "pipe2_after")

			return result, err
		},
	)

	d.Map(testCommand{}, func(ctx context.Context, _ any) (any, error) {
		order = append(order, "handler")

		return nil, nil
	})

	_, _ = d.Dispatch(context.Background(), testCommand{})

	expected := []string{"pipe1_before", "pipe2_before", "handler", "pipe2_after", "pipe1_after"}

	for i, e := range expected {
		if i >= len(order) || order[i] != e {
			t.Errorf("pipeline order[%d]: got %q, want %q", i, func() string {
				if i < len(order) {
					return order[i]
				}

				return "<missing>"
			}(), e)
		}
	}
}

func TestDispatcherDeferredFlush(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	count := 0
	d.Map(testCommand{}, func(_ context.Context, _ any) (any, error) {
		count++

		return nil, nil
	})

	_ = d.DispatchAfterResponse(context.Background(), testCommand{})
	_ = d.DispatchAfterResponse(context.Background(), testCommand{})

	if count != 0 {
		t.Error("deferred commands should not run before Flush")
	}

	if err := d.FlushDeferred(context.Background()); err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Errorf("expected 2 deferred commands run, got %d", count)
	}
}

func TestDispatcherNoHandlerError(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	_, err := d.Dispatch(context.Background(), testCommand{})

	if err == nil {
		t.Error("expected error for unregistered command")
	}
}

func (c *selfHandlingCommand) Handle(_ context.Context) (any, error) {
	c.Called = true

	return c.Value + "_self", nil
}

func TestDispatcherSelfHandlingCommand(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	cmd := &selfHandlingCommand{Value: "hello"}
	result, err := d.Dispatch(context.Background(), cmd)

	if err != nil {
		t.Fatal(err)
	}

	if !cmd.Called {
		t.Error("expected self-handling Handle() to be called")
	}

	if result != "hello_self" {
		t.Errorf("got %v, want hello_self", result)
	}
}

func TestDispatcherMappedHandlerOverSelfHandling(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	d.Map(&selfHandlingCommand{}, func(_ context.Context, cmd any) (any, error) {
		return "mapped", nil
	})

	cmd := &selfHandlingCommand{Value: "hello"}
	result, err := d.Dispatch(context.Background(), cmd)

	if err != nil {
		t.Fatal(err)
	}

	if cmd.Called {
		t.Error("expected mapped handler to take precedence; self-handling should not be called")
	}

	if result != "mapped" {
		t.Errorf("got %v, want mapped", result)
	}
}

func TestDispatchSyncBypassesPipeline(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	pipeCalled := false
	d.PipeThrough(func(_ context.Context, cmd any, next bus.Handler) (any, error) {
		pipeCalled = true

		return next(context.Background(), cmd)
	})

	d.Map(testCommand{}, func(_ context.Context, cmd any) (any, error) {
		return cmd.(testCommand).Value, nil
	})

	_, _ = d.DispatchSync(context.Background(), testCommand{Value: "sync"})

	if pipeCalled {
		t.Error("expected DispatchSync to bypass pipeline")
	}

	_, _ = d.Dispatch(context.Background(), testCommand{Value: "dispatch"})

	if !pipeCalled {
		t.Error("expected Dispatch to go through pipeline")
	}
}

func TestDispatchNowAliasesDispatchSync(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)
	d.Map(testCommand{}, func(_ context.Context, cmd any) (any, error) {
		return cmd.(testCommand).Value + "_result", nil
	})

	syncResult, syncErr := d.DispatchSync(context.Background(), testCommand{Value: "a"})
	nowResult, nowErr := d.DispatchNow(context.Background(), testCommand{Value: "a"})

	if syncErr != nil || nowErr != nil {
		t.Fatalf("unexpected errors: sync=%v, now=%v", syncErr, nowErr)
	}

	if syncResult != nowResult {
		t.Errorf("DispatchSync returned %v, DispatchNow returned %v", syncResult, nowResult)
	}
}

func TestDispatchToQueueNoBackendError(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	err := d.DispatchToQueue(context.Background(), testCommand{Value: "x"})

	if err == nil {
		t.Error("expected error when no queue backend configured")
	}
}

func TestDispatchToQueueDefaultQueueName(t *testing.T) {
	q := newMockQueue()
	d := bus.NewDispatcher(q, nil)

	err := d.DispatchToQueue(context.Background(), testCommand{Value: "test"})

	if err != nil {
		t.Fatal(err)
	}

	q.mu.Lock()

	defer q.mu.Unlock()

	if len(q.pushes) != 1 {
		t.Fatalf("expected 1 push, got %d", len(q.pushes))
	}

	if q.pushes[0].Queue != "default" {
		t.Errorf("expected queue name 'default', got %q", q.pushes[0].Queue)
	}
}

func (c customQueueCommand) GetQueue() string { return "emails" }

func TestDispatchToQueueCustomQueueName(t *testing.T) {
	q := newMockQueue()
	d := bus.NewDispatcher(q, nil)

	err := d.DispatchToQueue(context.Background(), customQueueCommand{Value: "test"})

	if err != nil {
		t.Fatal(err)
	}

	q.mu.Lock()

	defer q.mu.Unlock()

	if len(q.pushes) != 1 {
		t.Fatalf("expected 1 push, got %d", len(q.pushes))
	}

	if q.pushes[0].Queue != "emails" {
		t.Errorf("expected queue name 'emails', got %q", q.pushes[0].Queue)
	}
}

func TestDispatchToQueueSerializesPayload(t *testing.T) {
	q := newMockQueue()
	d := bus.NewDispatcher(q, nil)

	err := d.DispatchToQueue(context.Background(), testCommand{Value: "hello"})

	if err != nil {
		t.Fatal(err)
	}

	q.mu.Lock()

	defer q.mu.Unlock()

	if len(q.pushes) == 0 {
		t.Fatal("expected at least 1 push")
	}

	payload := string(q.pushes[0].Payload)

	if payload == "" {
		t.Error("expected non-empty payload")
	}
}

func TestDispatcherMapMultipleCommands(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	d.Map(testCommand{}, func(_ context.Context, cmd any) (any, error) {
		return "test:" + cmd.(testCommand).Value, nil
	})
	d.Map(anotherCommand{}, func(_ context.Context, cmd any) (any, error) {
		return "another:" + cmd.(anotherCommand).Name, nil
	})

	r1, _ := d.Dispatch(context.Background(), testCommand{Value: "A"})
	r2, _ := d.Dispatch(context.Background(), anotherCommand{Name: "B"})

	if r1 != "test:A" {
		t.Errorf("got %v, want test:A", r1)
	}

	if r2 != "another:B" {
		t.Errorf("got %v, want another:B", r2)
	}
}

func TestFlushDeferredStopsOnError(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	count := 0
	d.Map(testCommand{}, func(_ context.Context, _ any) (any, error) {
		count++

		if count == 2 {
			return nil, errTestFailure
		}

		return nil, nil
	})

	_ = d.DispatchAfterResponse(context.Background(), testCommand{})
	_ = d.DispatchAfterResponse(context.Background(), testCommand{})
	_ = d.DispatchAfterResponse(context.Background(), testCommand{})

	err := d.FlushDeferred(context.Background())

	if err == nil {
		t.Error("expected error from FlushDeferred")
	}

	if count != 2 {
		t.Errorf("expected 2 commands attempted, got %d", count)
	}
}

func TestFlushDeferredClearsBuffer(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	count := 0
	d.Map(testCommand{}, func(_ context.Context, _ any) (any, error) {
		count++

		return nil, nil
	})

	_ = d.DispatchAfterResponse(context.Background(), testCommand{})

	if err := d.FlushDeferred(context.Background()); err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("expected 1, got %d", count)
	}

	// Second flush should be a no-op.
	if err := d.FlushDeferred(context.Background()); err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Errorf("expected count to remain 1 after second flush, got %d", count)
	}
}

func TestDispatcherConcurrentMapAndDispatch(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)
	d.Map(testCommand{}, func(_ context.Context, _ any) (any, error) {
		return "ok", nil
	})

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()

			d.Map(testCommand{}, func(_ context.Context, _ any) (any, error) {
				return "ok", nil
			})
		}()

		go func() {
			defer wg.Done()

			_, _ = d.Dispatch(context.Background(), testCommand{Value: "concurrent"})
		}()
	}

	wg.Wait()
}

func TestFindBatchNoRepoError(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	_, err := d.FindBatch(context.Background(), "any-id")

	if err == nil {
		t.Error("expected error when no batch repository configured")
	}
}

func TestFindBatchDelegatesToRepo(t *testing.T) {
	repo := newMockBatchRepo()
	expected := &bus.Batch{ID: "batch-42", Name: "test"}
	repo.batch = expected

	d := bus.NewDispatcher(nil, repo)

	result, err := d.FindBatch(context.Background(), "batch-42")

	if err != nil {
		t.Fatal(err)
	}

	if result.ID != expected.ID {
		t.Errorf("expected batch ID %q, got %q", expected.ID, result.ID)
	}
}

func TestBatchReturnsPendingBatch(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	pb := d.Batch([]any{"job1", "job2"})

	if pb == nil {
		t.Fatal("expected non-nil PendingBatch")
	}
}

func TestHasCommandHandler(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	if d.HasCommandHandler(testCommand{}) {
		t.Error("expected HasCommandHandler to be false before mapping")
	}

	d.Map(testCommand{}, func(_ context.Context, _ any) (any, error) { return nil, nil })

	if !d.HasCommandHandler(testCommand{}) {
		t.Error("expected HasCommandHandler to be true after mapping")
	}
}

func TestGetCommandHandler(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	_, ok := d.GetCommandHandler(testCommand{})

	if ok {
		t.Error("expected GetCommandHandler to return false for unmapped type")
	}

	d.Map(testCommand{}, func(_ context.Context, _ any) (any, error) { return "found", nil })

	handler, ok := d.GetCommandHandler(testCommand{})

	if !ok {
		t.Error("expected GetCommandHandler to return true for mapped type")
	}

	result, _ := handler(context.Background(), testCommand{})

	if result != "found" {
		t.Errorf("expected 'found', got %v", result)
	}
}

func TestChainReturnsNonNil(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	chain := d.Chain([]any{"j1", "j2"})

	if chain == nil {
		t.Fatal("expected non-nil PendingChain")
	}
}
