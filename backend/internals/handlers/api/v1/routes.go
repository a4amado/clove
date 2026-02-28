package v1

import (
	"clove/internals/auth"
	"clove/internals/cache"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/data/valkeyPool"
	AppHandlersV1 "clove/internals/handlers/api/v1/app"
	AppKeysHandlersV1 "clove/internals/handlers/api/v1/app/keys"
	AppRegionsHandlersV1 "clove/internals/handlers/api/v1/app/regions"
	AppTokensHandlersV1 "clove/internals/handlers/api/v1/app/tokens"
	AuthHandlersV1 "clove/internals/handlers/api/v1/auth"
	credentialservice "clove/internals/services/credential"
	repository "clove/internals/services/generatedRepo"
	"net/http"

	"github.com/swaggest/rest/web"
)

// V1Routes registers all v1 API routes on the given web.Service.
func V1Routes(service *web.Service) {
	credSvc := credentialservice.New(
		repository.New(postgresPool.Client()),
		cache.New(valkeyPool.Client(valkeyPool.ValkeyStore)),
	)
	authMiddleware := auth.AuthMiddleware(credSvc)

	// Auth routes (public)
	service.Post("/v1/auth/sign-up", AuthHandlersV1.Signup())
	service.Post("/v1/auth/sign-in", AuthHandlersV1.SignIn())
	service.Get("/v1/auth/me", AuthHandlersV1.Me())

	// App management — session populated by global SessionMiddleware in index.go
	service.Post("/v1/apps", AppHandlersV1.CreateApp())
	service.Get("/v1/apps/{app_id}/keys", AppKeysHandlersV1.ListAppApiKeys())
	service.Post("/v1/apps/{app_id}/keys", AppKeysHandlersV1.CreateAppApiKey())
	service.Delete("/v1/apps/{app_id}/keys/{key_id}", AppKeysHandlersV1.DeleteAppApiKey())
	service.Post("/v1/apps/{app_id}/tokens", AppTokensHandlersV1.CreateAppOneTimeToken())
	service.Get("/v1/apps/{app_id}/regions", AppRegionsHandlersV1.ListAppRegions())
	service.Patch("/v1/apps/{app_id}/regions", AppRegionsHandlersV1.UpdateAppRegions())

	// Raw handlers (WebSocket + raw-body entry) — auth middleware blocks unauthenticated requests
	service.Method(http.MethodGet, "/v1/apps/{app_id}/ws",
		authMiddleware(http.HandlerFunc(AppHandlersV1.UserConnect)))
	service.Method(http.MethodPost, "/v1/apps/{app_id}/entry",
		authMiddleware(http.HandlerFunc(AppHandlersV1.WSMessageEntry)))
}
