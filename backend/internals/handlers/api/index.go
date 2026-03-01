package Api

import (
	"net/http"

	"clove/internals/cache"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/data/valkeyPool"
	"clove/internals/handlers/api/httpctx"
	v1 "clove/internals/handlers/api/v1"
	"clove/internals/middleware"
	authservice "clove/internals/services/auth"
	repository "clove/internals/services/generatedRepo"

	swgui "github.com/swaggest/swgui/v5cdn"

	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"
)

// NewService creates a web.Service with OpenAPI 3.1 documentation and mounts
// all API v1 routes.
func NewService() *web.Service {
	reflector := openapi31.NewReflector()
	reflector.Spec.Info.WithTitle("Clove API")
	service := web.NewService(reflector)
	authSvc := authservice.New(
		repository.New(postgresPool.Client()),
		cache.New(valkeyPool.Client(valkeyPool.ValkeyStore)),
	)

	service.Use(httpctx.Middleware)
	service.Use(middleware.SessionMiddleware(authSvc))
	v1.V1Routes(service)
	service.Docs("/api/v1/docs", func(title, specURL, basePath string) http.Handler {
		return swgui.New(title, specURL, basePath)
	})

	return service
}
