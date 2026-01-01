package appservice

import (
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (as *RegionsCtx) Update() error {
	queries := as.App.Queries
	if as.App.Tx != nil {
		queries = queries.WithTx(*as.App.Tx)
	}

	if as.App.Cache {
		// log in the cache
	}

	return queries.UpdateAppRegions(as.App.ReqCtx, repository.UpdateAppRegionsParams{
		Regions: *as.Regions,
		AppID: pgtype.UUID{
			Bytes: as.AppId,
			Valid: true,
		},
	})

}
