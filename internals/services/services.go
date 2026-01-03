package services

import (
	postgresPool "clove/internals/data/postgres/pool"
	appservice "clove/internals/services/app-service/apps"
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
	userservice "clove/internals/services/usr-service"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type baseLineCtxWrapper struct {
	*types.BaseLineServiceCtx
}

func (a *baseLineCtxWrapper) App(id uuid.UUID) *appservice.AppCtx {
	return &appservice.AppCtx{
		BaseCtx: *a.BaseLineServiceCtx,
		AppID:   id,
	}
}
func (a *baseLineCtxWrapper) Apps() *appservice.AppsCtx {
	return &appservice.AppsCtx{
		BaseCtx: *a.BaseLineServiceCtx,
	}
}

func User() *userservice.UserService {
	newPool := postgresPool.Client()
	return &userservice.UserService{
		Pool:    newPool,
		Queries: repository.New(newPool),
	}
}
func C(ctx context.Context, tx *pgx.Tx, cache bool) *baseLineCtxWrapper {
	newPool := postgresPool.Client()
	return &baseLineCtxWrapper{
		BaseLineServiceCtx: &types.BaseLineServiceCtx{
			Pool:    newPool,
			Queries: repository.New(newPool),
			Tx:      tx,
			ReqCtx:  ctx,
			Cache:   cache,
		},
	}
}
