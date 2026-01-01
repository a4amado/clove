package appservice

import (
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (as *RegionsCtx) List() (*[]repository.Region, error) {
	queries := as.App.Queries
	if as.App.Tx != nil {
		queries = queries.WithTx(*as.App.Tx)
	}

	if as.App.Cache {
		// log in the cache
	}
	region, err := queries.ListAppRegions(as.App.ReqCtx, pgtype.UUID{
		Bytes: as.AppId,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	return &region, nil

}
