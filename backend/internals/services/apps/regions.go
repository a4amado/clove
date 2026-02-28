package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type GetAppRegions struct {
	AppID uuid.UUID
}

func (s *RegionsService) List(args GetAppRegions) ([]repository.Region, error) {
	cacheKey := cache.FormatAppRegionsCacheKey(args.AppID)

	if s.CacheReady() {
		var regions []repository.Region
		if err := cache.Get(s.GetCtx(), s.Cache, cacheKey, &regions); err == nil {
			return regions, nil
		}
	}

	regions, err := s.DB.App_Region_Select(s.GetCtx(), s.ToPgUUID(args.AppID))
	if err != nil {
		return nil, fmt.Errorf("failed to list regions: %w", err)
	}

	if s.CacheReady() {
		s.CacheAsync(func(ctx context.Context) error {
			return cache.Set(ctx, s.Cache, cacheKey, regions, cache.TTLRegions)
		})
	}

	return regions, nil
}

type UpdateRegions struct {
	Regions []repository.Region
	AppID   uuid.UUID
}

func (s *RegionsService) Update(args UpdateRegions) error {
	err := s.DB.App_Region_Update(s.GetCtx(), repository.App_Region_UpdateParams{
		ID:     s.ToPgUUID(args.AppID),
		Region: args.Regions,
	})
	if err != nil {
		return fmt.Errorf("failed to update regions: %w", err)
	}

	if s.CacheReady() {
		s.CacheAsync(func(ctx context.Context) error {
			return s.Cache.Delete(ctx, string(cache.FormatAppRegionsCacheKey(args.AppID)))
		})
	}

	return nil
}
