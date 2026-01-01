package types

import (
	repository "clove/internals/services/generatedRepo"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseLineServiceCtx struct {
	Pool    *pgxpool.Pool
	Queries *repository.Queries
	Tx      *pgx.Tx
	ReqCtx  context.Context
	Cache   bool
}
