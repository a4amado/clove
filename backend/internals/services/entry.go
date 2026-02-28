package services

import (
	"clove/internals/cache"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/data/valkeyPool"
	appservice "clove/internals/services/apps"
	credentialservice "clove/internals/services/credential"
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
	Credentials *credentialservice.CredentialService
}

func New(ctx context.Context) *router {
	q := repository.New(postgresPool.Client())

	router := &router{
		baseService: &types.BaseService{
			DB: q,
		},
	}
	router.baseService.WithCtx(ctx)
	router.Apps = &appservice.AppsService{
		BaseService: router.baseService,
	}
	router.Users = &userservice.Users{
		BaseService: router.baseService,
	}
	router.Credentials = &credentialservice.CredentialService{
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
