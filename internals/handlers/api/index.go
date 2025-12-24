package Api

import (
	ApiV1 "clove/internals/handlers/api/v1"

	"github.com/go-chi/chi/v5"
)

func Routes() chi.Router {
	route := chi.NewRouter()
	route.Mount("/v1", ApiV1.V1Routes())
	return route
}
