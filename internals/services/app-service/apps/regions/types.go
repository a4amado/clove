package appservice

import (
	appservicetypes "clove/internals/services/app-service/types"
	repository "clove/internals/services/generatedRepo"

	"github.com/google/uuid"
)

type RegionsCtx struct {
	App   *appservicetypes.BaseLineAppCtx
	AppId uuid.UUID
}
type RegionCtx struct {
	App    *appservicetypes.BaseLineAppCtx
	AppId  uuid.UUID
	Region repository.Region
}
