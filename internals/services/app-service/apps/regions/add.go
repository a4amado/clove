package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"errors"
	"slices"

	"github.com/jackc/pgx/v5/pgtype"
)

func (as *RegionCtx) Add() (*[]repository.Region, error) {
	queries := as.App.Queries
	if as.App.Tx != nil {
		queries = queries.WithTx(*as.App.Tx)
	}

	if as.App.Cache {
		// log in the cache
	}
	region, err := queries.AddRegionToApp(as.App.ReqCtx, repository.AddRegionToAppParams{
		Region: as.Region,
		ID: pgtype.UUID{
			Bytes: as.AppId,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}
	if !slices.Contains(region, as.Region) {
		return nil, errors.New("failed to add region")
	}
	return &region, nil

}
