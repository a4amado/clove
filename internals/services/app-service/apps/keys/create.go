package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"context"
	"time"

	"github.com/google/uuid"
)

func (a *KeysCtx) Create(args repository.CreateAppApiKeyParams) (*repository.AppApiKey, error) {
	q := a.BaseCtx.Queries
	if a.BaseCtx.Tx != nil {
		q = q.WithTx(*a.BaseCtx.Tx)
	}

	key, err := q.CreateAppApiKey(a.BaseCtx.ReqCtx, args)
	if err != nil {
		return nil, err
	}

	if a.BaseCtx.Cache {

		go func(appId uuid.UUID, keyId uuid.UUID, apiKey string) {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
			defer cancel()
			err := cache.Apps().Keys().Set(ctx, appId, keyId, apiKey)
			if err != nil {
				cancel()
				return
			}
		}(key.ID.Bytes, key.ID.Bytes, key.Key.String)

	}

	return &key, nil
}
