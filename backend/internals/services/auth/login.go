package authservice

import (
	"errors"
	"time"

	repository "clove/internals/services/generatedRepo"

	"github.com/google/uuid"
)

type LoginParams struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token     string
	ExpiresAt time.Time
}

// Login looks up the user by email, verifies the password, creates a session
// credential, and returns the token with its expiry.
func (s *AuthService) Login(params LoginParams) (*LoginResult, error) {
	emailRecord, err := s.DB.UserEmail_SelectByEmail(s.GetCtx(), params.Email)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	user, err := s.DB.User_Select(s.GetCtx(), emailRecord.UserID)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	if user.Hash != params.Password {
		return nil, errors.New("unauthorized")
	}

	token, err := GenRandKey(32)
	if err != nil {
		return nil, err
	}

	perms := NewPermissionsBuilder()
	perms.Allow(RESOURCES_ALL, OPERATIONS_ALL)
	permsStr, err := perms.String()
	if err != nil {
		return nil, err
	}

	expires := time.Now().Add(time.Hour * 24 * 7)
	_, err = s.Create(InsertParams{
		Token:       token,
		UserID:      uuid.UUID(emailRecord.UserID.Bytes),
		Type:        repository.CredentialTypeSession,
		Permissions: permsStr,
		ExpiresAt:   expires,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: token, ExpiresAt: expires}, nil
}
