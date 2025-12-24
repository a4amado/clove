package mongoDB

import (
	envConsts "clove/internals/consts/env"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var monogdbClinet *mongo.Client
var monogdbClinetOnce = sync.Once{}

// Client returns the package-level singleton MongoDB client, initializing it on first call.
// It connects using the MONGO_HISTORY_DATABASE_URL environment variable and panics if the initial connection fails.
// Each invocation also accesses the database and collection named by MONGO_HISTORY_DATABASE_NAME and MONGO_HISTORY_DATABASE_APP_HISTORY_COLLECTION_NAME.
func Client() *mongo.Client {
	monogdbClinetOnce.Do(func() {
		client, err := mongo.Connect(options.Client().ApplyURI(os.Getenv(string(envConsts.MONGO_HISTORY_DATABASE_URL))))
		if err != nil {
			panic(err)
		}
		monogdbClinet = client
	})

	monogdbClinet.Database(os.Getenv(string(envConsts.MONGO_HISTORY_DATABASE_NAME))).Collection(string(envConsts.MONGO_HISTORY_DATABASE_APP_HISTORY_COLLECTION_NAME))
	return monogdbClinet

}