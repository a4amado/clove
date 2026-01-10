package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Keys returns a service scoped to this app's keys (collection)
func (s *AppService) Keys() *KeysService {
	return &KeysService{
		BaseService: s.BaseService,
		appID:       s.appID,
	}
}

// Key returns a service scoped to a specific key
func (s *AppService) Key(keyID uuid.UUID) *KeyService {
	return &KeyService{
		BaseService: s.BaseService,
		appID:       s.appID,
		keyID:       keyID,
	}
}

func (s *KeysService) List(page int32) ([]repository.AppApiKey, error) {
	if page <= 0 {
		page = 0
	} else {
		page = page - 1
	}

	keys, err := s.Q().App_Key_List(s.CTX(), repository.App_Key_ListParams{
		AppID:   s.ToPgUUID(s.appID),
		PageIdx: page,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	if len(keys) == 0 {
		return []repository.AppApiKey{}, nil
	}

	return keys, nil
}

func (s *KeysService) Create(name, keyString string) (*repository.AppApiKey, error) {
	key, err := s.Q().App_Key_Insert(s.CTX(), repository.App_Key_InsertParams{
		AppID: s.ToPgUUID(s.appID),
		Key:   s.ToPgText(keyString),
		Name:  s.ToPgText(name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create key: %w", err)
	}

	// Cache the key async
	s.CacheAsync(func(ctx context.Context) error {
		return cache.Apps().Keys().Set(s.CTX(), s.appID, key.ID.Bytes, key.Key.String)
	})

	return &key, nil
}

func (s *KeyService) Get() (string, error) {
	// Try cache first
	if s.Cache() {
		if keyStr, err := cache.Apps().Keys().Get(s.CTX(), s.appID, s.keyID); err == nil {
			return *keyStr, nil
		}
	}

	// Fetch from DB
	key, err := s.Q().App_Key_Select(s.CTX(), repository.App_Key_SelectParams{
		KeyID: s.ToPgUUID(s.keyID),
		AppID: s.ToPgUUID(s.appID),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get key: %w", err)
	}

	// Cache async
	s.CacheAsync(func(ctx context.Context) error {
		return cache.Apps().Keys().Set(s.CTX(), s.appID, s.keyID, key.Key.String)
	})

	return key.Key.String, nil
}

func (s *KeyService) Delete() error {
	rowsAffected, err := s.Q().App_Key_Delete(s.CTX(), repository.App_Key_DeleteParams{
		ID:    s.ToPgUUID(s.keyID),
		AppID: s.ToPgUUID(s.appID),
	})
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("key not found")
	}

	// Invalidate cache async
	s.CacheAsync(func(ctx context.Context) error {
		return cache.Apps().Keys().Delete(s.CTX(), s.appID, s.keyID)
	})

	return nil
}
