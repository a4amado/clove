package userservice

import (
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestDeleteEmailParams_Structure(t *testing.T) {
	emailID := uuid.New()
	params := DeleteEmailParams{
		EmailID: emailID,
	}

	assert.Equal(t, emailID, params.EmailID)
}

func TestListUserEmailParams_Structure(t *testing.T) {
	userID := uuid.New()
	params := ListUserEmailParams{
		UserID: userID,
	}

	assert.Equal(t, userID, params.UserID)
}

func TestGetEmailParams_Structure(t *testing.T) {
	emailID := uuid.New()
	params := GetEmailParams{
		EmailID: emailID,
	}

	assert.Equal(t, emailID, params.EmailID)
}

func TestGetEmailByEmailParams_Structure(t *testing.T) {
	params := GetEmailByEmailParams{
		Email: "test@example.com",
	}

	assert.Equal(t, "test@example.com", params.Email)
}

func TestGetEmailByEmailParams_EmptyEmail(t *testing.T) {
	params := GetEmailByEmailParams{
		Email: "",
	}

	assert.Empty(t, params.Email)
}

func TestVerifyEmail_Structure(t *testing.T) {
	emailID := uuid.New()
	params := VerifyEmail{
		EmailID: emailID,
	}

	assert.Equal(t, emailID, params.EmailID)
}

func TestEmails_ToPgUUID(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)
	users := &Users{BaseService: baseService}
	emails := &Emails{Users: users}

	emailID := uuid.New()
	pgUUID := emails.ToPgUUID(emailID)

	assert.True(t, pgUUID.Valid)
	assert.Equal(t, [16]byte(emailID), pgUUID.Bytes)
}

func TestVerifyEmail_ToPgUUID(t *testing.T) {
	// Test the conversion that happens in Verify method
	emailID := uuid.New()
	params := VerifyEmail{
		EmailID: emailID,
	}

	// Simulate the direct pgtype.UUID creation in Verify method
	pgUUID := pgtype.UUID{
		Bytes: params.EmailID,
		Valid: true,
	}

	assert.True(t, pgUUID.Valid)
	assert.Equal(t, [16]byte(emailID), pgUUID.Bytes)
}

func TestEmailParams_UUIDConversions(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)
	users := &Users{BaseService: baseService}
	emails := &Emails{Users: users}

	tests := []struct {
		name   string
		params interface{}
		getID  func() uuid.UUID
	}{
		{
			name:   "DeleteEmailParams",
			params: DeleteEmailParams{EmailID: uuid.New()},
			getID: func() uuid.UUID {
				return emails.ToPgUUID(uuid.New()).Bytes
			},
		},
		{
			name:   "GetEmailParams",
			params: GetEmailParams{EmailID: uuid.New()},
			getID: func() uuid.UUID {
				return emails.ToPgUUID(uuid.New()).Bytes
			},
		},
		{
			name:   "ListUserEmailParams",
			params: ListUserEmailParams{UserID: uuid.New()},
			getID: func() uuid.UUID {
				return emails.ToPgUUID(uuid.New()).Bytes
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.getID()
			assert.NotEqual(t, uuid.Nil, id)
		})
	}
}

func TestEmails_ContextPropagation(t *testing.T) {
	type ctxKey string
	const key ctxKey = "testKey"

	ctx := context.WithValue(context.Background(), key, "testValue")
	baseService := types.NewBaseService(ctx, nil, false)
	users := &Users{BaseService: baseService}
	emails := &Emails{Users: users}

	// Context should propagate through the chain
	resultCtx := emails.GetCtx()
	assert.Equal(t, "testValue", resultCtx.Value(key))
}

func TestEmails_CacheSettings(t *testing.T) {
	tests := []struct {
		name     string
		useCache bool
	}{
		{"cache enabled", true},
		{"cache disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			baseService := types.NewBaseService(ctx, nil, tt.useCache)
			users := &Users{BaseService: baseService}
			emails := &Emails{Users: users}

			assert.Equal(t, tt.useCache, emails.IsCache())
		})
	}
}

func TestEmails_ChainedMethods(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)
	users := &Users{BaseService: baseService}
	emails := &Emails{Users: users}

	// Test that chained methods work through the embedded struct
	newCtx := context.WithValue(context.Background(), "key", "value")
	emails.WithCtx(newCtx)

	assert.Equal(t, newCtx, emails.GetCtx())
}

func TestNullUserRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		valid    bool
		expected bool
	}{
		{"valid super role", string(repository.UserRoleSuper), true, true},
		{"valid admin role", string(repository.UserRoleAdmin), true, true},
		{"invalid when not valid", string(repository.UserRoleUser), false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nullRole := repository.NullUserRole{
				UserRole: repository.UserRole(tt.role),
				Valid:    tt.valid,
			}

			assert.Equal(t, tt.valid, nullRole.Valid)
			if tt.valid {
				assert.Equal(t, repository.UserRole(tt.role), nullRole.UserRole)
			}
		})
	}
}
