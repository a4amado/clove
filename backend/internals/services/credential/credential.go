package credentialservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CredentialService struct {
	*types.BaseService
}

// New creates a CredentialService backed by the given DB and optional cache client.
func New(db *repository.Queries, cacheClient *cache.Client) *CredentialService {
	base := &types.BaseService{DB: db}
	if cacheClient != nil {
		base.WithCache(cacheClient)
	}
	return &CredentialService{BaseService: base}
}

type InsertParams struct {
	Token       string
	UserID      uuid.UUID
	AppID       uuid.UUID
	ChannelID   string
	Type        repository.CredentialType
	Permissions string
	ExpiresAt   time.Time
}

func nullUUID(id uuid.UUID) pgtype.UUID {
	if id == (uuid.UUID{}) {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: id, Valid: true}
}

func (s *CredentialService) Create(args InsertParams) (*repository.Credential, error) {
	cred, err := s.DB.Credential_Insert(s.GetCtx(), repository.Credential_InsertParams{
		Token:  args.Token,
		UserID: nullUUID(args.UserID),
		AppID:  nullUUID(args.AppID),
		ChannelID: pgtype.Text{
			String: args.ChannelID,
			Valid:  args.ChannelID != "",
		},
		Type:        args.Type,
		Permissions: args.Permissions,
		ExpiresAt:   pgtype.Timestamp{Time: args.ExpiresAt, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert credential: %w", err)
	}
	return &cred, nil
}

func (s *CredentialService) GetByToken(ctx context.Context, token string) (*repository.Credential, error) {
	cacheKey := cache.FormatCredentialCacheKey(token)

	if s.CacheReady() {
		var cred repository.Credential
		if err := cache.Get(ctx, s.Cache, cacheKey, &cred); err == nil {
			return &cred, nil
		}
	}

	cred, err := s.DB.Credential_SelectByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("credential not found: %w", err)
	}

	if s.CacheReady() {
		s.CacheAsync(func(aCtx context.Context) error {
			return cache.Set(aCtx, s.Cache, cacheKey, cred, cache.TTLCredential)
		})
	}

	return &cred, nil
}

func (s *CredentialService) Delete(id uuid.UUID) error {
	err := s.DB.Credential_Delete(s.GetCtx(), s.ToPgUUID(id))
	if err != nil {
		return fmt.Errorf("failed to delete credential: %w", err)
	}
	return nil
}

func (s *CredentialService) DeleteExpired() error {
	err := s.DB.Credential_DeleteExpired(s.GetCtx())
	if err != nil {
		return fmt.Errorf("failed to delete expired credentials: %w", err)
	}
	return nil
}
