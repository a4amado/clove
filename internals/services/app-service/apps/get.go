package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (as AppCtx) Get() (*repository.App, error) {
	queries := as.BaseCtx.Queries
	if as.BaseCtx.Tx != nil {
		queries = queries.WithTx(*as.BaseCtx.Tx)
	}

	if as.BaseCtx.Cache {
		app, err := cache.Apps().Get(as.BaseCtx.ReqCtx, as.AppID)
		if err == nil {
			return app, nil
		}
	}
	app, err := queries.App_Select(as.BaseCtx.ReqCtx, pgtype.UUID{
		Bytes: as.AppID,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}
	err = cache.Apps().Set(as.BaseCtx.ReqCtx, app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}
