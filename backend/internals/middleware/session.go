package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	authorization          = "Authorization"
	sessionTokenCookieName = "euf_k15"
)

type SessionType string

const (
	OTT            SessionType = "OTT"
	SDKToken       SessionType = "SDK_TOKEN"
	RegularSession SessionType = "REGULAR_SESSION"
)

type Session struct {
	Permissions PermissionsBuilder
	SessionType SessionType
	UserID      pgtype.UUID
	AppID       pgtype.UUID
	ChannelID   string
}

type sessionCtxKey struct{}

// AuthMiddleware validates the request credential against the DB and stores
// the resulting Session in the request context.
func AuthMiddleware(db repository.DBTX) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := ParseSession(r, db)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), sessionCtxKey{}, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SessionFromContext(ctx context.Context) (*Session, bool) {
	s, ok := ctx.Value(sessionCtxKey{}).(*Session)
	return s, ok
}

func UnAuthResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

// SetSessionCookie writes the session token as an HTTP-only cookie.
func SetSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionTokenCookieName,
		Value:    token,
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})
}

// ParseSession extracts the token from the request and looks it up in the
// credential table. Returns an error if the token is missing, not found, or expired.
func ParseSession(r *http.Request, db repository.DBTX) (*Session, error) {
	token := extractToken(r)
	if token == "" {
		return nil, errors.New("unauthorized")
	}

	cred, err := repository.New(db).Credential_SelectByToken(r.Context(), token)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	if !cred.ExpiresAt.Valid || cred.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("unauthorized")
	}

	perms, err := ParsePermissions([]byte(cred.Permissions))
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

func extractToken(r *http.Request) string {
	if bearer := extractBearerTokenFrom(r.Header.Get(authorization)); bearer != "" {
		return bearer
	}
	if cookie, err := r.Cookie(sessionTokenCookieName); err == nil {
		return cookie.Value
	}
	return ""
}

func extractBearerTokenFrom(header string) string {
	_, after, ok := strings.Cut(header, "Bearer ")
	if !ok {
		return ""
	}
	return after
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
