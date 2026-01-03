package appservice

import (
	"clove/internals/services/types"

	"github.com/google/uuid"
)

type KeysCtx struct {
	BaseCtx *types.BaseLineServiceCtx
	AppID   uuid.UUID
}
type KeyCtx struct {
	BaseCtx *types.BaseLineServiceCtx
	AppId   uuid.UUID
	KeyId   uuid.UUID
}
