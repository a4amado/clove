package services

import (
	postgresPool "clove/internals/data/postgres/pool"
	appservice "clove/internals/services/apps"
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
	userservice "clove/internals/services/user"
	"context"

	"github.com/jackc/pgx/v5"
)

type router struct {
	baseService *types.BaseService
	Apps        *appservice.AppsService
	Users       *userservice.Users
}

func New(ctx context.Context) *router {

	q := repository.New(postgresPool.Client())

	router := &router{
		baseService: &types.BaseService{
			DB: q,
		},
	}
	router.Apps = &appservice.AppsService{
		BaseService: router.baseService,
	}
	router.Users = &userservice.Users{
		BaseService: router.baseService,
	}

	return router
}

func (b *router) WithCache() *router {
	b.baseService.WithCache()
	return b

}
func (b *router) WithTx() (*router, pgx.Tx, error) {
	tx, err := postgresPool.NewTx(b.baseService.GetCtx(), pgx.TxOptions{})
	if err != nil {
		return nil, nil, err
	}
	b.baseService.DB = b.baseService.DB.WithTx(tx)

	return b, tx, nil
}
