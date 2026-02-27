package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionTypeConstants(t *testing.T) {
	assert.Equal(t, SessionType("OTT"), OTT)
	assert.Equal(t, SessionType("SDK_TOKEN"), SDKToken)
	assert.Equal(t, SessionType("REGULAR_SESSION"), RegularSession)
}

func TestSession_Structure(t *testing.T) {
	perms := NewPermissionsBuilder()
	perms.Allow(APP, CREATE)

	session := Session{
		Permissions: perms,
		SessionType: RegularSession,
	}

	assert.Equal(t, RegularSession, session.SessionType)
	assert.True(t, session.Permissions.Can(APP, CREATE))
}

func TestUnAuthResponse(t *testing.T) {
	w := httptest.NewRecorder()
	UnAuthResponse(w)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Unauthorized(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// nil DB is safe here — ParseSession returns early when no token is present
	middleware := AuthMiddleware(nil)(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestParseSession_NoValidTokens(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// nil DB is safe — returns early when token string is empty
	session, err := ParseSession(req, nil)

	assert.Nil(t, session)
	assert.Error(t, err)
	assert.Equal(t, "unauthorized", err.Error())
}
