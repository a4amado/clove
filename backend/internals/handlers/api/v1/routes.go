package v1

import (
	"clove/internals/cache"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/data/valkeyPool"
	AppHandlersV1 "clove/internals/handlers/api/v1/app"
	AppKeysHandlersV1 "clove/internals/handlers/api/v1/app/keys"
	AppRegionsHandlersV1 "clove/internals/handlers/api/v1/app/regions"
	AppTokensHandlersV1 "clove/internals/handlers/api/v1/app/tokens"
	AuthHandlersV1 "clove/internals/handlers/api/v1/auth"
	"clove/internals/middleware"
	authservice "clove/internals/services/auth"
	repository "clove/internals/services/generatedRepo"
	"net/http"

	"github.com/swaggest/rest/web"
)

// V1Routes registers all v1 API routes on the given web.Service.
func V1Routes(service *web.Service) {
	authSvc := authservice.New(
		repository.New(postgresPool.Client()),
		cache.New(valkeyPool.Client(valkeyPool.ValkeyStore)),
	)
	authMiddleware := middleware.AuthMiddleware(authSvc)

	// Auth routes (public)
	service.Post("/api/v1/auth/sign-up", AuthHandlersV1.Signup())
	service.Post("/api/v1/auth/sign-in", AuthHandlersV1.SignIn())
	service.Get("/api/v1/auth/me", AuthHandlersV1.Me())

	// App management — session populated by global SessionMiddleware in index.go
	service.Post("/api/v1/apps", AppHandlersV1.CreateApp())
	service.Get("/api/v1/apps", AppHandlersV1.ListApps())
	service.Get("/api/v1/apps/{app_id}/keys", AppKeysHandlersV1.ListAppApiKeys())
	service.Post("/api/v1/apps/{app_id}/keys", AppKeysHandlersV1.CreateAppApiKey())
	service.Delete("/api/v1/apps/{app_id}/keys/{key_id}", AppKeysHandlersV1.DeleteAppApiKey())
	service.Post("/api/v1/apps/{app_id}/tokens", AppTokensHandlersV1.CreateAppOneTimeToken())
	service.Get("/api/v1/apps/{app_id}/regions", AppRegionsHandlersV1.ListAppRegions())
	service.Patch("/api/v1/apps/{app_id}/regions", AppRegionsHandlersV1.UpdateAppRegions())

	// Raw handlers (WebSocket + raw-body entry) — auth middleware blocks unauthenticated requests
	service.Method(http.MethodGet, "/api/v1/apps/{app_id}/ws",
		authMiddleware(http.HandlerFunc(AppHandlersV1.UserConnect)))
	service.Method(http.MethodPost, "/api/v1/apps/{app_id}/entry",
		authMiddleware(http.HandlerFunc(AppHandlersV1.WSMessageEntry)))
}
