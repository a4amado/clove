package main

import (
	"clove/internals/auth"
	repository "clove/internals/services/generatedRepo"
	_ "embed"
	"log"
)

//go:embed .env.example
var envExample string

// main is the program entry point.
// It is intentionally empty.
func main() {

	claims := auth.NewClaimsBuilder()
	claims.Allow(repository.ResourceAPP, repository.OperationCREATE)
	token, _ := claims.String()
	log.Println("token", token)
	parsedClaims, _ := auth.ParseClaims(token)
	log.Println("can (true):", parsedClaims.Can(repository.ResourceAPP, repository.OperationCREATE))
	log.Println("can (false):", parsedClaims.Can(repository.ResourceKEY, repository.OperationCREATE))
	// meridian.Client().ReplicateMessage()
	// meridian.Client().ReplicateApp()
	// postgresPool.Client()
	// valkeyPool.Client(valkeyPool.ValkeyFanout)
	// valkeyPool.Client(valkeyPool.ValkeyHeartbeat)
	// valkeyPool.Client(valkeyPool.ValkeyStore)

	// go meridian.Client().ReplicateApp().BridgeRabbitMQAppReplicatorToRedis(context.Background())
	// go meridian.Client().ReplicateMessage().BridgeRabbitMQInternalDeliveryReplicatorToRedis(context.Background())
	// router := chi.NewMux()
	// router.Use(func(next http.Handler) http.Handler {
	// 	return logger.WrapWithSentry(next.ServeHTTP)
	// })

	// router.Mount("/api/", Api.Routes())
	// fmt.Println("listening at :8080")
	// http.ListenAndServe(":8080", router)
}
