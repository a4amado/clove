package services

import (
	"clove/internals/cache"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/data/valkeyPool"
	authservice "clove/internals/services/auth"
	appservice "clove/internals/services/apps"
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
	userservice "clove/internals/services/user"
	"context"

	"github.com/jackc/pgx/v5"
)

type router struct {
	baseService *types.BaseService
	Auth        *authservice.AuthService
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
	router.baseService.WithCtx(ctx)
	router.Auth = &authservice.AuthService{
		BaseService: router.baseService,
	}
	router.Apps = &appservice.AppsService{
		BaseService: router.baseService,
		Keys:        &appservice.KeysService{BaseService: router.baseService},
		Regions:     &appservice.RegionsService{BaseService: router.baseService},
	}
	router.Users = &userservice.Users{
		BaseService: router.baseService,
	}
	return router
}

func (b *router) WithCache() *router {
	cacheClient := cache.New(valkeyPool.Client(valkeyPool.ValkeyStore))
	b.baseService.WithCache(cacheClient)
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
