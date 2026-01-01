package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"errors"
	"slices"

	"github.com/jackc/pgx/v5/pgtype"
)

func (as *RegionCtx) Remove() (*[]repository.Region, error) {
	queries := as.App.Queries
	if as.App.Tx != nil {
		queries = queries.WithTx(*as.App.Tx)
	}

	if as.App.Cache {
		// log in the cache
	}
	region, err := queries.RemoveRegionFromApp(as.App.ReqCtx, repository.RemoveRegionFromAppParams{
		Region: as.Region,
		ID: pgtype.UUID{
			Bytes: as.AppId,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}
	if slices.Contains(region, as.Region) {
		return nil, errors.New("failed to delete region")
	}
	return &region, nil

}
