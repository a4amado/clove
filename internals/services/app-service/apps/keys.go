package appservice

import (
	appservice "clove/internals/services/app-service/apps/keys"

	"github.com/google/uuid"
)

func (as *AppServiceCtx) Keys() *appservice.KeysCtx {
	return &appservice.KeysCtx{
		App:   as.App,
		AppId: as.AppId,
	}
}
func (as *AppServiceCtx) Key(id uuid.UUID) *appservice.KeyCtx {
	return &appservice.KeyCtx{
		App:   as.App,
		AppId: as.AppId,
		KeyId: id,
	}
}
