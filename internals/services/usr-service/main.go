package userservice

import (
	repository "clove/internals/services/generatedRepo"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	Pool    *pgxpool.Pool
	Queries *repository.Queries
	tx      *pgx.Tx
	ctx     context.Context
}

func (us *UserService) WithTx(tx pgx.Tx) {
	us.tx = &tx
}
func (us *UserService) WithContext(ctx context.Context) {
	us.ctx = ctx
}
func (us *UserService) EnsureTransactionAndContext() {
	if us.tx == nil {
		panic("each query need to be wrapped in a trasaction")
	}
	if us.ctx == nil {
		panic("each query needs to be in a context")
	}
}
