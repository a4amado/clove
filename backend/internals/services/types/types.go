package types

import (
	"clove/internals/cache"
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

// NewBaseService creates a BaseService. The queries parameter is ignored (set DB directly).
func NewBaseService(ctx context.Context, _ *repository.Queries, useCache bool) *BaseService {
	return &BaseService{
		ctx:      ctx,
		useCache: useCache,
	}
}

type BaseService struct {
	ctx      context.Context
	DB       *repository.Queries
	Cache    *cache.Client
	useCache bool
}

// WithCache enables cache and sets the cache client.
func (b *BaseService) WithCache(client *cache.Client) *BaseService {
	b.useCache = true
	b.Cache = client
	return b
}

// IsCache reports whether cache is enabled.
func (b *BaseService) IsCache() bool {
	return b.useCache
}

// CacheReady reports whether a cache client is available and cache is enabled.
func (b *BaseService) CacheReady() bool {
	return b.useCache && b.Cache != nil
}

func (b *BaseService) GetCtx() context.Context {
	return b.ctx
}

func (b *BaseService) WithCtx(ctx context.Context) *BaseService {
	b.ctx = ctx
	return b
}

func (b *BaseService) ToPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func (b *BaseService) ToPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

// CacheAsync runs a cache operation in the background with a short timeout.
// No-ops if cache is not enabled or no client is set.
func (b *BaseService) CacheAsync(fn func(context.Context) error) {
	if !b.IsCache() {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		fn(ctx) //nolint:errcheck
	}()
}
