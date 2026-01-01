package appservice

import (
	"clove/internals/services/types"

	"github.com/google/uuid"
)

type AppsCtx struct {
	BaseCtx types.BaseLineServiceCtx
}
type AppCtx struct {
	BaseCtx types.BaseLineServiceCtx
	AppID   uuid.UUID
}
