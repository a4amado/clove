package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"context"
	"fmt"
	"time"

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
	cacheKey := cache.FormatKeyCacheKey(args.AppID, args.Key)

	if s.CacheReady() {
		var key repository.AppApiKey
		if err := cache.Get(s.GetCtx(), s.Cache, cacheKey, &key); err == nil {
			return &GetKeyReturn{AppApiKey: key}, nil
		}
	}

	key, err := s.DB.App_Key_Select(s.GetCtx(), repository.App_Key_SelectParams{
		AppID: s.ToPgUUID(args.AppID),
		Key:   args.Key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	if s.CacheReady() {
		s.CacheAsync(func(ctx context.Context) error {
			return cache.Set(ctx, s.Cache, cacheKey, key, 24*time.Hour)
		})
	}

	return &GetKeyReturn{AppApiKey: key}, nil
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

	if s.CacheReady() {
		s.CacheAsync(func(ctx context.Context) error {
			cacheKey := cache.FormatKeyCacheKey(args.AppID, args.KeyID.String())
			return s.Cache.Delete(ctx, string(cacheKey))
		})
	}

	return nil
}
