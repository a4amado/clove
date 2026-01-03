package appservice

import (
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (as *KeysCtx) List(page int32) ([]repository.AppApiKey, error) {
	queries := as.BaseCtx.Queries
	if as.BaseCtx.Tx != nil {
		queries = queries.WithTx(*as.BaseCtx.Tx)
	}
	if page <= 0 {
		page = 0
	} else {
		page = page - 1
	}
	keys, err := queries.App_Key_List(as.BaseCtx.ReqCtx, repository.App_Key_ListParams{
		AppID: pgtype.UUID{
			Bytes: as.AppID,
			Valid: true,
		},
		PageIdx: page,
	})

	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return []repository.AppApiKey{}, nil
	}
	return keys, nil
}
