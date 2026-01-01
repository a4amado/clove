package appservice

import (
	appservice "clove/internals/services/app-service/apps/regions"
	repository "clove/internals/services/generatedRepo"
)

func (as *AppServiceCtx) Region(region repository.Region) *appservice.RegionCtx {
	return &appservice.RegionCtx{
		App:    as.App,
		AppId:  as.AppId,
		Region: region,
	}
}

func (as *AppServiceCtx) Regions() *appservice.RegionsCtx {
	return &appservice.RegionsCtx{
		App:   as.App,
		AppId: as.AppId,
	}
}
