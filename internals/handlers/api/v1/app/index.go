package AppHandlersV1

import "github.com/go-chi/chi/v5"

// Routes creates a chi.Router configured with handlers for version 1 of the app API.
// Routes creates and configures a chi.Router with API v1 routes.
// It registers a GET route at /{app_id}/ws and a POST route at /{app_id}/channels/{channel_id}, both handled by UserConnect, and returns the configured router.
func Routes() chi.Router {
	router := chi.NewRouter()
	router.Get("/{app_id}/ws", UserConnect)
	router.Post("/{app_id}/channels/{channel_id}", UserConnect)
	return router
}