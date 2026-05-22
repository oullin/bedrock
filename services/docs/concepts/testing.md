# Testing

<!-- ref: @bedrock/code-0175 -->
<!-- ref: @bedrock/code-0173 -->
<!-- ref: @bedrock/code-0174 -->

A tour of Bedrock's test doubles and testing patterns.

## Philosophy

Every Bedrock package that performs I/O ships with at least one in-memory or
no-op implementation. This lets you write fast, deterministic tests without
mocks, network, or disk.

## Built-in Test Doubles

| Package       | Fake / stub                       | Purpose                             |
| ------------- | --------------------------------- | ----------------------------------- |
| `cache`       | `MemoryStore` / `NullStore`       | Swap out Redis or database in tests |
| `queue`       | `SyncConnector` / `NullConnector` | Run jobs inline or discard them     |
| `events`      | `NullDispatcher`                  | Silence all event dispatch          |
| `mailx`       | `ArrayTransport`                  | Capture sent messages in memory     |
| `session`     | `ArrayHandler`                    | In-memory session storage           |
| `concurrency` | `SyncDriver`                      | Run parallel code sequentially      |
| `log`         | `NullHandler`                     | Discard all log records             |
| `httpx`       | `TestRequest` / `TestResponse`    | Build requests, assert responses    |
| `redis`       | Custom `Client` interface         | Inject a fake Redis client          |

## Example: Testing a Mail-Sending Service

```go
func TestWelcomeService_SendsWelcomeMail(t *testing.T) {
    arr := mailx.NewArrayTransport()
    mailer := mailx.NewMailer("test", arr)

    svc := &WelcomeService{Mailer: mailer}

    err := svc.SendWelcome(context.Background(), &User{Email: "a@b.com"})
    require.NoError(t, err)

    sent := arr.Messages()
    require.Len(t, sent, 1)
    require.Equal(t, "a@b.com", sent[0].Envelope().To[0].Email)
}
```

## Example: Testing a Queued Job

```go
func TestProcessOrder_DispatchesShippingJob(t *testing.T) {
    connector := queue.NewSyncConnector() // runs jobs inline
    q := connector.Connect(nil)

    manager := queue.NewManager()
    manager.AddConnector("sync", connector)

    handler := &OrderHandler{Queue: manager}
    err := handler.Process(context.Background(), &Order{ID: 42})
    require.NoError(t, err)

    // Sync connector executed the job, so side effects are observable
    require.True(t, shippingLabelWasCreated(42))
}
```

## Example: Testing Events

```go
func TestRegisterUser_FiresUserRegisteredEvent(t *testing.T) {
    d := events.NewDispatcher()

    var fired bool
    d.Listen(UserRegistered{}, func(ctx context.Context, e any) error {
        fired = true
        return nil
    })

    svc := &RegistrationService{Events: d}
    _, _ = svc.Register(context.Background(), "a@b.com", "password")

    require.True(t, fired)
}
```

## Example: Freezing the Clock

`contracts.Clock` exists precisely so tests can freeze time:

```go
type FakeClock struct{ now time.Time }
func (c *FakeClock) Now() time.Time { return c.now }
func (c *FakeClock) Advance(d time.Duration) { c.now = c.now.Add(d) }

clock := &FakeClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
svc := &TokenService{Clock: clock}

token, _ := svc.Issue("user-42")
clock.Advance(2 * time.Hour)
require.True(t, svc.IsExpired(token)) // deterministic, no sleep
```

## Example: Fake Sleep

`support.Sleep` can be faked in tests to avoid real delays:

```go
support.FakeSleep()
defer support.ResetSleep()

svc := &RetryService{}
svc.DoWithBackoff()

support.AssertSlept(100*time.Millisecond, 3) // slept 3 times
```

## Example: HTTP Handler Tests

The `httpx/testing` sub-package provides a fluent request builder:

```go
func TestShowUser(t *testing.T) {
    req := httptesting.NewTestRequest("GET", "/users/42").
        WithHeader("Accept", "application/json")

    res := httptesting.Call(router, req)

    res.AssertStatus(http.StatusOK)
    res.AssertJsonPath("id", 42)
    res.AssertJsonPath("email", "alice@example.com")
}
```

## Upstream Parity Tests

Upstream-parity tests previously lived alongside each package as
`*_upstream_test.go` and `*_inventory_test.go` files. They have been moved to
the `bedrock-compliance` repository, which tracks parity against pinned
upstream sources on its own schedule. See that repo for parity status,
divergence rationales, and feature audit coverage.

## Running Tests

```bash
# All packages
pnpm run test

# A single package
go test ./packages/validation/...

# With coverage
go test -cover ./packages/cache/...

# Verbose
go test -v ./packages/routing/...
```
