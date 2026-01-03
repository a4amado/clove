package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (a *KeysCtx) Create(name string, keyString string) (*repository.AppApiKey, error) {
	q := a.BaseCtx.Queries
	if a.BaseCtx.Tx != nil {
		q = q.WithTx(*a.BaseCtx.Tx)
	}

	key, err := q.App_Key_Insert(a.BaseCtx.ReqCtx, repository.App_Key_InsertParams{
		AppID: pgtype.UUID{Bytes: a.AppID, Valid: true},
		Key: pgtype.Text{
			String: keyString,
			Valid:  true,
		},
		Name: pgtype.Text{
			String: name,
			Valid:  true,
		},
	})
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
