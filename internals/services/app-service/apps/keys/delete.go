package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (a *KeyCtx) Delete() (int64, error) {
	q := a.App.Queries
	if a.App.Tx != nil {
		q = q.WithTx(*a.App.Tx)
	}
	go cache.Apps().Keys().Delete(a.App.ReqCtx, a.AppId, a.KeyId)

	return q.DeleteApiKey(a.App.ReqCtx, repository.DeleteApiKeyParams{
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
