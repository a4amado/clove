package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"

	"github.com/google/uuid"
)

func NewAppsService(args types.ServiceParams, queries *repository.Queries) *AppsService {
	return &AppsService{
		BaseService: types.NewBaseService(args.Ctx, queries, args.UseCache),
	}
}
func NewAppService(args types.ServiceParams, queries *repository.Queries) *AppService {
	return &AppService{
		BaseService: types.NewBaseService(args.Ctx, queries, args.UseCache),
	}
}

type AppsService struct {
	*types.BaseService
}
type AppService struct {
	*types.BaseService
	appID uuid.UUID
}

func (as *AppService) AppID() uuid.UUID {
	return as.appID
}
func (as *AppService) SetAppID(id uuid.UUID) {
	as.appID = id
}

type RegionsService struct {
	*types.BaseService

	appID uuid.UUID
}
type KeysService struct {
	*types.BaseService
	appID uuid.UUID
}

type KeyService struct {
	*types.BaseService
	appID uuid.UUID
	keyID uuid.UUID
}
