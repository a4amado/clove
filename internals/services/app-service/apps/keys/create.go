package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"context"
	"time"

	"github.com/google/uuid"
)

func (a *KeysCtx) Generate(args repository.CreateAppApiKeyParams) (*repository.AppApiKey, error) {
	q := a.App.Queries
	if a.App.Tx != nil {
		q = q.WithTx(*a.App.Tx)
	}

	key, err := q.CreateAppApiKey(a.App.ReqCtx, args)
	if err != nil {
		return nil, err
	}

	if a.App.Cache {

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
