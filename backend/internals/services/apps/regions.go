package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"fmt"

	"github.com/google/uuid"
)

type GetAppRegions struct {
	AppID uuid.UUID
}

func (s *RegionsService) List(args GetAppRegions) ([]repository.Region, error) {
	// Try cache first if you implement it
	// if s.useCache {
	//     if regions, err := cache.Apps().Regions().List(s.GetCtx(), s.appID); err == nil {
	//         return regions, nil
	//     }
	// }

	regions, err := s.DB.App_Region_Select(s.GetCtx(), s.ToPgUUID(args.AppID))
	if err != nil {
		return nil, fmt.Errorf("failed to list regions: %w", err)
	}

	// Update cache async if implemented
	// s.CacheAsync(func(ctx context.Context) error {
	//     return cache.Apps().Regions().Set(ctx, s.appID, regions)
	// })

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

	// Invalidate cache async if implemented
	// s.CacheAsync(func(ctx context.Context) error {
	//     return cache.Apps().Regions().Delete(ctx, s.appID)
	// })

	return nil
}
