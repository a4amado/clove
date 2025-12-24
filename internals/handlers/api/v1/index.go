package v1

import (
	AppHandlersV1 "clove/internals/handlers/api/v1/app"
	AuthHandlersV1 "clove/internals/handlers/api/v1/auth"
	UserHandlersV1 "clove/internals/handlers/api/v1/user"

	"github.com/go-chi/chi/v5"
)

func V1Routes() chi.Router {
	r := chi.NewRouter()

	r.Mount("/auth", AuthHandlersV1.Routes())
	r.Mount("/user", UserHandlersV1.Routes())
	r.Mount("/app", AppHandlersV1.Routes())

	return r
}
