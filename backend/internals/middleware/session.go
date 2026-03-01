package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	envConsts "clove/internals/consts/env"
	authservice "clove/internals/services/auth"
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

const cookieName = "euf_k15"

type SessionType string

const (
	OTT            SessionType = "OTT"
	SDKToken       SessionType = "SDK_TOKEN"
	RegularSession SessionType = "REGULAR_SESSION"
)

type Session struct {
	Permissions authservice.PermissionsBuilder
	SessionType SessionType
	UserID      pgtype.UUID
	AppID       pgtype.UUID
	ChannelID   string
}

// CredentialLookup is satisfied by any service that can fetch a credential by token.
type CredentialLookup interface {
	GetByToken(ctx context.Context, token string) (*repository.Credential, error)
}

type sessionCtxKey struct{}

// SessionFromContext retrieves the session stored in the context by SessionMiddleware or AuthMiddleware.
func SessionFromContext(ctx context.Context) (*Session, bool) {
	s, ok := ctx.Value(sessionCtxKey{}).(*Session)
	return s, ok
}

// UnAuthResponse writes a 401 Unauthorized response.
func UnAuthResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

// SetSessionCookie writes the session token as an HTTP-only cookie.
func SetSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Expires:  expires,
		HttpOnly: true,
		Secure:   envConsts.IsProd(),
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	})
}

// ClearSessionCookie expires the session cookie immediately.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   envConsts.IsProd(),
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	})
}

// AuthMiddleware validates the request credential and stores the resulting Session
// in the request context. Responds 401 if auth fails.
// Use this for routes that must always be authenticated (e.g. WebSocket).
func AuthMiddleware(creds CredentialLookup) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := ParseSession(r, creds)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), sessionCtxKey{}, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SessionMiddleware enriches the request context with a Session if the request
// carries valid credentials. Unlike AuthMiddleware, it never rejects requests —
// handlers are responsible for checking SessionFromContext.
func SessionMiddleware(creds CredentialLookup) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if sess, err := ParseSession(r, creds); err == nil {
				ctx := context.WithValue(r.Context(), sessionCtxKey{}, sess)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ParseSession extracts the token from the request and looks it up via creds.
// Returns an error if the token is missing, not found, or expired.
func ParseSession(r *http.Request, creds CredentialLookup) (*Session, error) {
	token := extractToken(r)
	if token == "" {
		return nil, errors.New("unauthorized")
	}

	cred, err := creds.GetByToken(r.Context(), token)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	if !cred.ExpiresAt.Valid || cred.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("unauthorized")
	}

	perms, err := authservice.ParsePermissions([]byte(cred.Permissions))
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	return &Session{
		Permissions: *perms,
		SessionType: credTypeToSessionType(cred.Type),
		UserID:      cred.UserID,
		AppID:       cred.AppID,
		ChannelID:   cred.ChannelID.String,
	}, nil
}

func credTypeToSessionType(t repository.CredentialType) SessionType {
	switch t {
	case repository.CredentialTypeOtt:
		return OTT
	case repository.CredentialTypeSdk:
		return SDKToken
	default:
		return RegularSession
	}
}

func extractToken(r *http.Request) string {
	if bearer := bearerFrom(r.Header.Get("Authorization")); bearer != "" {
		return bearer
	}
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func bearerFrom(header string) string {
	_, after, ok := strings.Cut(header, "Bearer ")
	if !ok {
		return ""
	}
	return after
}
