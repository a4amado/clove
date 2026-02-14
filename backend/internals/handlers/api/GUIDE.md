# API Layer Architecture Guide

This document explains the API layer after migrating from raw `chi` handlers to `swaggest/rest`. If you're new to the codebase, read this first.

---

## Table of Contents

1. [The Big Picture](#the-big-picture)
2. [Why We Migrated](#why-we-migrated)
3. [Key Packages](#key-packages)
4. [How a Request Flows](#how-a-request-flows)
5. [Writing a New Handler](#writing-a-new-handler)
6. [Input Struct Tags](#input-struct-tags)
7. [Error Handling](#error-handling)
8. [The httpctx Package](#the-httpctx-package)
9. [Raw Handlers (WebSocket & Entry)](#raw-handlers-websocket--entry)
10. [Registering Routes](#registering-routes)
11. [OpenAPI Docs](#openapi-docs)

---

## The Big Picture

```
HTTP Request
    │
    ▼
web.Service (replaces chi.Router)
    │
    ├── httpctx.Middleware   ← stores request/response in context
    │
    ├── usecase handlers    ← typed input/output, auto-documented
    │   (signup, signin, keys, regions, tokens, etc.)
    │
    └── raw handlers        ← plain http.HandlerFunc
        (websocket, entry)
```

Before the migration, every handler was a plain `http.HandlerFunc`:

```go
// OLD — you manually decode JSON, parse path params, write responses
func CreateApp(w http.ResponseWriter, r *http.Request) {
    body := CreateAppStruct{}
    json.NewDecoder(r.Body).Decode(&body)
    appId := r.PathValue("app_id")
    // ... do work ...
    json.NewEncoder(w).Encode(result)
}
```

After the migration, most handlers are **usecase interactors**:

```go
// NEW — the framework decodes input and encodes output for you
func CreateApp() usecase.Interactor {
    return usecase.NewInteractor(func(ctx context.Context, input CreateAppInput, output *CreateAppOutput) error {
        // input is already decoded from JSON body + URL path + query string
        // just do your business logic and set output fields
        output.App = *app
        return nil  // framework writes the JSON response
    })
}
```

---

## Why We Migrated

1. **Auto-generated OpenAPI docs** — The framework reads your Go structs and produces an OpenAPI 3.1 spec automatically. No need to write or maintain a separate spec file.

2. **Less boilerplate** — No more manual `json.NewDecoder`, `r.PathValue`, `strconv.ParseInt`, or `json.NewEncoder` in every handler.

3. **Type safety** — Path params, query params, and JSON bodies are decoded into typed structs. If `app_id` isn't a valid UUID, the framework returns 400 before your code runs.

4. **Swagger UI** — Visit `/docs` in the browser to see and test every endpoint interactively.

---

## Key Packages

| Package | What it is | Why we use it |
|---------|-----------|---------------|
| `github.com/swaggest/rest/web` | Wraps chi router with OpenAPI support | Creates the service, registers routes |
| `github.com/swaggest/usecase` | Defines typed request/response handlers | `usecase.NewInteractor` creates handlers with typed I/O |
| `github.com/swaggest/openapi-go/openapi31` | OpenAPI 3.1 reflector | Generates the spec from Go structs |
| `github.com/swaggest/swgui/v5cdn` | Swagger UI served from CDN | Renders interactive docs at `/docs` |
| `httpctx` (our own, see below) | Stores `http.Request` and `http.ResponseWriter` in context | Lets interactors access cookies, headers, etc. |

---

## How a Request Flows

Here's what happens when someone calls `POST /v1/apps/{app_id}/keys`:

```
1. HTTP request arrives
        │
2. httpctx.Middleware
   └── Stores *http.Request and http.ResponseWriter in ctx
        │
3. swaggest/rest request decoder
   └── Reads the CreateKeyInput struct tags:
       - path:"app_id"  → extracts from URL, parses as uuid.UUID
       - json:"name"    → extracts from JSON body
   └── If decoding fails → automatic 400 response (your code never runs)
        │
4. Your interactor function runs
   └── func(ctx, input CreateKeyInput, output *CreateKeyOutput) error
   └── input.AppID and input.Name are ready to use
   └── On error: return &apperrors.AppError{...}
   └── On success: set output fields and return nil
        │
5. swaggest/rest response encoder
   └── If you returned nil  → JSON-encodes *output and writes 200
   └── If you returned error → uses HTTPStatus() for status code,
                                JSON-encodes the error as the body
```

---

## Writing a New Handler

Let's say you need a `GET /v1/apps/{app_id}/settings?include_secrets=true` endpoint.

### Step 1: Define Input and Output structs

```go
// In a file like app/settings/get.go
package AppSettingsHandlersV1

type GetSettingsInput struct {
    AppID          uuid.UUID `path:"app_id"`           // from URL
    IncludeSecrets bool      `query:"include_secrets"`  // from query string
}

type GetSettingsOutput struct {
    Theme    string `json:"theme"`
    Language string `json:"language"`
}
```

The struct tags tell the framework WHERE to find each value:
- `path:"app_id"` → from the URL path `/v1/apps/{app_id}/settings`
- `query:"include_secrets"` → from `?include_secrets=true`
- `json:"theme"` → goes into the JSON response body

### Step 2: Write the interactor

```go
func GetSettings() usecase.Interactor {
    u := usecase.NewInteractor(func(ctx context.Context, input GetSettingsInput, output *GetSettingsOutput) error {
        // Auth check (see "The httpctx Package" section below)
        r := httpctx.Request(ctx)
        session, err := auth.ParseSessionFromRequest(r)
        if err != nil {
            return &apperrors.AppError{
                StatusCode: http.StatusUnauthorized,
                Code:       "UNAUTHORIZED",
            }
        }

        // Your business logic
        settings, err := getSettingsFromDB(ctx, input.AppID)
        if err != nil {
            return &apperrors.AppError{
                StatusCode: http.StatusInternalServerError,
            }
        }

        // Set output — the framework JSON-encodes this as the response
        output.Theme = settings.Theme
        output.Language = settings.Language
        return nil
    })

    // These show up in the OpenAPI docs
    u.SetTitle("Get App Settings")
    u.SetTags("Settings")

    return u
}
```

### Step 3: Register the route

In `routes.go`:

```go
service.Get("/v1/apps/{app_id}/settings", AppSettingsHandlersV1.GetSettings())
```

That's it. The framework handles JSON decoding, path param parsing, response encoding, and OpenAPI documentation.

---

## Input Struct Tags

These tags control how the framework decodes the request:

| Tag | Source | Example |
|-----|--------|---------|
| `path:"name"` | URL path parameter `{name}` | `AppID uuid.UUID \`path:"app_id"\`` |
| `query:"name"` | Query string `?name=value` | `Page int64 \`query:"page_idx"\`` |
| `json:"name"` | JSON request body | `Email string \`json:"email"\`` |
| `header:"name"` | HTTP header | `Token string \`header:"Authorization"\`` |
| `cookie:"name"` | HTTP cookie | `Session string \`cookie:"session"\`` |

You can mix them in one struct. Fields with `path` or `query` tags come from the URL. Fields with only `json` tags come from the request body.

```go
type CreateKeyInput struct {
    AppID uuid.UUID `path:"app_id"`  // from URL: /v1/apps/{app_id}/keys
    Name  string    `json:"name"`    // from JSON body: {"name": "my-key"}
}
```

**Important:** If a field has no tag, the framework ignores it during decoding. Always add the appropriate tag.

---

## Error Handling

When something goes wrong, return an `*apperrors.AppError`:

```go
return &apperrors.AppError{
    Code:       "ERROR_NOT_FOUND",       // machine-readable error code
    Message:    "App not found",          // human-readable message
    StatusCode: http.StatusNotFound,      // HTTP status code (404)
}
```

`AppError` implements two interfaces that the framework uses:

```go
// error interface — makes it a valid Go error
func (e *AppError) Error() string { ... }

// Used by swaggest/rest to set the HTTP status code
func (e *AppError) HTTPStatus() int { return e.StatusCode }
```

When you return an `AppError`, the framework:
1. Calls `HTTPStatus()` to set the HTTP response status (e.g. 404)
2. JSON-encodes the error as the response body

When you return `nil`, the framework:
1. JSON-encodes your `*output` struct as the response body with status 200

---

## The httpctx Package

**Location:** `handlers/api/httpctx/httpctx.go`

### The problem it solves

In a usecase interactor, you only get `context.Context` — not `http.Request` or `http.ResponseWriter`. But you still need them for:

- **Auth parsing** — `auth.ParseSessionFromRequest(r)` reads cookies and headers from `*http.Request`
- **Setting cookies** — `http.SetCookie(w, cookie)` writes to `http.ResponseWriter`

### How it works

It's a middleware + two helper functions. That's the entire package.

**The middleware** (runs on every request, registered in `index.go`):
```go
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Stash w and r into the context
        ctx := context.WithValue(r.Context(), reqKey{}, r)
        ctx = context.WithValue(ctx, rwKey{}, w)
        // Continue with the enriched context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**The helpers** (used inside interactors):
```go
r := httpctx.Request(ctx)          // get *http.Request back
w := httpctx.ResponseWriter(ctx)   // get http.ResponseWriter back
```

### When to use it

Use `httpctx.Request(ctx)` whenever you need to read cookies, headers, or anything from the raw request that isn't already in your Input struct:

```go
func(ctx context.Context, input MyInput, output *MyOutput) error {
    r := httpctx.Request(ctx)
    session, err := auth.ParseSessionFromRequest(r)
    // ...
}
```

Use `httpctx.ResponseWriter(ctx)` only when you need to set response headers before the framework writes the body (e.g. setting cookies in the signup handler).

**You do NOT need httpctx for:**
- Path params → use `path` tag on input struct
- Query params → use `query` tag on input struct
- JSON body → use `json` tag on input struct
- Writing JSON responses → just set output fields and return nil

---

## Raw Handlers (WebSocket & Entry)

Two endpoints are NOT converted to usecase interactors:

1. **`/v1/apps/{app_id}/ws`** — WebSocket connection. Uses `gorilla/websocket` to upgrade the HTTP connection. The usecase pattern doesn't work here because WebSocket is a long-lived bidirectional connection, not a request/response cycle.

2. **`/v1/apps/{app_id}/entry`** — Message entry. Reads a raw binary body (not JSON) with size limits. The usecase pattern expects JSON input/output, which doesn't fit raw body handling.

These are registered differently in `routes.go`:

```go
// Usecase handler — framework decodes/encodes for you
service.Post("/v1/apps/{app_id}/keys", AppKeysHandlersV1.CreateAppApiKey())

// Raw handler — you handle everything yourself, just like before
service.Method(http.MethodGet, "/v1/apps/{app_id}/ws",
    auth.AuthMiddleware(http.HandlerFunc(AppHandlersV1.UserConnect)))
```

For raw handlers, `auth.AuthMiddleware` is wrapped around them since they don't go through the usecase flow.

---

## Registering Routes

All routes are registered in two files:

### `handlers/api/index.go` — Creates the service

```go
func NewService() *web.Service {
    service := web.NewService(openapi31.NewReflector())  // create service
    service.Use(httpctx.Middleware)                        // global middleware
    v1.V1Routes(service)                                  // register routes
    service.Docs("/docs", swgui.New)                      // serve docs UI
    return service
}
```

### `handlers/api/v1/routes.go` — Registers all v1 routes

Use the right method for the right handler type:

```go
// For usecase interactors (most handlers):
service.Get("/path", handler())
service.Post("/path", handler())
service.Delete("/path", handler())
service.Patch("/path", handler())

// For raw http.HandlerFunc (websocket, entry):
service.Method(http.MethodGet, "/path", http.HandlerFunc(rawHandler))
```

---

## OpenAPI Docs

The OpenAPI spec is generated automatically from your Input/Output structs. Visit `/docs` in the browser to see the Swagger UI.

Every usecase handler you register appears in the docs with:
- Request parameters (from `path`, `query` tags)
- Request body schema (from `json` tags)
- Response body schema (from output struct `json` tags)
- Tags and title (from `SetTags` and `SetTitle`)

You don't need to write any YAML or JSON spec files. Just define your Go structs with the right tags and the docs update themselves.
