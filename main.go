package main

import (
	envConsts "clove/internals/consts/env"
	dbPool "clove/internals/data/database/pool"
	mongoDB "clove/internals/data/mongo"
	redisPool "clove/internals/data/redispool"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

//go:embed .env.example
var envExample string

// init loads a .env file when ENV equals envConsts.DEV and verifies that required Postgres, MongoDB (database name, URLs, and collection names), and Redis environment variables are set, panicking if any are missing.
func init() {
	if os.Getenv("ENV") != string(envConsts.PROD) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	msg := strings.Builder{}
	for line := range strings.SplitSeq(envExample, "\n") {
		if strings.Index(line, "#") == 0 || len(strings.Trim(line, " ")) == 0 {
			continue
		}
		if os.Getenv(line) == "" {
			fmt.Fprintf(&msg, "%s env is missing\n", line)
		}
	}
	if msg.String() != "" {
		panic(errors.New(msg.String()))
	}

	err := dbPool.Init()
	if err != nil {

	}
	redisPool.Client(redisPool.RedisFanout)
	redisPool.Client(redisPool.RedisHeartbeat)
	redisPool.Client(redisPool.RedisStore)
	mongoDB.Client()

}

// main is the program entry point.
// It is intentionally empty.
func main() {

}
