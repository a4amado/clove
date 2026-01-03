package appservice

import (
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (as *RegionsCtx) List() (*[]repository.Region, error) {
	queries := as.BaseCtx.Queries
	if as.BaseCtx.Tx != nil {
		queries = queries.WithTx(*as.BaseCtx.Tx)
	}

	if as.BaseCtx.Cache {
		// log in the cache
	}
	region, err := queries.App_Region_Select(as.BaseCtx.ReqCtx, pgtype.UUID{
		Bytes: as.AppId,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	return &region, nil

}
