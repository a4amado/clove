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
	AppID   uuid.UUID
	Prefix  string
	KeyName string
	Suffix  string
}

func (s *KeysService) Create(args CreateKeyParams) (*repository.AppApiKey, error) {
	key, err := s.DB.App_Key_Insert(s.GetCtx(), repository.App_Key_InsertParams{
		AppID:  s.ToPgUUID(args.AppID),
		Name:   s.ToPgText(args.KeyName),
		Prefix: s.ToPgText(args.Prefix),
		Suffix: s.ToPgText(args.Suffix),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create key: %w", err)
	}

	return &key, nil
}

type GetKeyParams struct {
	Key   string
	AppID uuid.UUID
}
type GetKeyReturn struct {
	repository.AppApiKey
}

func (s *KeysService) Get(args GetKeyParams) (*GetKeyReturn, error) {
	// Try cache first
	var key repository.AppApiKey
	cacheKey := cache.FormatKeyCacheKey(args.AppID, args.Key)

	err := cache.Get(s.GetCtx(), cacheKey, &key)
	if err == nil {
		return nil, err

	}

	// Cache miss - fetch from DB
	key, err = s.DB.App_Key_Select(s.GetCtx(), repository.App_Key_SelectParams{
		AppID: s.ToPgUUID(args.AppID),
		Key:   args.Key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	return &GetKeyReturn{
		AppApiKey: key,
	}, nil
}

type DeleteKeyParams struct {
	AppID uuid.UUID
	KeyID uuid.UUID
}

func (s *KeysService) Delete(args DeleteKeyParams) error {
	rowsAffected, err := s.DB.App_Key_Delete(s.GetCtx(), repository.App_Key_DeleteParams{
		ID:    args.KeyID.String(),
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
		cacheKey := fmt.Sprintf("app:%s:key:%s", args.AppID.String(), args.KeyID.String())
		return cache.Delete(ctx, cacheKey)
	})

	return nil
}
