package appservice

import (
	appservice "clove/internals/services/app-service/apps/regions"
)

func (as *AppCtx) Regions() *appservice.RegionsCtx {
	return &appservice.RegionsCtx{}
}
