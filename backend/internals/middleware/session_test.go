package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	authservice "clove/internals/services/auth"

	"github.com/stretchr/testify/assert"
)

func TestSessionTypeConstants(t *testing.T) {
	assert.Equal(t, SessionType("OTT"), OTT)
	assert.Equal(t, SessionType("SDK_TOKEN"), SDKToken)
	assert.Equal(t, SessionType("REGULAR_SESSION"), RegularSession)
}

func TestSession_Structure(t *testing.T) {
	perms := authservice.NewPermissionsBuilder()
	perms.Allow(authservice.APP, authservice.CREATE)

	sess := Session{
		Permissions: perms,
		SessionType: RegularSession,
	}

	assert.Equal(t, RegularSession, sess.SessionType)
	assert.True(t, sess.Permissions.Can(authservice.APP, authservice.CREATE))
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

	// nil is safe here — ParseSession returns early when no token is present
	mw := AuthMiddleware(nil)(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	mw.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestParseSession_NoValidTokens(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// nil is safe — returns early when token string is empty
	sess, err := ParseSession(req, nil)

	assert.Nil(t, sess)
	assert.Error(t, err)
	assert.Equal(t, "unauthorized", err.Error())
}
