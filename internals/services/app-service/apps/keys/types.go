package appservice

import (
	appservicetypes "clove/internals/services/app-service/types"

	"github.com/google/uuid"
)

type KeysCtx struct {
	App   *appservicetypes.BaseLineAppCtx
	AppId uuid.UUID
}
type KeyCtx struct {
	App   *appservicetypes.BaseLineAppCtx
	AppId uuid.UUID
	KeyId uuid.UUID
}
