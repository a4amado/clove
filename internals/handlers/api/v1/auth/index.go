package AuthHandlersV1

import (
	"github.com/go-chi/chi/v5"
)

func Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/sign-up", SignUp)
	r.Post("/sign-in", SignIn)
	r.Post("/reset-password", ResetPassword)

	return r
}
