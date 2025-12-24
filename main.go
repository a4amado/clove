package main

import (
	envConsts "clove/internals/consts/env"
	mongoDB "clove/internals/data/mongo"
	postgresPool "clove/internals/data/postgres/pool"
	redisPool "clove/internals/data/redispool"
	emailTemplates "clove/internals/email/email-templates"
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
		varName := line
		if idx := strings.Index(line, "="); idx > 0 {
			varName = line[:idx]
		}
		varName = strings.TrimSpace(varName)
		if os.Getenv(varName) == "" {
			fmt.Fprintf(&msg, "%s env is missing\n", varName)
		}
	}
	if msg.String() != "" {
		panic(errors.New(msg.String()))
	}

	// paincs on startup if any of these failed
	postgresPool.Init()
	redisPool.Init()
	mongoDB.Init()
	emailTemplates.Init()
}

// main is the program entry point.
// It is intentionally empty.
func main() {

}
