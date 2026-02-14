package httpctx

import (
	"context"
	"net/http"
)

type reqKey struct{}
type rwKey struct{}

// Middleware stores the *http.Request and http.ResponseWriter in context so
// usecase interactors can access headers, cookies, etc.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), reqKey{}, r)
		ctx = context.WithValue(ctx, rwKey{}, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Request retrieves the *http.Request stored by Middleware.
func Request(ctx context.Context) *http.Request {
	r, _ := ctx.Value(reqKey{}).(*http.Request)
	return r
}

// ResponseWriter retrieves the http.ResponseWriter stored by Middleware.
func ResponseWriter(ctx context.Context) http.ResponseWriter {
	w, _ := ctx.Value(rwKey{}).(http.ResponseWriter)
	return w
}
