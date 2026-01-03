package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (a *KeyCtx) Delete() (int64, error) {
	q := a.BaseCtx.Queries
	if a.BaseCtx.Tx != nil {
		q = q.WithTx(*a.BaseCtx.Tx)
	}
	go cache.Apps().Keys().Delete(a.BaseCtx.ReqCtx, a.AppId, a.KeyId)

	return q.App_Key_Delete(a.BaseCtx.ReqCtx, repository.App_Key_DeleteParams{
		ID: pgtype.UUID{
			Bytes: a.KeyId,
			Valid: true,
		},
		AppID: pgtype.UUID{
			Bytes: a.AppId,
			Valid: true,
		},
	})

}
