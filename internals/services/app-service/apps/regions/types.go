package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"

	"github.com/google/uuid"
)

type RegionsCtx struct {
	BaseCtx *types.BaseLineServiceCtx
	AppId   uuid.UUID
	Regions *[]repository.Region
}
