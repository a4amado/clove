package appservice

import (
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (as *RegionCtx) List() (*[]repository.AppApiKey, error) {
	queries := as.App.Queries
	if as.App.Tx != nil {
		queries = queries.WithTx(*as.App.Tx)
	}

	if as.App.Cache {
		// log in the cache
	}
	region, err := queries.ListAppApiKeys(as.App.ReqCtx, repository.ListAppApiKeysParams{
		AppID: pgtype.UUID{
			Bytes: as.AppId,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &region, nil

}
