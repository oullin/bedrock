# CSRF Protection

Bedrock protects against Cross-Site Request Forgery using per-session tokens
stored in the session and verified on state-changing requests.

## How It Works

1. On session start, a cryptographically random token is generated and stored
   as `_token` inside the session.
2. For every non-GET request, the CSRF middleware compares the `_token` from
   the session against either:
   - The `_token` form field, **or**
   - The `X-CSRF-TOKEN` header, **or**
   - The `X-XSRF-TOKEN` header (Base64-decoded)
3. On mismatch, the middleware responds with `419 Page Expired`.

## Accessing the Token

Inside a handler, read the token from the session store:

```go
store := session.FromRequest(r)
token := store.Token() // string
```

## Injecting the Token into Forms

```go
tmpl := template.Must(template.New("form").Parse(`
<form method="POST" action="/users">
    <input type="hidden" name="_token" value="{{.Token}}">
    <input name="email">
    <button type="submit">Create</button>
</form>
`))

tmpl.Execute(w, map[string]string{"Token": store.Token()})
```

## AJAX Requests

Send the token in the `X-CSRF-TOKEN` header:

```javascript
fetch('/api/users', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
        'X-CSRF-TOKEN': document.querySelector('meta[name="csrf-token"]').content,
    },
    body: JSON.stringify({ email: 'a@b.com' }),
})
```

Make the token available in a `<meta>` tag in your layout:

```html
<meta name="csrf-token" content="{{.Token}}">
```

## Excluding Routes

Some routes — webhooks, Stripe callbacks, Passport tokens — legitimately
cannot carry a CSRF token. Exclude them when wiring the middleware:

```go
csrf := session.NewCsrfMiddleware(store, session.CsrfOptions{
    Except: []string{
        "/webhooks/*",
        "/stripe/*",
        "/api/oauth/token",
    },
})

router.Use(csrf.Handle)
```

## Regenerating the Token

After login, regenerate the session ID _and_ the CSRF token to prevent session
fixation:

```go
store.Regenerate(true) // also rotates the CSRF token
```

## Token Size and Storage

- Tokens are 40-character URL-safe base64 strings derived from 30 random bytes.
- They live in the session store, so the storage guarantees of your session
  handler apply (file, database, Redis, encrypting cookie, …).
- Tokens are **not** rotated on every request — only on session regeneration.

## Double-Submit Cookie Pattern

For stateless APIs where a session doesn't fit, the `cookie` package can set
an `XSRF-TOKEN` cookie that the client echoes back in the `X-XSRF-TOKEN`
header. The server verifies equality. This is the mechanism used by Axios and
similar libraries out of the box.

## Testing

`session.ArrayHandler` lets tests boot a session without real storage. Set a
known token and include it in test requests:

```go
store := session.NewStore("test", session.NewArrayHandler(), nil)
store.Start()
store.PutToken("test-token")

req := httptest.NewRequest("POST", "/users", body)
req.Header.Set("X-CSRF-TOKEN", "test-token")
```
