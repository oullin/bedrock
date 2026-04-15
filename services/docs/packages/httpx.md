# httpx

HTTP utilities, middleware, and testing helpers.

## Overview

The `httpx` package provides rich HTTP request and response primitives built
on Go's `net/http`, inspired by Upstream's HTTP layer.

**Module:** `github.com/gocanto/bedrock/packages/httpx`

```bash
go get github.com/gocanto/bedrock/packages/httpx@latest
```

## Request

`httpx.Request` wraps `*http.Request` and adds:

- **Input helpers** — typed access to query, form, and JSON body values
- **Content-type negotiation** — `WantsJSON()`, `ExpectsJSON()`, etc.
- **Flash data** — single-request carry-over values stored in the session
- **Precognitive** — marks a request as a form pre-validation probe

## Response Writers

| Type               | Description                                 |
| ------------------ | ------------------------------------------- |
| `Response`         | HTML/text responses with status and headers |
| `JsonResponse`     | Marshals any value to a JSON body           |
| `RedirectResponse` | 3xx redirects with optional flash data      |
| `StreamedEvent`    | Server-Sent Events (SSE) writer             |

## File Uploads

`httpx.UploadedFile` wraps `multipart.FileHeader` with helpers for
validation, path generation, and storage via pluggable backends.

## Usage

```go
// In an http.Handler — wrap the stdlib request
func MyHandler(w http.ResponseWriter, r *http.Request) {
    req := httpx.NewRequest(r)

    // Typed input access
    name  := req.Input("name")          // any
    email := req.QueryString("email")   // string
    page  := req.QueryInt("page", 1)    // int (with default)

    // Content negotiation
    if req.WantsJSON() {
        httpx.JSON(w, map[string]any{"hello": name})
        return
    }

    httpx.Text(w, "Hello, "+name)
}
```

### JSON Response

```go
res := httpx.NewJsonResponse(map[string]any{
    "users": users,
    "total": 42,
})
res.WithStatus(http.StatusOK).
    WithHeader("X-Request-ID", requestID).
    Send(w)
```

### Redirect

```go
httpx.NewRedirectResponse("/dashboard").
    WithStatus(http.StatusFound).
    With("status", "Login successful.").
    Send(w)
```

### File Uploads

```go
file, err := req.File("avatar")
path, err := file.StoreAs("/uploads/avatars", file.ClientName())
```

### Outbound HTTP Client

```go
client := httpx.NewClient().
    WithToken("my-api-token").
    AcceptJSON()

resp, err := client.Post("https://api.example.com/data", map[string]any{
    "key": "value",
})
data, err := resp.JSON()
```

## Sub-packages

| Sub-package        | Purpose                                                |
| ------------------ | ------------------------------------------------------ |
| `httpx/middleware` | Request logging, CORS, throttle, and more              |
| `httpx/client`     | Outbound HTTP client with fluent API and test fakes    |
| `httpx/resources`  | JSON API resource and resource collection transformers |
| `httpx/testing`    | `TestRequest` / `TestResponse` helpers                 |
