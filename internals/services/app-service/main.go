package AppService

import (
	dbPool "clove/internals/data/postgres/pool"
	"clove/internals/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppService struct {
	Pool    *pgxpool.Pool
	Queries *repository.Queries
}

// New constructs and returns an AppService with an initialized database pool and query helper.
// The returned AppService has Pool set to a database connection pool and Queries initialized to use that pool.
func New() *AppService {
	newPool := dbPool.Client()
	return &AppService{
		Pool:    newPool,
		Queries: repository.New(newPool),
	}
}
