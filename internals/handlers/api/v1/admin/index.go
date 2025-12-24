package ApiV1Admin

import "github.com/go-chi/chi/v5"

// Routes creates and returns a new chi.Router configured for the API v1 admin routes.
func Routes() chi.Router {
	r := chi.NewRouter()
	return r
}