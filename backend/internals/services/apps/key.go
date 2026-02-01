package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ListKeysParams struct {
	Page  int32
	AppId uuid.UUID
}

func (s *KeysService) List(args ListKeysParams) ([]repository.AppApiKey, error) {
	if args.Page <= 0 {
		args.Page = 0
	} else {
		args.Page = args.Page - 1
	}

	keys, err := s.DB.App_Key_List(s.GetCtx(), repository.App_Key_ListParams{
		AppID:   s.ToPgUUID(args.AppId),
		PageIdx: args.Page,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	if len(keys) == 0 {
		return []repository.AppApiKey{}, nil
	}

	return keys, nil
}

type CreateKeyParams struct {
	AppID     uuid.UUID
	KeyString string
	KeyName   string
}

func (s *KeysService) Create(args CreateKeyParams) (*repository.AppApiKey, error) {
	key, err := s.DB.App_Key_Insert(s.GetCtx(), repository.App_Key_InsertParams{
		AppID: s.ToPgUUID(args.AppID),
		Key:   s.ToPgText(args.KeyString),
		Name:  s.ToPgText(args.KeyName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create key: %w", err)
	}

	// Cache the key async
	s.CacheAsync(func(ctx context.Context) error {
		return cache.Apps().Keys().Set(s.GetCtx(), args.AppID, key.ID.Bytes, key.Key.String)
	})

	return &key, nil
}

type GetKeyParams struct {
	KeyID uuid.UUID
	AppID uuid.UUID
}

func (s *KeysService) Get(args GetKeyParams) (*repository.AppApiKey, error) {
	// Try cache first

	// Fetch from DB
	key, err := s.DB.App_Key_Select(s.GetCtx(), repository.App_Key_SelectParams{
		KeyID: s.ToPgUUID(args.KeyID),
		AppID: s.ToPgUUID(args.KeyID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	// Cache async
	s.CacheAsync(func(ctx context.Context) error {
		return cache.Apps().Keys().Set(s.GetCtx(), args.AppID, args.KeyID, key.Key.String)
	})

	return &key, nil
}

type DeleteKeyParams struct {
	AppID uuid.UUID
	KeyID uuid.UUID
}

func (s *KeysService) Delete(args DeleteKeyParams) error {
	rowsAffected, err := s.DB.App_Key_Delete(s.GetCtx(), repository.App_Key_DeleteParams{
		ID:    s.ToPgUUID(args.KeyID),
		AppID: s.ToPgUUID(args.AppID),
	})
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("key not found")
	}

	// Invalidate cache async
	s.CacheAsync(func(ctx context.Context) error {
		return cache.Apps().Keys().Delete(s.GetCtx(), args.AppID, args.KeyID)
	})

	return nil
}
