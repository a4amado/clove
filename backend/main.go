package main

import (
	Api "clove/internals/handlers/api"
	"clove/internals/logger"
	_ "embed"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed .env.example
var envExample string

// main is the program entry point.
// It is intentionally empty.
func main() {

	// meridian.Client().ReplicateMessage()
	// postgresPool.Client()
	// valkeyPool.Client(valkeyPool.ValkeyFanout)
	// valkeyPool.Client(valkeyPool.ValkeyHeartbeat)
	// valkeyPool.Client(valkeyPool.ValkeyStore)

	// go meridian.Client().ReplicateMessage().BridgeRabbitMQInternalDeliveryReplicatorToRedis(context.Background())
	router := chi.NewMux()
	router.Use(func(next http.Handler) http.Handler {
		return logger.WrapWithSentry(next.ServeHTTP)
	})
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	router.Mount("/", Api.NewService())
	fmt.Println("listening at :8080")
	http.ListenAndServe(":8080", router)
}
