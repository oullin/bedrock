# URL Generation

<!-- laravel-docs: urls.md#url-generation -->
<!-- laravel-docs: urls.md#urls-for-named-routes -->
<!-- laravel-docs: urls.md#default-values -->

Generating URLs for routes is a first-class feature of the router. Prefer
named-route lookups over hand-rolled string concatenation — renaming a path
in one place should not break every link in your templates.

## Generating a Named Route URL

```go
router.Get("/users/{id}", showUser).Name("users.show")

// Later...
url := router.URL("users.show", map[string]any{"id": 42})
// "/users/42"
```

## Required vs Optional Parameters

If the route pattern declares a required parameter that you don't provide,
`URL` returns an error:

```go
_, err := router.URL("users.show", nil) // missing "id"
// err != nil
```

Optional parameters (`/posts/{page?}`) can be omitted:

```go
router.URL("posts.index", nil)                           // "/posts"
router.URL("posts.index", map[string]any{"page": 3})     // "/posts/3"
```

## Query Strings

Extra values that don't match a route parameter are appended as query string
parameters:

```go
router.URL("users.index", map[string]any{
    "sort":   "name",
    "filter": "active",
})
// "/users?filter=active&sort=name"
```

## Absolute URLs

By default, `URL` returns a path. Use `AbsoluteURL` when you need a
fully-qualified URL (for emails, webhooks, redirects to external tools):

```go
url := router.AbsoluteURL("users.show", map[string]any{"id": 42})
// "https://example.com/users/42"
```

Configure the base URL via the router's `SetBaseURL`:

```go
router.SetBaseURL("https://example.com")
```

## Signed URLs

For password reset links, email verification, and similar short-lived actions,
generate a signed URL that includes an HMAC signature and optional expiry:

```go
signed := router.SignedURL("verify.email", map[string]any{
    "user": user.ID,
}, time.Now().Add(24*time.Hour))
// "/verify-email/42?expires=1740000000&signature=abcdef..."
```

Verify the signature on the receiving handler:

```go
if err := router.VerifySignedURL(r); err != nil {
    http.Error(w, "Invalid or expired link", http.StatusUnauthorized)
    return
}
```

## Current URL Helpers

Inside a handler, introspect the current request:

```go
req := httpx.NewRequest(r)

req.URL()           // "/users/42?tab=settings"
req.FullURL()       // "https://example.com/users/42?tab=settings"
req.Path()          // "/users/42"
req.Host()          // "example.com"
```

## URL Defaults

Set default values that apply to every URL generation call (e.g. current
locale, current tenant):

```go
router.Defaults(map[string]any{
    "locale": "en",
    "tenant": currentTenantID,
})

// Now any route with {locale} is pre-populated:
router.URL("users.show", map[string]any{"id": 42})
// "/en/users/42"
```

## Redirects Using Named Routes

Prefer redirecting by name, not by path:

```go
httpx.NewRedirectResponse(router.URL("users.show", map[string]any{"id": user.ID})).
    Send(w)

// Or with a helper
routing.RedirectToRoute(w, r, router, "users.show", map[string]any{"id": user.ID})
```
