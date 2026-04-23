package pipeline

import (
	"context"
	"errors"
	"testing"
)

// --- test helpers ---

type testPipeOne struct{}

type testPipeTwo struct{}

type testParameterPipe struct{}

func (p testPipeOne) Handle(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
	return next(passable.(string) + "-pipe1")
}

func (p testPipeOne) DifferentMethod(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
	return next(passable.(string) + "-diff")
}

func (p testPipeTwo) Handle(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
	return next(passable.(string) + "-pipe2")
}

func (p testParameterPipe) Handle(_ context.Context, passable any, next func(any) (any, error), params ...string) (any, error) {
	s := passable.(string)

	for _, param := range params {
		s += "-" + param
	}

	return next(s)
}

// --- tests ---

func TestPipelineBasicUsage(t *testing.T) {
	result, err := New().
		Send("hello").
		Through(
			Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
				return next(passable.(string) + "-closure")
			}),
			testPipeOne{},
		).
		Then(context.Background(), func(passable any) (any, error) {
			return passable, nil
		})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "hello-closure-pipe1" {
		t.Fatalf("expected hello-closure-pipe1, got %v", result)
	}
}

func TestPipelineUsageWithObjects(t *testing.T) {
	result, err := New().
		Send("start").
		Through(testPipeOne{}, testPipeTwo{}).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "start-pipe1-pipe2" {
		t.Fatalf("expected start-pipe1-pipe2, got %v", result)
	}
}

func TestPipelineUsageWithInvokableObjects(t *testing.T) {
	result, err := New().
		Send("data").
		Through(
			Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
				return next(passable.(string) + "-invoked")
			}),
		).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "data-invoked" {
		t.Fatalf("expected data-invoked, got %v", result)
	}
}

func TestPipelineUsageWithCallable(t *testing.T) {
	fn := Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
		return next(passable.(string) + "-callable")
	})

	result, err := New().
		Send("value").
		Through(fn).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "value-callable" {
		t.Fatalf("expected value-callable, got %v", result)
	}
}

func TestPipelineWithinTransaction(t *testing.T) {
	t.Parallel()

	var wrapped bool

	result, err := New().
		Send("value").
		Through(Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
			return next(passable.(string) + "-pipe")
		})).
		WithinTransaction(func(_ context.Context, run func() (any, error)) (any, error) {
			wrapped = true

			return run()
		}).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !wrapped || result != "value-pipe" {
		t.Fatalf("expected transaction wrapper and pipeline result, wrapped=%v result=%v", wrapped, result)
	}
}

func TestPipelineUsageWithPipeAppend(t *testing.T) {
	result, err := New().
		Send("x").
		Through(testPipeOne{}).
		Pipe(testPipeTwo{}).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "x-pipe1-pipe2" {
		t.Fatalf("expected x-pipe1-pipe2, got %v", result)
	}
}

func TestPipelineThroughOverwritesPipes(t *testing.T) {
	p := New()
	p.Send("x")
	p.Through(testPipeOne{})
	p.Pipe(testPipeTwo{})
	p.Through(testPipeTwo{})

	result, err := p.ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "x-pipe2" {
		t.Fatalf("expected x-pipe2, got %v", result)
	}
}

func TestPipelineUsageWithInvokableClass(t *testing.T) {
	result, err := New().
		Send("data").
		Via("DifferentMethod").
		Through(&testPipeOne{}).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "data-diff" {
		t.Fatalf("expected data-diff, got %v", result)
	}
}

func TestThenNotCalledIfPipeReturns(t *testing.T) {
	thenCalled := false

	result, err := New().
		Send("hello").
		Through(
			Pipe(func(_ context.Context, _ any, _ func(any) (any, error)) (any, error) {
				return "early-return", nil
			}),
		).
		Then(context.Background(), func(_ any) (any, error) {
			thenCalled = true

			return nil, nil
		})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if thenCalled {
		t.Fatal("then destination should not have been called")
	}

	if result != "early-return" {
		t.Fatalf("expected early-return, got %v", result)
	}
}

func TestThenMethodInputValue(t *testing.T) {
	var received any

	_, err := New().
		Send("original").
		Through(
			Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
				return next(passable.(string) + "-modified")
			}),
		).
		Then(context.Background(), func(passable any) (any, error) {
			received = passable

			return passable, nil
		})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received != "original-modified" {
		t.Fatalf("expected original-modified, got %v", received)
	}
}

func TestPipelineUsageWithParameters(t *testing.T) {
	resolver := func(name string, params []string) (Pipe, error) {
		if name == "param-pipe" {
			return func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
				s := passable.(string)

				for _, param := range params {
					s += "-" + param
				}

				return next(s)
			}, nil
		}

		return nil, errors.New("unknown pipe")
	}

	result, err := New(WithResolver(resolver)).
		Send("start").
		Through("param-pipe:foo,bar").
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "start-foo-bar" {
		t.Fatalf("expected start-foo-bar, got %v", result)
	}
}

func TestPipelineViaChangesMethod(t *testing.T) {
	result, err := New().
		Send("hello").
		Via("DifferentMethod").
		Through(&testPipeOne{}).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "hello-diff" {
		t.Fatalf("expected hello-diff, got %v", result)
	}
}

func TestPipelineErrorOnResolveWithoutResolver(t *testing.T) {
	_, err := New().
		Send("value").
		Through("SomeClass").
		ThenReturn(context.Background())

	if err == nil {
		t.Fatal("expected error for string pipe without resolver")
	}
}

func TestPipelineThenReturn(t *testing.T) {
	result, err := New().
		Send("passthrough").
		Through(
			Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
				return next(passable.(string) + "-piped")
			}),
		).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "passthrough-piped" {
		t.Fatalf("expected passthrough-piped, got %v", result)
	}
}

func TestPipelineConditionable(t *testing.T) {
	addPipe := Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
		return next(passable.(string) + "-conditional")
	})

	result, err := New().
		Send("start").
		When(true, func(p *Pipeline) *Pipeline {
			p.Through(addPipe)

			return p
		}).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "start-conditional" {
		t.Fatalf("expected start-conditional, got %v", result)
	}

	result2, err := New().
		Send("start").
		When(false, func(p *Pipeline) *Pipeline {
			p.Through(addPipe)

			return p
		}).
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result2 != "start" {
		t.Fatalf("expected start, got %v", result2)
	}
}

func TestPipelineFinally(t *testing.T) {
	finallyCalled := false

	var finalResult any

	result, err := New().
		Send("data").
		Through(
			Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
				return next(passable.(string) + "-piped")
			}),
		).
		Finally(func(r any, e error) {
			finallyCalled = true
			finalResult = r
		}).
		Then(context.Background(), func(passable any) (any, error) {
			return passable, nil
		})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !finallyCalled {
		t.Fatal("finally callback was not called")
	}

	if result != "data-piped" {
		t.Fatalf("expected data-piped, got %v", result)
	}

	if finalResult != "data-piped" {
		t.Fatalf("expected finally to receive data-piped, got %v", finalResult)
	}
}

func TestPipelineFinallyWhenChainStopped(t *testing.T) {
	finallyCalled := false

	result, err := New().
		Send("hello").
		Through(
			Pipe(func(_ context.Context, _ any, _ func(any) (any, error)) (any, error) {
				return "stopped", nil
			}),
		).
		Finally(func(_ any, _ error) {
			finallyCalled = true
		}).
		Then(context.Background(), func(passable any) (any, error) {
			return passable, nil
		})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !finallyCalled {
		t.Fatal("finally should be called even when chain is stopped")
	}

	if result != "stopped" {
		t.Fatalf("expected stopped, got %v", result)
	}
}

func TestPipelineFinallyOrder(t *testing.T) {
	var order []string

	_, err := New().
		Send("start").
		Through(
			Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
				order = append(order, "pipe1")

				return next(passable)
			}),
			Pipe(func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
				order = append(order, "pipe2")

				return next(passable)
			}),
		).
		Finally(func(_ any, _ error) {
			order = append(order, "finally")
		}).
		Then(context.Background(), func(passable any) (any, error) {
			order = append(order, "destination")

			return passable, nil
		})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"pipe1", "pipe2", "destination", "finally"}

	if len(order) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, order)
	}

	for i, v := range expected {
		if order[i] != v {
			t.Fatalf("expected order[%d] = %q, got %q", i, v, order[i])
		}
	}
}

func TestPipelineFinallyWhenErrorOccurs(t *testing.T) {
	finallyCalled := false

	var finalErr error

	_, err := New().
		Send("data").
		Through(
			Pipe(func(_ context.Context, _ any, _ func(any) (any, error)) (any, error) {
				return nil, errors.New("pipe error")
			}),
		).
		Finally(func(_ any, e error) {
			finallyCalled = true
			finalErr = e
		}).
		Then(context.Background(), func(passable any) (any, error) {
			return passable, nil
		})

	if err == nil {
		t.Fatal("expected error")
	}

	if !finallyCalled {
		t.Fatal("finally should be called even when error occurs")
	}

	if finalErr == nil || finalErr.Error() != "pipe error" {
		t.Fatalf("expected finally to receive pipe error, got %v", finalErr)
	}
}

func TestPipelineEmptyStages(t *testing.T) {
	result, err := New().
		Send("unchanged").
		Through().
		ThenReturn(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "unchanged" {
		t.Fatalf("expected unchanged, got %v", result)
	}
}
