package main

import (
	mongoDB "clove/internals/data/mongo"
	postgresPool "clove/internals/data/postgres/pool"
	redisPool "clove/internals/data/redispool"
	emailTemplates "clove/internals/email/email-templates"
	Api "clove/internals/handlers/api"
	"clove/internals/meridian"
	"clove/internals/meridian/fanout"
	"clove/internals/repository"
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

//go:embed .env.example
var envExample string

// main is the program entry point.
// It is intentionally empty.
func main() {
	// Load config FIRST - make it obvious
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env:", err)
	}

	postgresPool.Init()
	redisPool.Init()
	mongoDB.Init()
	emailTemplates.Init()

	replicateClient := meridian.Client().ReplicateApp()

	fanoutClient := meridian.Client().Fanout()

	go replicateClient.BridgeKafkaAppReplicatorToRedis(context.Background())
	user, err := repository.New(postgresPool.Client()).InsertUser(context.Background(), repository.InsertUserParams{
		Email: uuid.NewString() + "a4addel@gmail.com",
		Hash:  "ssssssssssssssssss",
	})
	if err != nil {
		panic(err)
	}
	app, err := repository.New(postgresPool.Client()).InsertApp(context.Background(), repository.InsertAppParams{
		Appslug: uuid.NewString() + "test",
		Region:  []repository.Region{repository.RegionDk1},
		Apptype: repository.AppTypeFree,
		Userid:  user.ID,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(app.ID.String())

	err = replicateClient.SaveApp(context.Background(), app)
	if err != nil {
		panic(err)
	}
	fmt.Println("sending to: ", fanoutClient.FormatChannelKey(fanout.ChannelKey{
		AppId:     app.ID.Bytes,
		ChannelId: "test",
	}))
	go meridian.Client().ReplicateMessage().BridgeKafkaInternalDelevieryReplicatorToRedis(context.Background())
	go func() {
		t := time.NewTicker(time.Second)
		for {
			<-t.C

		}

	}()
	chi := chi.NewMux()
	chi.Mount("/api/", Api.Routes())
	fmt.Println("listening at :3000")
	http.ListenAndServe(":3000", chi)
}
