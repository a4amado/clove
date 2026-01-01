package appservice

import (
	appservice "clove/internals/services/app-service/apps/keys"

	"github.com/google/uuid"
)

func (as *AppCtx) Keys() *appservice.KeysCtx {
	return &appservice.KeysCtx{
		BaseCtx: &as.BaseCtx,
		AppID:   as.AppID,
	}
}
func (as *AppCtx) Key(id uuid.UUID) *appservice.KeyCtx {
	return &appservice.KeyCtx{
		BaseCtx: &as.BaseCtx,
		AppId:   as.AppID,
		KeyId:   id,
	}
}
