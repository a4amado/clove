package appservice

import (
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (as *KeysCtx) List(page int32) ([]repository.AppApiKey, error) {
	queries := as.App.Queries
	if as.App.Tx != nil {
		queries = queries.WithTx(*as.App.Tx)
	}
	if page <= 0 {
		page = 1
	}
	return queries.ListAppApiKeys(as.App.ReqCtx, repository.ListAppApiKeysParams{
		AppID: pgtype.UUID{
			Bytes: as.AppId,
			Valid: true,
		},
		PageIdx: page,
	})
}
