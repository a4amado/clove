package Api

import (
	"clove/internals/cache"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/data/valkeyPool"
	"clove/internals/handlers/api/httpctx"
	"clove/internals/auth"
	credentialservice "clove/internals/services/credential"
	repository "clove/internals/services/generatedRepo"
	v1 "clove/internals/handlers/api/v1"

	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5cdn"
)

// NewService creates a web.Service with OpenAPI 3.1 documentation and mounts
// all API v1 routes.
func NewService() *web.Service {
	service := web.NewService(openapi31.NewReflector())
	service.OpenAPISchema().SetTitle("Clove API")
	service.OpenAPISchema().SetVersion("v1.0.0")

	credSvc := credentialservice.New(
		repository.New(postgresPool.Client()),
		cache.New(valkeyPool.Client(valkeyPool.ValkeyStore)),
	)

	service.Use(httpctx.Middleware)
	service.Use(auth.SessionMiddleware(credSvc))

	v1.V1Routes(service)

	service.Docs("/docs", swgui.New)
	return service
}
