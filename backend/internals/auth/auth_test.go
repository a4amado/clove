package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionTypeConstants(t *testing.T) {
	assert.Equal(t, SessionType("OOT"), OOT)
	assert.Equal(t, SessionType("SDK_TOKEN"), SDKToken)
	assert.Equal(t, SessionType("REGULAR_SESSION"), RegualrSession)
}

func TestSession_Structure(t *testing.T) {
	perms := NewPermissionsBuilder()
	perms.Allow(APP, CREATE)

	session := Session{
		Permissions: perms,
		SessionType: RegualrSession,
	}

	assert.Equal(t, RegualrSession, session.SessionType)
	assert.True(t, session.Permissions.Can(APP, CREATE))
}

func TestUnAuthResponse(t *testing.T) {
	w := httptest.NewRecorder()

	UnAuthResponse(w)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Unauthorized(t *testing.T) {
	// Create a simple handler to test middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := AuthMiddleware(nextHandler)

	// Create request without any auth
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	// Should return 401 since no valid session
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestParseSession_NoValidTokens(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	session, err := ParseSession(req)

	assert.Nil(t, session)
	assert.Error(t, err)
	assert.Equal(t, "unauthrized", err.Error())
}
