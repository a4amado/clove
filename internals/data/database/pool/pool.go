package dbPool

import (
	envConsts "clove/internals/consts/env"
	"context"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool
var dbConnOnce = sync.Once{}

func Client() *pgxpool.Pool {
	dbConnOnce.Do(func() {
		pool, err := pgxpool.New(context.Background(), os.Getenv(string(envConsts.POSTGRES_DATABASE_URL)))
		if err != nil {
			panic(err)
		}

		dbPool = pool
	})
	return dbPool
}
