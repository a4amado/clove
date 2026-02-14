package userservice

import (
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUsers_Structure(t *testing.T) {
	users := &Users{}
	assert.Nil(t, users.BaseService)
}

func TestUsers_WithBaseService(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)

	users := &Users{
		BaseService: baseService,
	}

	assert.NotNil(t, users.BaseService)
	assert.Equal(t, ctx, users.GetCtx())
	assert.False(t, users.IsCache())
}

func TestInsertUserParams_Structure(t *testing.T) {
	params := InsertUserParams{
		Email:           "test@example.com",
		Password:        "securepassword123",
		EmailVerifycode: "ABC123",
	}

	assert.Equal(t, "test@example.com", params.Email)
	assert.Equal(t, "securepassword123", params.Password)
	assert.Equal(t, "ABC123", params.EmailVerifycode)
}

func TestInsertUserParams_EmptyFields(t *testing.T) {
	params := InsertUserParams{}

	assert.Empty(t, params.Email)
	assert.Empty(t, params.Password)
	assert.Empty(t, params.EmailVerifycode)
}

func TestUserEmail_InsertParams_Conversion(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)

	userID := uuid.New()
	insertParams := InsertUserParams{
		Email:           "test@example.com",
		Password:        "password123",
		EmailVerifycode: "VERIFY123",
	}

	// Simulate the conversion for UserEmail_Insert
	emailParams := repository.UserEmail_InsertParams{
		Email:  insertParams.Email,
		UserID: baseService.ToPgUUID(userID),
		Code:   insertParams.EmailVerifycode,
	}

	assert.Equal(t, "test@example.com", emailParams.Email)
	assert.True(t, emailParams.UserID.Valid)
	assert.Equal(t, [16]byte(userID), emailParams.UserID.Bytes)
	assert.Equal(t, "VERIFY123", emailParams.Code)
}

func TestEmails_Structure(t *testing.T) {
	emails := &Emails{}
	assert.Nil(t, emails.Users)
}

func TestEmails_WithUsers(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)
	users := &Users{BaseService: baseService}

	emails := &Emails{
		Users: users,
	}

	assert.NotNil(t, emails.Users)
	assert.Equal(t, ctx, emails.GetCtx())
}

func TestEmails_InheritsFromUsers(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, true)
	users := &Users{BaseService: baseService}

	emails := &Emails{
		Users: users,
	}

	// Emails should inherit context and cache settings from Users
	assert.Equal(t, ctx, emails.GetCtx())
	assert.True(t, emails.IsCache())
}

func TestUserRole_Valid(t *testing.T) {
	tests := []struct {
		name     string
		role     repository.UserRole
		expected bool
	}{
		{"super is valid", repository.UserRoleSuper, true},
		{"admin is valid", repository.UserRoleAdmin, true},
		{"user is valid", repository.UserRoleUser, true},
		{"empty is invalid", repository.UserRole(""), false},
		{"unknown is invalid", repository.UserRole("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.Valid())
		})
	}
}

func TestAllUserRoleValues(t *testing.T) {
	roles := repository.AllUserRoleValues()

	assert.Len(t, roles, 3)
	assert.Contains(t, roles, repository.UserRoleSuper)
	assert.Contains(t, roles, repository.UserRoleAdmin)
	assert.Contains(t, roles, repository.UserRoleUser)
}

func TestAppType_Valid(t *testing.T) {
	tests := []struct {
		name     string
		appType  repository.AppType
		expected bool
	}{
		{"free is valid", repository.AppTypeFree, true},
		{"standard is valid", repository.AppTypeStandard, true},
		{"pro is valid", repository.AppTypePro, true},
		{"empty is invalid", repository.AppType(""), false},
		{"unknown is invalid", repository.AppType("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.appType.Valid())
		})
	}
}

func TestAllAppTypeValues(t *testing.T) {
	appTypes := repository.AllAppTypeValues()

	assert.Len(t, appTypes, 3)
	assert.Contains(t, appTypes, repository.AppTypeFree)
	assert.Contains(t, appTypes, repository.AppTypeStandard)
	assert.Contains(t, appTypes, repository.AppTypePro)
}
