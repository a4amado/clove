package appservice

import (
	appservicetypes "clove/internals/services/app-service/types"

	"github.com/google/uuid"
)

type AppsServiceCtx struct {
	App *appservicetypes.BaseLineAppCtx
}
type AppServiceCtx struct {
	App   *appservicetypes.BaseLineAppCtx
	AppId uuid.UUID
}
