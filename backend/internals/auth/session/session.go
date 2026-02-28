package session

import (
	envConsts "clove/internals/consts/env"
	"net/http"
	"strings"
	"time"
)

const CookieName = "euf_k15"

const authorizationHeader = "Authorization"

// Set writes the session token as an HTTP-only cookie.
// Secure is set automatically in production.
func Set(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Expires:  expires,
		HttpOnly: true,
		Secure:   envConsts.IsProd(),
		SameSite: http.SameSiteDefaultMode,
	})
}

// Get returns the session token from the cookie, or empty string if absent.
func Get(r *http.Request) string {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// ExtractToken returns the bearer token from the Authorization header,
// falling back to the session cookie.
func ExtractToken(r *http.Request) string {
	if bearer := bearerFrom(r.Header.Get(authorizationHeader)); bearer != "" {
		return bearer
	}
	return Get(r)
}

func bearerFrom(header string) string {
	_, after, ok := strings.Cut(header, "Bearer ")
	if !ok {
		return ""
	}
	return after
}
