package v1

import (
	AppHandlersV1 "clove/internals/handlers/api/v1/app"
	AuthHandlersV1 "clove/internals/handlers/api/v1/auth"
	UserHandlersV1 "clove/internals/handlers/api/v1/user"

	"github.com/go-chi/chi/v5"
)

// V1Routes creates a chi.Router configured with the v1 API subroutes.
// The returned router mounts the v1 auth, user, and app handlers at /auth, /user, and /app respectively.
func V1Routes() chi.Router {
	r := chi.NewRouter()

	r.Mount("/auth", AuthHandlersV1.Routes())
	r.Mount("/user", UserHandlersV1.Routes())
	r.Mount("/apps", AppHandlersV1.Routes())

	return r
}
