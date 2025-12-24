package dbPool

import (
	envConsts "clove/internals/consts/env"
	"context"
	"errors"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool
var dbConnOnce = sync.Once{}

func Init() error {
	var ConnectionError error
	dbConnOnce.Do(func() {
		pool, err := pgxpool.New(context.Background(), os.Getenv(string(envConsts.POSTGRES_DATABASE_URL)))
		if err != nil {
			ConnectionError = err
		}

		dbPool = pool
	})
	return ConnectionError
}

// Client returns the package's singleton PostgreSQL connection pool.
// It lazily initializes the pool on first call using the POSTGRES_DATABASE_URL
// environment variable and panics if pool creation fails.
func Client() *pgxpool.Pool {
	if dbPool == nil {
		panic(errors.New("dbPool has not been initialized"))
	}
	return dbPool
}
