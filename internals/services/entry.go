package services

import (
	postgresPool "clove/internals/data/postgres/pool"
	appservice "clove/internals/services/apps"
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
	userservice "clove/internals/services/user"

	"github.com/google/uuid"
)

type router struct {
	baseService *types.BaseService
}

func New(args types.ServiceParams) *router {

	q := repository.New(postgresPool.Client())
	if args.Tx != nil {
		q = q.WithTx(*args.Tx)
	}
	return &router{
		baseService: types.NewBaseService(args.Ctx, q, args.UseCache),
	}
}
func (r *router) App(appId uuid.UUID) *appservice.AppService {
	appsrvs := &appservice.AppService{
		BaseService: r.baseService,
	}
	appsrvs.SetAppID(appId)
	return appsrvs
}
func (r *router) Apps() *appservice.AppsService {
	appsrvs := &appservice.AppsService{
		BaseService: r.baseService,
	}
	return appsrvs
}
func (r *router) User(userId uuid.UUID) *userservice.User {
	usrsrvs := &userservice.User{
		BaseService: r.baseService,
	}
	usrsrvs.SetUserId(userId)
	return usrsrvs
}
func (r *router) Users() *userservice.Users {
	userssrvs := &userservice.Users{
		BaseService: r.baseService,
	}
	return userssrvs
}
