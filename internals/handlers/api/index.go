package Api

import (
	ApiV1 "clove/internals/handlers/api/v1"

	"github.com/go-chi/chi/v5"
)

// Routes constructs a chi.Router and mounts the API v1 routes at the "/v1" path.
// The returned router is ready to be used as the top-level handler for the API.
func Routes() chi.Router {
	route := chi.NewRouter()
	route.Mount("/v1", ApiV1.V1Routes())
	return route
}