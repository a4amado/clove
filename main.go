package main

import (
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/data/valkeyPool"
	Api "clove/internals/handlers/api"
	"clove/internals/meridian"
	"context"
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
	meridian.Client().ReplicateMessage()
	meridian.Client().ReplicateApp()
	postgresPool.Client()
	valkeyPool.Client(valkeyPool.ValkeyFanout)
	valkeyPool.Client(valkeyPool.ValkeyHeartbeat)
	valkeyPool.Client(valkeyPool.ValkeyStore)

	go meridian.Client().ReplicateApp().BridgeRabbitMQAppReplicatorToRedis(context.Background())
	go meridian.Client().ReplicateMessage().BridgeRabbitMQInternalDeliveryReplicatorToRedis(context.Background())
	router := chi.NewMux()
	router.Mount("/api/", Api.Routes())
	fmt.Println("listening at :8080")
	http.ListenAndServe(":8080", router)
}
