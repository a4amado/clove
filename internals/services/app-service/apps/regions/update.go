package appservice

import (
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (as *RegionsCtx) Update(region []repository.Region) error {
	queries := as.BaseCtx.Queries
	if as.BaseCtx.Tx != nil {
		queries = queries.WithTx(*as.BaseCtx.Tx)
	}

	if as.BaseCtx.Cache {
		// log in the cache
	}

	return queries.UpdateAppRegions(as.BaseCtx.ReqCtx, repository.UpdateAppRegionsParams{
		Regions: region,
		AppID: pgtype.UUID{
			Bytes: as.AppId,
			Valid: true,
		},
	})

}
