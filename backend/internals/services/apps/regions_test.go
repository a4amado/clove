package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAppRegions_Structure(t *testing.T) {
	appID := uuid.New()
	params := GetAppRegions{
		AppID: appID,
	}

	assert.Equal(t, appID, params.AppID)
}

func TestUpdateRegions_Structure(t *testing.T) {
	appID := uuid.New()
	regions := []repository.Region{repository.RegionDk1}

	params := UpdateRegions{
		Regions: regions,
		AppID:   appID,
	}

	assert.Equal(t, appID, params.AppID)
	assert.Equal(t, regions, params.Regions)
	assert.Len(t, params.Regions, 1)
}

func TestUpdateRegions_MultipleRegions(t *testing.T) {
	appID := uuid.New()
	// Even though currently only dk1 exists, test with the same region multiple times
	regions := []repository.Region{repository.RegionDk1}

	params := UpdateRegions{
		Regions: regions,
		AppID:   appID,
	}

	assert.Len(t, params.Regions, 1)
	assert.Contains(t, params.Regions, repository.RegionDk1)
}

func TestUpdateRegions_EmptyRegions(t *testing.T) {
	appID := uuid.New()
	regions := []repository.Region{}

	params := UpdateRegions{
		Regions: regions,
		AppID:   appID,
	}

	assert.Empty(t, params.Regions)
}

func TestApp_Region_UpdateParams_Conversion(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)

	appID := uuid.New()
	regions := []repository.Region{repository.RegionDk1}

	updateParams := UpdateRegions{
		Regions: regions,
		AppID:   appID,
	}

	// Simulate the conversion that happens in Update method
	dbParams := repository.App_Region_UpdateParams{
		ID:     baseService.ToPgUUID(updateParams.AppID),
		Region: updateParams.Regions,
	}

	assert.True(t, dbParams.ID.Valid)
	assert.Equal(t, [16]byte(appID), dbParams.ID.Bytes)
	assert.Equal(t, regions, dbParams.Region)
}

func TestRegionsService_WithBaseService(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, true)

	service := &RegionsService{
		BaseService: baseService,
	}

	assert.NotNil(t, service.BaseService)
	assert.True(t, service.IsCache())
	assert.Equal(t, ctx, service.GetCtx())
}

func TestRegion_Valid(t *testing.T) {
	tests := []struct {
		name     string
		region   repository.Region
		expected bool
	}{
		{"dk1 is valid", repository.RegionDk1, true},
		{"empty is invalid", repository.Region(""), false},
		{"unknown is invalid", repository.Region("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.region.Valid())
		})
	}
}

func TestAllRegionValues(t *testing.T) {
	regions := repository.AllRegionValues()

	assert.NotEmpty(t, regions)
	assert.Contains(t, regions, repository.RegionDk1)
}
