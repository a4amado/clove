package AppHandlersV1

import "github.com/go-chi/chi/v5"

func Routes() chi.Router {
	router := chi.NewRouter()
	router.Get("/app/:app_id/ws", UserConnect)
	return router
}
