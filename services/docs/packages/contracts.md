# contracts

Shared interface definitions for every Bedrock package.

## Overview

The `contracts` package is the interface hub of Bedrock. Every concrete package
(`auth`, `cache`, `mailx`, `queue`, …) depends on interfaces declared here
rather than on other concrete packages, enabling loose coupling, test doubles,
and drop-in replacements.

**Module:** `github.com/gocanto/bedrock/packages/contracts`

```bash
go get github.com/gocanto/bedrock/packages/contracts@latest
```

This package has no runtime logic — only interface definitions.

## Sub-packages

Interfaces are grouped by domain. Import only the sub-package you need:

| Sub-package               | Key interfaces                                                   |
| ------------------------- | ---------------------------------------------------------------- |
| `contracts/auth`          | `Authenticatable`, `Guard`, `UserProvider`, `PasswordHasher`     |
| `contracts/concurrency`   | `Driver`                                                         |
| `contracts/encryption`    | `Encrypter`, `StringEncrypter`                                   |
| `contracts/events`        | `Dispatcher`, `Listener`, `Subscriber`                           |
| `contracts/filesystem`    | `Filesystem`                                                     |
| `contracts/hashing`       | `Hasher`, `HashInfo`                                             |
| `contracts/log`           | `Logger`                                                         |
| `contracts/mail`          | `Mailer`, `Mailable`, `Message`, `Envelope`, `Attachment`        |
| `contracts/notifications` | `Channel`, `Notifiable`, `Dispatcher`                            |
| `contracts/pagination`    | `Paginator`, `LengthAwarePaginator`, `CursorPaginator`, `Cursor` |
| `contracts/pipeline`      | `Pipeline`, `Pipe`                                               |
| `contracts/provider`      | `ServiceProvider`, `Bootable`, `Provides`                        |
| `contracts/validation`    | `Rule`, `Validator`, `MessageBag`                                |

## Top-level Contracts

The package root exposes framework-wide utilities:

### Clock

`Clock` abstracts wall-clock time so tests can freeze or advance the clock:

```go
type Clock interface {
    Now() time.Time
}
```

Inject a `Clock` into any component that needs timestamps:

```go
type SystemClock struct{}
func (SystemClock) Now() time.Time { return time.Now() }

type FakeClock struct{ t time.Time }
func (c *FakeClock) Now() time.Time { return c.t }
func (c *FakeClock) Advance(d time.Duration) { c.t = c.t.Add(d) }
```

## Using Contracts Instead of Concretes

When you write a service that needs to send mail, depend on
`contracts/mail.Mailer` — not on `mailx.Mailer`:

```go
import cmail "github.com/bedrock/packages/contracts/mail"

type WelcomeService struct {
    mailer cmail.Mailer // interface, not the concrete *mailx.Mailer
}

func (s *WelcomeService) SendWelcome(ctx context.Context, user *User) error {
    return s.mailer.Send(ctx, &WelcomeMailable{User: user})
}
```

This lets you swap `mailx` for any other implementation — or a test double —
without touching `WelcomeService`.

## Writing Test Doubles

Because everything is an interface, fakes are trivial:

```go
type FakeMailer struct {
    Sent []cmail.Mailable
}

func (f *FakeMailer) Send(ctx context.Context, m cmail.Mailable) error {
    f.Sent = append(f.Sent, m)
    return nil
}
```

See the `testing.md` concept page for a complete catalogue of Bedrock's
built-in fakes.
