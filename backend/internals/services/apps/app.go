package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *AppsService) Create(arg repository.App_InsertParams) (*repository.App, error) {
	app, err := s.DB.App_Insert(s.GetCtx(), arg)
	if err != nil {
		return nil, fmt.Errorf("failed to insert app: %w", err)
	}
	return &app, nil
}

type GetParams struct {
	AppID uuid.UUID
}

func (s *AppsService) Get(args GetParams) (*repository.App, error) {
	cacheKey := cache.FormatAppCacheKey(args.AppID)

	if s.CacheReady() {
		var app repository.App
		if err := cache.Get(s.GetCtx(), s.Cache, cacheKey, &app); err == nil {
			return &app, nil
		}
	}

	app, err := s.DB.App_Select(s.GetCtx(), s.ToPgUUID(args.AppID))
	if err != nil {
		return nil, fmt.Errorf("failed to get app: %w", err)
	}

	if s.CacheReady() {
		s.CacheAsync(func(ctx context.Context) error {
			return cache.Set(ctx, s.Cache, cacheKey, app, cache.TTLApp)
		})
	}

	return &app, nil
}

type ListParams struct {
	UserID uuid.UUID
}

func (s *AppsService) List(args ListParams) ([]repository.App, error) {
	apps, err := s.DB.App_List_ByUserID(s.GetCtx(), s.ToPgUUID(args.UserID))
	if err != nil {
		return nil, fmt.Errorf("failed to list apps: %w", err)
	}
	if len(apps) == 0 {
		return []repository.App{}, nil
	}
	return apps, nil
}

type DeleteParams struct {
	AppID uuid.UUID
}

func (s *AppsService) Delete(args DeleteParams) error {
	err := s.DB.App_Delete(s.GetCtx(), s.ToPgUUID(args.AppID))
	if err != nil {
		return fmt.Errorf("failed to delete app: %w", err)
	}

	if s.CacheReady() {
		s.CacheAsync(func(ctx context.Context) error {
			return s.Cache.Delete(ctx, string(cache.FormatAppCacheKey(args.AppID)))
		})
	}

	return nil
}

type UpdateParams struct {
	AppID          uuid.UUID
	AppSlug        string
	AllowedOrigins []string
	AppType        repository.AppType
}

func (s *AppsService) Update(args UpdateParams) (*repository.App, error) {
	app, err := s.DB.App_Update(s.GetCtx(), repository.App_UpdateParams{
		ID:             s.ToPgUUID(args.AppID),
		AppSlug:        args.AppSlug,
		AllowedOrigins: args.AllowedOrigins,
		AppType:        args.AppType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update app: %w", err)
	}

	if s.CacheReady() {
		cacheKey := cache.FormatAppCacheKey(args.AppID)
		s.CacheAsync(func(ctx context.Context) error {
			return cache.Set(ctx, s.Cache, cacheKey, app, cache.TTLApp)
		})
	}

	return &app, nil
}
