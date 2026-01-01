package services

import (
	postgresPool "clove/internals/data/postgres/pool"
	appservice "clove/internals/services/app-service/apps"
	"clove/internals/services/app-service/types"
	repository "clove/internals/services/generatedRepo"
	userservice "clove/internals/services/usr-service"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func User() *userservice.UserService {
	newPool := postgresPool.Client()
	return &userservice.UserService{
		Pool:    newPool,
		Queries: repository.New(newPool),
	}
}

// New constructs and returns an AppServiceCtx with an initialized database pool and query helper.
// The returned AppServiceCtx has Pool set to a database connection pool and Queries initialized to use that pool.
func Apps(ctx context.Context, tx *pgx.Tx, cache bool) *appservice.AppsServiceCtx {
	newPool := postgresPool.Client()
	return &appservice.AppsServiceCtx{
		App: &types.BaseLineAppCtx{
			Pool:    newPool,
			Queries: repository.New(newPool),
			Tx:      tx,
			ReqCtx:  ctx,
			Cache:   cache,
		},
	}
}
func App(ctx context.Context, tx *pgx.Tx, cache bool, id uuid.UUID) *appservice.AppServiceCtx {
	newPool := postgresPool.Client()
	return &appservice.AppServiceCtx{
		App: &types.BaseLineAppCtx{
			Pool:    newPool,
			Queries: repository.New(newPool),
			Tx:      tx,
			ReqCtx:  ctx,
			Cache:   cache,
		},
		AppId: id,
	}
}
