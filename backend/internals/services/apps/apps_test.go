package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewAppsService(t *testing.T) {
	ctx := context.Background()
	params := types.ServiceParams{
		Ctx:      ctx,
		UseCache: false,
	}

	service := NewAppsService(params, nil)

	assert.NotNil(t, service)
	assert.NotNil(t, service.BaseService)
	assert.Equal(t, ctx, service.GetCtx())
	assert.False(t, service.IsCache())
}

func TestNewAppsService_WithCache(t *testing.T) {
	ctx := context.Background()
	params := types.ServiceParams{
		Ctx:      ctx,
		UseCache: true,
	}

	service := NewAppsService(params, nil)

	assert.True(t, service.IsCache())
}

func TestNewAppService(t *testing.T) {
	ctx := context.Background()
	params := types.ServiceParams{
		Ctx:      ctx,
		UseCache: false,
	}

	service := NewAppService(params, nil)

	assert.NotNil(t, service)
	assert.NotNil(t, service.BaseService)
}

func TestAppsService_Structure(t *testing.T) {
	service := &AppsService{}

	// Verify struct fields exist
	assert.Nil(t, service.Regions)
	assert.Nil(t, service.Keys)
	assert.Nil(t, service.BaseService)
}

func TestGetParams_Structure(t *testing.T) {
	appID := uuid.New()
	params := GetParams{
		AppID: appID,
	}

	assert.Equal(t, appID, params.AppID)
}

func TestApp_InsertParams_Structure(t *testing.T) {
	userID := uuid.New()
	params := repository.App_InsertParams{
		AppSlug:        "test-app",
		Regions:        []repository.Region{repository.RegionDk1},
		AppType:        repository.AppTypeFree,
		UserID:         types.NewBaseService(context.Background(), nil, false).ToPgUUID(userID),
		AllowedOrigins: []string{"http://localhost:3000"},
	}

	assert.Equal(t, "test-app", params.AppSlug)
	assert.Contains(t, params.Regions, repository.RegionDk1)
	assert.Equal(t, repository.AppTypeFree, params.AppType)
	assert.True(t, params.UserID.Valid)
	assert.Len(t, params.AllowedOrigins, 1)
}

func TestRegionsService_Structure(t *testing.T) {
	service := &RegionsService{}
	assert.Nil(t, service.BaseService)
}

func TestKeysService_Structure(t *testing.T) {
	service := &KeysService{}
	assert.Nil(t, service.BaseService)
}
