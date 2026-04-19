# mailx

Driver-based email sending.

## Overview

The `mailx` package provides a `Manager` that creates named `Mailer` instances.
Each mailer delegates to a swappable `Transport` for delivery. Lifecycle events
are fired before and after each send.

**Module:** `github.com/gocanto/bedrock/packages/mailx`

```bash
go get github.com/gocanto/bedrock/packages/mailx@latest
```

## Transports

| Transport        | Description                                   |
| ---------------- | --------------------------------------------- |
| `SmtpTransport`  | Production SMTP delivery                      |
| `LogTransport`   | Writes message to a logger (development)      |
| `ArrayTransport` | Stores messages in memory (testing / preview) |

## Creating a Mailer

```go
transport := mailx.NewSmtpTransport(mailx.SmtpConfig{
    Host:     "smtp.example.com",
    Port:     587,
    Username: "user@example.com",
    Password: "secret",
})

mailer := mailx.NewMailer("default", transport)
```

## Sending Mail

Implement `mailx.Mailable` to describe a message:

```go
type WelcomeMail struct {
    User *User
}

func (m *WelcomeMail) Envelope() *cmail.Envelope {
    return &cmail.Envelope{
        Subject: "Welcome to Bedrock!",
        To:      []cmail.Address{{Email: m.User.Email, Name: m.User.Name}},
    }
}

func (m *WelcomeMail) Content() *cmail.Content {
    return &cmail.Content{HtmlBody: "<h1>Welcome</h1>"}
}

func (m *WelcomeMail) Attachments() []cmail.Attachment { return nil }
```

```go
err := mailer.Send(ctx, &WelcomeMail{User: user})
```

## Global Overrides

```go
mailer.AlwaysFrom("noreply@example.com", "Bedrock")
mailer.AlwaysReplyTo("support@example.com")

// Force all mail to a single address (useful in staging)
mailer.AlwaysTo("dev@example.com")
```

## Using the Manager

```go
manager := mailx.NewManager(map[string]mailx.MailerFactory{
    "default": func() (*mailx.Mailer, error) {
        return mailx.NewMailer("default", smtpTransport), nil
    },
    "log": func() (*mailx.Mailer, error) {
        return mailx.NewMailer("log", mailx.NewLogTransport(logger)), nil
    },
})

m, err := manager.Mailer("log")
m.Send(ctx, mailable)
```

## Testing

Use `ArrayTransport` to capture sent messages in tests:

```go
arr := mailx.NewArrayTransport()
mailer := mailx.NewMailer("test", arr)
mailer.Send(ctx, &WelcomeMail{})

sent := arr.Messages() // []*mailx.SentMessage
```
