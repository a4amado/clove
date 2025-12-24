package AppService

import (
	dbPool "clove/internals/data/database/pool"
	"clove/internals/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppService struct {
	Pool    *pgxpool.Pool
	Queries *repository.Queries
}

func New() *AppService {
	newPool := dbPool.Client()
	return &AppService{
		Pool:    newPool,
		Queries: repository.New(newPool),
	}
}
