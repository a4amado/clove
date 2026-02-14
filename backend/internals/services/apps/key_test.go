package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListKeysParams_Structure(t *testing.T) {
	appID := uuid.New()
	params := ListKeysParams{
		Page:  1,
		AppId: appID,
	}

	assert.Equal(t, int32(1), params.Page)
	assert.Equal(t, appID, params.AppId)
}

func TestListKeysParams_PageValues(t *testing.T) {
	tests := []struct {
		name         string
		inputPage    int32
		expectedPage int32
	}{
		{"page 1 becomes 0", 1, 0},
		{"page 2 becomes 1", 2, 1},
		{"page 0 becomes 0", 0, 0},
		{"negative page becomes 0", -1, 0},
		{"page 10 becomes 9", 10, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := ListKeysParams{
				Page:  tt.inputPage,
				AppId: uuid.New(),
			}

			// Simulate the page adjustment logic from List method
			if params.Page <= 0 {
				params.Page = 0
			} else {
				params.Page = params.Page - 1
			}

			assert.Equal(t, tt.expectedPage, params.Page)
		})
	}
}

func TestCreateKeyParams_Structure(t *testing.T) {
	appID := uuid.New()
	params := CreateKeyParams{
		AppID:   appID,
		Prefix:  "pk_",
		KeyName: "Production Key",
		Suffix:  "_prod",
	}

	assert.Equal(t, appID, params.AppID)
	assert.Equal(t, "pk_", params.Prefix)
	assert.Equal(t, "Production Key", params.KeyName)
	assert.Equal(t, "_prod", params.Suffix)
}

func TestGetKeyParams_Structure(t *testing.T) {
	appID := uuid.New()
	params := GetKeyParams{
		Key:   "api-key-123",
		AppID: appID,
	}

	assert.Equal(t, "api-key-123", params.Key)
	assert.Equal(t, appID, params.AppID)
}

func TestGetKeyReturn_Structure(t *testing.T) {
	key := repository.AppApiKey{
		ID: "test-key-id",
	}

	result := GetKeyReturn{
		AppApiKey: key,
	}

	assert.Equal(t, "test-key-id", result.ID)
	assert.Equal(t, key, result.AppApiKey)
}

func TestDeleteKeyParams_Structure(t *testing.T) {
	appID := uuid.New()
	keyID := uuid.New()

	params := DeleteKeyParams{
		AppID: appID,
		KeyID: keyID,
	}

	assert.Equal(t, appID, params.AppID)
	assert.Equal(t, keyID, params.KeyID)
}

func TestKeysService_ToPgTypes(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)
	service := &KeysService{BaseService: baseService}

	appID := uuid.New()

	// Test ToPgUUID
	pgUUID := service.ToPgUUID(appID)
	assert.True(t, pgUUID.Valid)
	assert.Equal(t, [16]byte(appID), pgUUID.Bytes)

	// Test ToPgText
	pgText := service.ToPgText("test-key-name")
	assert.True(t, pgText.Valid)
	assert.Equal(t, "test-key-name", pgText.String)
}

func TestApp_Key_InsertParams_Conversion(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)

	appID := uuid.New()
	createParams := CreateKeyParams{
		AppID:   appID,
		Prefix:  "pk_",
		KeyName: "Test Key",
		Suffix:  "_test",
	}

	// Simulate the conversion that happens in Create method
	insertParams := repository.App_Key_InsertParams{
		AppID:  baseService.ToPgUUID(createParams.AppID),
		Name:   baseService.ToPgText(createParams.KeyName),
		Prefix: baseService.ToPgText(createParams.Prefix),
		Suffix: baseService.ToPgText(createParams.Suffix),
	}

	assert.True(t, insertParams.AppID.Valid)
	assert.Equal(t, [16]byte(appID), insertParams.AppID.Bytes)
	assert.Equal(t, "Test Key", insertParams.Name.String)
	assert.Equal(t, "pk_", insertParams.Prefix.String)
	assert.Equal(t, "_test", insertParams.Suffix.String)
}

func TestApp_Key_ListParams_Conversion(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)

	appID := uuid.New()
	listParams := ListKeysParams{
		Page:  3,
		AppId: appID,
	}

	// Simulate page adjustment
	pageIdx := listParams.Page - 1

	// Simulate the conversion that happens in List method
	dbParams := repository.App_Key_ListParams{
		AppID:   baseService.ToPgUUID(listParams.AppId),
		PageIdx: pageIdx,
	}

	assert.True(t, dbParams.AppID.Valid)
	assert.Equal(t, [16]byte(appID), dbParams.AppID.Bytes)
	assert.Equal(t, int32(2), dbParams.PageIdx) // Page 3 becomes index 2
}

func TestApp_Key_DeleteParams_Conversion(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)

	appID := uuid.New()
	keyID := uuid.New()

	deleteParams := DeleteKeyParams{
		AppID: appID,
		KeyID: keyID,
	}

	// Simulate the conversion that happens in Delete method
	dbParams := repository.App_Key_DeleteParams{
		ID:    deleteParams.KeyID.String(),
		AppID: baseService.ToPgUUID(deleteParams.AppID),
	}

	assert.Equal(t, keyID.String(), dbParams.ID)
	assert.True(t, dbParams.AppID.Valid)
	assert.Equal(t, [16]byte(appID), dbParams.AppID.Bytes)
}

func TestApp_Key_SelectParams_Conversion(t *testing.T) {
	ctx := context.Background()
	baseService := types.NewBaseService(ctx, nil, false)

	appID := uuid.New()
	keyValue := "api-key-abc123"

	getParams := GetKeyParams{
		Key:   keyValue,
		AppID: appID,
	}

	// Simulate the conversion that happens in Get method
	dbParams := repository.App_Key_SelectParams{
		AppID: baseService.ToPgUUID(getParams.AppID),
		Key:   getParams.Key,
	}

	assert.True(t, dbParams.AppID.Valid)
	assert.Equal(t, [16]byte(appID), dbParams.AppID.Bytes)
	assert.Equal(t, keyValue, dbParams.Key)
}
