package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
)

func NewAppsService(args types.ServiceParams, queries *repository.Queries) *AppsService {
	return &AppsService{
		BaseService: types.NewBaseService(args.Ctx, queries, args.UseCache),
	}
}
func NewAppService(args types.ServiceParams, queries *repository.Queries) *AppsService {
	return &AppsService{
		BaseService: types.NewBaseService(args.Ctx, queries, args.UseCache),
	}
}

type AppsService struct {
	*types.BaseService
	Regions *RegionsService
	Keys    *KeysService
}

type RegionsService struct {
	*types.BaseService
}
type KeysService struct {
	*types.BaseService
}
