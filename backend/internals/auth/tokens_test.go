package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	repository "clove/internals/services/generatedRepo"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEnv(t *testing.T) {
	t.Helper()
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-1234567890")
	os.Setenv("REGION", "dk1")
}

func TestExtractBearerTokenFrom(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"valid bearer token", "Bearer abc123", "abc123"},
		{"bearer with spaces", "Bearer my-token-value", "my-token-value"},
		{"no bearer prefix", "abc123", ""},
		{"empty string", "", ""},
		{"just Bearer", "Bearer ", ""},
		{"lowercase bearer", "bearer abc123", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractBearerTokenFrom(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateAndValidateSDKToken(t *testing.T) {
	setupTestEnv(t)

	appID := uuid.New()

	token, err := GenerateSDKToken(appID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	claims, err := ValidateSDKToken(token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	assert.Equal(t, appID, claims.AppID)
	assert.True(t, claims.Permessions.Can(DELIVERY, CREATE))
	assert.True(t, claims.Permessions.Can(OneTimeToken, CREATE))
}

func TestValidateSDKToken_InvalidToken(t *testing.T) {
	setupTestEnv(t)

	claims, err := ValidateSDKToken("invalid-token")
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateSDKToken_TamperedToken(t *testing.T) {
	setupTestEnv(t)

	appID := uuid.New()

	token, err := GenerateSDKToken(appID)
	require.NoError(t, err)

	// Tamper with the token
	tamperedToken := token[:len(token)-5] + "xxxxx"

	claims, err := ValidateSDKToken(tamperedToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestParseSDKTokenFromRequest(t *testing.T) {
	setupTestEnv(t)

	appID := uuid.New()

	token, err := GenerateSDKToken(appID)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	claims, err := ParseSDKTokenFromRequest(req)
	require.NoError(t, err)
	require.NotNil(t, claims)

	assert.Equal(t, appID, claims.AppID)
}

func TestParseSDKTokenFromRequest_NoHeader(t *testing.T) {
	setupTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	claims, err := ParseSDKTokenFromRequest(req)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGenerateAndValidateSessionToken(t *testing.T) {
	setupTestEnv(t)

	user := repository.User{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Hash: "test-hash",
		Role: repository.NullUserRole{
			UserRole: repository.UserRoleUser,
			Valid:    true,
		},
	}

	token, err := GenerateSessionToken(user)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	claims, err := ValidateSessionToken(token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	assert.Equal(t, user.Hash, claims.User.Hash)
	// Session tokens have full permissions
	assert.True(t, claims.Permessions.Can(RESOURCES_ALL, OPERATIONS_ALL))
}

func TestValidateSessionToken_InvalidToken(t *testing.T) {
	setupTestEnv(t)

	claims, err := ValidateSessionToken("invalid-token")
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestParseSessionFromRequest_Cookie(t *testing.T) {
	setupTestEnv(t)

	user := repository.User{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Hash: "test-hash",
	}

	token, err := GenerateSessionToken(user)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  session_token_cookie_name,
		Value: token,
	})

	claims, err := ParseSessionFromRequest(req)
	require.NoError(t, err)
	require.NotNil(t, claims)
}

func TestParseSessionFromRequest_Mobile(t *testing.T) {
	setupTestEnv(t)

	user := repository.User{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Hash: "test-hash",
	}

	token, err := GenerateSessionToken(user)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Client-Type", "mobile")
	req.Header.Set("Authorization", "Bearer "+token)

	claims, err := ParseSessionFromRequest(req)
	require.NoError(t, err)
	require.NotNil(t, claims)
}

func TestParseSessionFromRequest_NoCookie(t *testing.T) {
	setupTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	claims, err := ParseSessionFromRequest(req)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGenerateAndValidateOneTimeToken(t *testing.T) {
	setupTestEnv(t)

	app := repository.App{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		AppSlug: "test-app",
		AppType: repository.AppTypeFree,
		Region:  []repository.Region{repository.RegionDk1},
	}
	channelID := "test-channel"

	token, err := GenerateOneTimeToken(app, channelID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	claims, err := ValidateOneTimeToken(token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	assert.Equal(t, app.AppSlug, claims.App.AppSlug)
	assert.Equal(t, channelID, claims.ChannelID)
	assert.True(t, claims.Permessions.Can(DELIVERY, READ))
}

func TestValidateOneTimeToken_InvalidToken(t *testing.T) {
	setupTestEnv(t)

	claims, err := ValidateOneTimeToken("invalid-token")
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestParseOneTimeTokenFromRequest(t *testing.T) {
	setupTestEnv(t)

	app := repository.App{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		AppSlug: "test-app",
		AppType: repository.AppTypeFree,
	}
	channelID := "test-channel"

	token, err := GenerateOneTimeToken(app, channelID)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	claims, err := ParseOneTimeTokenFromRequest(req)
	require.NoError(t, err)
	require.NotNil(t, claims)

	assert.Equal(t, channelID, claims.ChannelID)
}

func TestParseOneTimeTokenFromRequest_NoHeader(t *testing.T) {
	setupTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	claims, err := ParseOneTimeTokenFromRequest(req)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
