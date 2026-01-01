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
		page = 0
	} else {
		page = page - 1
	}
	keys, err := queries.ListAppApiKeys(as.App.ReqCtx, repository.ListAppApiKeysParams{
		AppID: pgtype.UUID{
			Bytes: as.AppId,
			Valid: true,
		},
		PageIdx: page,
	})
	if err != nil {
		return nil, err
	}
	if keys == nil || len(keys) == 0 {
		return []repository.AppApiKey{}, nil
	}
	return keys, nil
}
