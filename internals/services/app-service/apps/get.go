package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (as AppServiceCtx) Get() (*repository.App, error) {
	queries := as.App.Queries
	if as.App.Tx != nil {
		queries = queries.WithTx(*as.App.Tx)
	}

	if as.App.Cache {
		app, err := cache.Apps().Get(as.App.ReqCtx, as.AppId)
		if err == nil {
			return app, nil
		}
	}
	app, err := queries.FindAppById(as.App.ReqCtx, pgtype.UUID{
		Bytes: as.AppId,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}
	err = cache.Apps().Set(as.App.ReqCtx, app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}
