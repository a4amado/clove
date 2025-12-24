package envConsts

type EnvKey string
type EnvModeKey string

const (
	DEV  EnvModeKey = "DEV"
	PROD EnvModeKey = "PROD"
)
const (
	REDIS_STORE_URL                                    EnvKey = "REDIS_STORE_URL"
	REDIS_FANOUT_URL                                   EnvKey = "REDIS_FANOUT_URL"
	REDIS_HEARTBEAT_URL                                EnvKey = "REDIS_HEARTBEAT_URL"
	POSTGRES_DATABASE_URL                              EnvKey = "POSTGRES_DATABASE_URL"
	MONGO_HISTORY_DATABASE_URL                         EnvKey = "MONGO_HISTORY_DATABASE_URL"
	MONGO_HISTORY_DATABASE_NAME                        EnvKey = "MONGO_HISTORY_DATABASE_NAME"
	MONGO_HISTORY_DATABASE_USR_HISTORY_COLLECTION_NAME EnvKey = "MONGO_HISTORY_DATABASE_USR_HISTORY_COLLECTION_NAME"
	MONGO_HISTORY_DATABASE_APP_HISTORY_COLLECTION_NAME EnvKey = "MONGO_HISTORY_DATABASE_APP_HISTORY_COLLECTION_NAME"
	REGION                                             EnvKey = "REGION"
	KAFKA_BOOTSTRAP                                    EnvKey = "KAFKA_BOOTSTRAP"
	MAILJET_API_KEY                                    EnvKey = "MAILJET_API_KEY"
	MAILJET_API_SECRETS                                EnvKey = "MAILJET_API_SECRETS"
)
