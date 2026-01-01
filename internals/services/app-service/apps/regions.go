package appservice

import (
	appservice "clove/internals/services/app-service/apps/regions"
	repository "clove/internals/services/generatedRepo"
)

func (as *AppServiceCtx) Regions(region *[]repository.Region) *appservice.RegionsCtx {
	return &appservice.RegionsCtx{
		App:     as.App,
		AppId:   as.AppId,
		Regions: region,
	}
}
