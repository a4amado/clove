package envConsts

import (
	repository "clove/internals/services/generatedRepo"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type EnvType string

const (
	DEV  EnvType = "DEV"
	PROD EnvType = "PROD"
)

var once sync.Once

// loadEnv ensures .env is loaded only once
func loadEnv() {
	once.Do(func() {
		godotenv.Load()
	})
}

// Private helpers
func mustGetString(key string) string {
	loadEnv()
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("environment variable %s is not set", key))
	}
	return val
}

func mustGetInt(key string) int {
	loadEnv()
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("environment variable %s is not set", key))
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s has invalid integer value: %s", key, val))
	}
	return intVal
}

func mustGetFloat(key string) float64 {
	loadEnv()
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("environment variable %s is not set", key))
	}
	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s has invalid float value: %s", key, val))
	}
	return floatVal
}

// Public getters
func RedisStoreURL() string {
	return mustGetString("REDIS_STORE_URL")
}

func RedisFanoutURL() string {
	return mustGetString("REDIS_FANOUT_URL")
}

func RedisHeartbeatURL() string {
	return mustGetString("REDIS_HEARTBEAT_URL")
}

func PostgresDatabaseURL() string {
	return mustGetString("POSTGRES_DATABASE_URL")
}

func MongoHistoryDatabaseURL() string {
	return mustGetString("MONGO_HISTORY_DATABASE_URL")
}

func MongoHistoryDatabaseName() string {
	return mustGetString("MONGO_HISTORY_DATABASE_NAME")
}

func MongoHistoryDatabaseUsrCollectionName() string {
	return mustGetString("MONGO_HISTORY_DATABASE_USR_COLLECTION_NAME")
}

func MongoHistoryDatabaseAppCollectionName() string {
	return mustGetString("MONGO_HISTORY_DATABASE_APP_COLLECTION_NAME")
}

func Region() repository.Region {
	regionType := repository.Region(mustGetString("REGION"))
	if !regionType.Valid() {
		panic(fmt.Sprintf("%s is not a valid regions", mustGetString("REGION")))
	}
	return regionType
}

func KafkaBootstrap() string {
	return mustGetString("KAFKA_BOOTSTRAP")
}

func MailjetAPIKey() string {
	return mustGetString("MAILJET_API_KEY")
}

func MailjetAPISecrets() string {
	return mustGetString("MAILJET_API_SECRETS")
}

func KafkaReaderBufferSize() int {
	size := mustGetInt("KAFKA_READER_BUFFER_SIZE")
	if size <= 0 {
		panic(fmt.Sprintf("KAFKA_READER_BUFFER_SIZE must be greater than 0, got: %d", size))
	}
	return size
}
func KafkaCommitInterval() int {
	size := mustGetInt("KAFKA_COMMIT_INTERVAL")
	if size <= 0 {
		panic(fmt.Sprintf("KAFKA_COMMIT_INTERVAL must be greater than 0, got: %d", size))
	}
	return size
}
func KafkaNumReaders() int {
	nOfReaders := mustGetInt("KAFKA_NUM_READERS")
	if nOfReaders <= 0 {
		panic(fmt.Sprintf("KAFKA_NUM_READERS must be greater than 0, got: %d", nOfReaders))

	}
	return nOfReaders
}
