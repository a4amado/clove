package AppHandlersV1

import (
	AppKeysHandlersV1 "clove/internals/handlers/api/v1/app/keys"
	AppTokensHandlersV1 "clove/internals/handlers/api/v1/app/tokens"

	"github.com/go-chi/chi/v5"
)

// Routes creates a chi.Router configured with handlers for version 1 of the app API.
// It registers a GET route at /app/{app_id}/ws handled by UserConnect and returns the router.
func Routes() chi.Router {
	router := chi.NewRouter()
	router.Post("/", CreateApp)
	router.Get("/{app_id}/ws/", UserConnect)
	router.Post("/{app_id}/entry/", MessageEntry)
	router.Post("/{app_id}/auth/", MessageEntry)

	router.Mount("/{app_id}/keys/", AppKeysHandlersV1.Routes())
	router.Mount("/{app_id}/tokens/", AppTokensHandlersV1.Routes())
	return router
}
