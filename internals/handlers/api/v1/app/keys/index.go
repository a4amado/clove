package AppKeysHandlersV1

import (
	"github.com/go-chi/chi/v5"
)

// Routes creates a chi.Router configured with handlers for version 1 of the app API.
// It registers a GET route at /app/{app_id}/ws handled by UserConnect and returns the router.
func Routes() chi.Router {
	router := chi.NewRouter()
	router.Get("/", ListAppApiKeys)
	router.Post("/", CreateAppApiKey)
	router.Delete("/{api_token_id}/", DeleteAppApiKey)
	return router
}
