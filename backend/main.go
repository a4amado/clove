package main

import (
	"clove/internals/auth"
	_ "embed"
	"fmt"
	"time"

	"github.com/google/uuid"
)

//go:embed .env.example
var envExample string

// main is the program entry point.
// It is intentionally empty.
func main() {

	start := time.Now()
	token, _ := auth.GenerateSDKToken(uuid.New())
	generated := time.Now()
	c, _ := auth.ValidateSDKToken(token)
	validated := time.Now()
	fmt.Println("(true)", c.Permessions.Can(auth.DELIVERY, auth.CREATE))
	fmt.Println("(false)", c.Permessions.Can(auth.DELIVERY, auth.READ))

	fmt.Println("Generate:", generated.Sub(start))
	fmt.Println("Validate:", validated.Sub(generated))
	// meridian.Client().ReplicateMessage()
	// // meridian.Client().ReplicateApp()
	// // postgresPool.Client()
	// // valkeyPool.Client(valkeyPool.ValkeyFanout)
	// // valkeyPool.Client(valkeyPool.ValkeyHeartbeat)
	// // valkeyPool.Client(valkeyPool.ValkeyStore)

	// // go meridian.Client().ReplicateApp().BridgeRabbitMQAppReplicatorToRedis(context.Background())
	// // go meridian.Client().ReplicateMessage().BridgeRabbitMQInternalDeliveryReplicatorToRedis(context.Background())
	// router := chi.NewMux()
	// router.Use(func(next http.Handler) http.Handler {
	// 	return logger.WrapWithSentry(next.ServeHTTP)
	// })
	// router.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(404)
	// })
	// // router.Mount("/api/", Api.Routes())
	// fmt.Println("listening at :8080")
	// http.ListenAndServe(":8080", router)
}
