package bus_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/bus"
)

type chainJob struct {
	bus.Queueable
	Value string
}

func TestPendingChainDispatch(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	j1 := &chainJob{Value: "first"}
	j2 := &chainJob{Value: "second"}
	j3 := &chainJob{Value: "third"}

	d.Map(j1, func(_ context.Context, cmd any) (any, error) {
		j := cmd.(*chainJob)

		return j.Value, nil
	})

	chain := bus.NewPendingChain(d, []any{j1, j2, j3})
	result, err := chain.Dispatch(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if result != "first" {
		t.Errorf("expected 'first', got %v", result)
	}

	if len(j1.ChainJobs) != 2 {
		t.Errorf("expected 2 chain jobs on first job, got %d", len(j1.ChainJobs))
	}
}

func TestPendingChainOnConnectionOnQueue(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	j1 := &chainJob{Value: "first"}
	d.Map(j1, func(_ context.Context, cmd any) (any, error) {
		return nil, nil
	})

	chain := bus.NewPendingChain(d, []any{j1}).
		OnConnection("redis").
		OnQueue("high")

	_, err := chain.Dispatch(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if j1.Connection != "redis" {
		t.Errorf("expected connection 'redis', got %q", j1.Connection)
	}

	if j1.Queue != "high" {
		t.Errorf("expected queue 'high', got %q", j1.Queue)
	}
}

func TestPendingChainEmptyError(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)
	chain := bus.NewPendingChain(d, []any{})

	_, err := chain.Dispatch(context.Background())

	if err == nil {
		t.Error("expected error for empty chain")
	}
}
