package auth

import (
	"clove/internals/auth/apiguard"
	"clove/internals/auth/tokenguard"
	"context"
	"net/http"
	"time"

	"github.com/dineshgowda24/browser"
)

type AuthMethod string

const api_key_header_name = "token"
const session_cookies_name = "session"

const (
	COOKIE   AuthMethod = "cookie"
	KEY      AuthMethod = "key"
	Unknowen AuthMethod = "Unknowen"
)

type AuthManager struct {
	authMethod AuthMethod
	r          *http.Request
}

func (a *AuthManager) IsLoggedIn() bool {
	return a.authMethod != Unknowen
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := ParseAuthFromRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		new_ctx := context.WithValue(r.Context(), "session", session)
		r.WithContext(new_ctx)
		next.ServeHTTP(w, r)
	})
}

func ParseAuthFromRequest(r *http.Request) (*tokenguard.SessionTokenClaim, error) {
	d, err := browser.NewBrowser(r.UserAgent())
	if err != nil {
		new_ctx := context.WithValue(r.Context(), "session", nil)
		r.WithContext(new_ctx)
	}

	if d.IsBrowserKnown() {
		cookie, err := r.Cookie(session_cookies_name)
		if err != nil {
			new_ctx := context.WithValue(r.Context(), "session", nil)
			r.WithContext(new_ctx)
		}
		if cookie.Expires.After(time.Now()) {
			new_ctx := context.WithValue(r.Context(), "session", nil)
			r.WithContext(new_ctx)
		}
		content := cookie.Value
		return tokenguard.ValidateSessionToken(content)

	} else {
		content := apiguard.GetHeaderApi(r)
		return tokenguard.ValidateSessionToken(content)

	}
}
