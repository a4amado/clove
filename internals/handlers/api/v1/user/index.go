package UserHandlersV1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns a router configured with the user-related API v1 endpoints.
// It currently registers a POST handler for the path "/:user_id".
func Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/:user_id", func(w http.ResponseWriter, r *http.Request) {})

	return r
}