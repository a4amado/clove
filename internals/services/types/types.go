package types

import (
	repository "clove/internals/services/generatedRepo"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ServiceParams struct {
	Ctx      context.Context
	Tx       *pgx.Tx
	UseCache bool
}

func NewBaseService(ctx context.Context, queries *repository.Queries, useCache bool) *BaseService {

	return &BaseService{
		ctx:      ctx,
		queries:  queries,
		useCache: useCache,
	}
}

type BaseService struct {
	ctx      context.Context
	queries  *repository.Queries
	useCache bool
	tx       pgx.Tx
}

func (b *BaseService) Cache() bool {
	return b.useCache
}
func (b *BaseService) CTX() context.Context {
	return b.ctx
}

func (b *BaseService) Q() *repository.Queries {
	return b.queries
}

func (b *BaseService) ToPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func (b *BaseService) ToPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

// cacheAsync runs cache operations in background with proper timeout
func (b *BaseService) CacheAsync(fn func(context.Context) error) {
	if !b.useCache {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		if err := fn(ctx); err != nil {
			// TODO: Add structured logging here
			// log.Error("cache operation failed", "error", err)
		}
	}()
}
