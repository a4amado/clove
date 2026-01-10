package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"context"
	"fmt"
)

func (s *AppService) Get() (*repository.App, error) {
	// Try cache first
	if s.Cache() {
		if app, err := cache.Apps().Get(s.CTX(), s.appID); err == nil {
			return app, nil
		}
	}

	// Fetch from DB
	app, err := s.Q().App_Select(s.CTX(), s.ToPgUUID(s.appID))
	if err != nil {
		return nil, fmt.Errorf("failed to get app: %w", err)
	}

	// Update cache async
	s.CacheAsync(func(ctx context.Context) error {
		return cache.Apps().Set(ctx, app)
	})

	return &app, nil
}

func (s *AppService) Delete() error {
	if err := s.Q().App_Delete(s.CTX(), s.ToPgUUID(s.appID)); err != nil {
		return fmt.Errorf("failed to delete app: %w", err)
	}

	// Invalidate cache async
	s.CacheAsync(func(ctx context.Context) error {
		// TODO:
		return nil
	})

	return nil
}

// Regions returns a service scoped to this app's regions
func (s *AppService) Regions() *RegionsService {
	return &RegionsService{
		BaseService: s.BaseService,
		appID:       s.appID,
	}
}
func (s *RegionsService) List() ([]repository.Region, error) {
	// Try cache first if you implement it
	// if s.useCache {
	//     if regions, err := cache.Apps().Regions().List(s.CTX(), s.appID); err == nil {
	//         return regions, nil
	//     }
	// }

	regions, err := s.Q().App_Region_Select(s.CTX(), s.ToPgUUID(s.appID))
	if err != nil {
		return nil, fmt.Errorf("failed to list regions: %w", err)
	}

	// Update cache async if implemented
	// s.CacheAsync(func(ctx context.Context) error {
	//     return cache.Apps().Regions().Set(ctx, s.appID, regions)
	// })

	return regions, nil
}

func (s *RegionsService) Update(regions []repository.Region) error {
	err := s.Q().App_Region_Update(s.CTX(), repository.App_Region_UpdateParams{
		ID:     s.ToPgUUID(s.appID),
		Region: regions,
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
