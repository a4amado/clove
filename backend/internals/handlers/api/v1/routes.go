package v1

import (
	"clove/internals/auth"
	AppHandlersV1 "clove/internals/handlers/api/v1/app"
	AppKeysHandlersV1 "clove/internals/handlers/api/v1/app/keys"
	AppRegionsHandlersV1 "clove/internals/handlers/api/v1/app/regions"
	AppTokensHandlersV1 "clove/internals/handlers/api/v1/app/tokens"
	AuthHandlersV1 "clove/internals/handlers/api/v1/auth"

	"github.com/go-chi/chi/v5"
)

// V1Routes creates a chi.Router configured with the v1 API subroutes.
// The returned router mounts the v1 auth, user, and app handlers at /auth, /user, and /app respectively.
func V1Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", AuthHandlersV1.User_Signup)
		r.Post("/sign-in", AuthHandlersV1.SignIn)
	})

	r.Route("/apps", func(r chi.Router) {
		r.Route("/{app_id}/", func(r chi.Router) {
			r.Get("/ws/", AppHandlersV1.UserConnect) // this will use one tiem token, should i Hanle it in Auth Middlware and use the create another middle ware to enfore tir for each route, and what is the most maintanable way ?
			r.With(auth.AuthMiddleware).Post("/entry/", AppHandlersV1.WSMessageEntry)
			r.With(auth.AuthMiddleware).Route("/keys/", func(r chi.Router) {
				r.Get("/", AppKeysHandlersV1.ListAppApiKeys)
				r.Post("/", AppKeysHandlersV1.CreateAppApiKey)
				r.Route("/{key_id}/", func(r chi.Router) {
					r.Delete("/", AppKeysHandlersV1.DeleteAppApiKey)
				})
			})
			r.With(auth.AuthMiddleware).Route("/tokens/", func(r chi.Router) {
				r.Post("/", AppTokensHandlersV1.CreateAppOneTimeToken)
			})
			r.With(auth.AuthMiddleware).Route("/regions/", func(r chi.Router) {
				r.Get("/", AppRegionsHandlersV1.ListAppRegions)
				r.Patch("/", AppRegionsHandlersV1.UpdateAppRegions)
			})
		})

	})

	return r
}
