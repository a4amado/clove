package UserHandlersV1

import (
	"github.com/go-chi/chi/v5"
)

// Routes returns a router configured with the user-related API v1 endpoints.
// It currently registers a POST handler for the path "/{user_id}".
func Routes() chi.Router {
	r := chi.NewRouter()

	r.Patch("/users/{user_id}", UpdateUser)

	return r
}
