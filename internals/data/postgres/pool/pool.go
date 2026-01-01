package postgresPool

import (
	envConsts "clove/internals/consts/env"
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool
var dbPoolOnce = sync.Once{}

func NewTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return dbPool.BeginTx(ctx, txOptions)
}

// Client returns the package's singleton PostgreSQL connection pool.
// It lazily initializes the pool on first call using the POSTGRES_DATABASE_URL
// environment variable and panics if pool creation fails.
func Client() *pgxpool.Pool {
	dbPoolOnce.Do(func() {
		config, _ := pgxpool.ParseConfig(envConsts.PostgresDatabaseURL())
		config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			// Load the base enum
			regionType, err := conn.LoadType(ctx, "region")
			if err != nil {
				return fmt.Errorf("failed to load region type: %w", err)
			}
			conn.TypeMap().RegisterType(regionType)

			// Load the array type (note the underscore prefix)
			regionArrayType, err := conn.LoadType(ctx, "_region")
			if err != nil {
				return fmt.Errorf("failed to load _region type: %w", err)
			}
			conn.TypeMap().RegisterType(regionArrayType)

			// Do the same for app_type if you use it in queries
			appType, err := conn.LoadType(ctx, "app_type")
			if err != nil {
				return fmt.Errorf("failed to load app_type: %w", err)
			}
			conn.TypeMap().RegisterType(appType)

			UserRole, err := conn.LoadType(ctx, "UserRole")
			if err != nil {
				return fmt.Errorf("failed to load UserRole: %w", err)
			}
			conn.TypeMap().RegisterType(UserRole)

			return nil
		}
		pool, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {

		}

		dbPool = pool
	})
	return dbPool
}
