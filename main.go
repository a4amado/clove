package main

import (
	envConsts "clove/internals/consts/env"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("ENV") == string(envConsts.DEV) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	if os.Getenv(string(envConsts.POSTGRES_DATABASE_URL)) == "" {
		panic(errors.New(string(envConsts.POSTGRES_DATABASE_URL) + " env is missig"))
	}
	if os.Getenv(string(envConsts.MONGO_HISTORY_DATABASE_APP_HISTORY_COLLECTION_NAME)) == "" {
		panic(errors.New(string(envConsts.MONGO_HISTORY_DATABASE_APP_HISTORY_COLLECTION_NAME) + " env is missig"))

	}
	if os.Getenv(string(envConsts.MONGO_HISTORY_DATABASE_NAME)) == "" {
		panic(errors.New(string(envConsts.MONGO_HISTORY_DATABASE_NAME) + " env is missig"))

	}
	if os.Getenv(string(envConsts.MONGO_HISTORY_DATABASE_URL)) == "" {
		panic(string(envConsts.MONGO_HISTORY_DATABASE_URL) + " env is missig")

	}
	if os.Getenv(string(envConsts.MONGO_HISTORY_DATABASE_USR_HISTORY_COLLECTION_NAME)) == "" {
		panic(errors.New(string(envConsts.MONGO_HISTORY_DATABASE_USR_HISTORY_COLLECTION_NAME) + " env is missig"))

	}
	if os.Getenv(string(envConsts.REDIS_FANOUT_URL)) == "" {
		panic(errors.New(string(envConsts.REDIS_FANOUT_URL) + " env is missig"))

	}
	if os.Getenv(string(envConsts.REDIS_HEARTBEAT_URL)) == "" {
		panic(errors.New(string(envConsts.REDIS_HEARTBEAT_URL) + " env is missig"))

	}
	if os.Getenv(string(envConsts.REDIS_STORE_URL)) == "" {
		panic(errors.New(string(envConsts.REDIS_STORE_URL) + " env is missig"))

	}

}

func main() {

}
