package appservice

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (a *KeyCtx) Get() (*pgtype.Text, error) {
	q := a.App.Queries
	if a.App.Tx != nil {
		q = q.WithTx(*a.App.Tx)
	}

	if a.App.Cache {
		res, err := cache.Apps().Keys().Get(a.App.ReqCtx, a.AppId, a.KeyId)
		if err == nil {
			return &pgtype.Text{
				String: *res,
				Valid:  true,
			}, nil
		}
		key, err := q.GetAppApiKey(a.App.ReqCtx, repository.GetAppApiKeyParams{
			KeyID: pgtype.UUID{
				Bytes: a.KeyId,
				Valid: true,
			},
			AppID: pgtype.UUID{
				Bytes: a.AppId,
				Valid: true,
			},
		})
		if err != nil {
			return nil, err
		}
		go func(args repository.GetAppApiKeyParams, key pgtype.Text) {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(500*time.Millisecond))
			defer cancel()
			err := cache.Apps().Keys().Set(ctx, args.AppID.Bytes, args.KeyID.Bytes, key.String)
			if err != nil {

				return
			}
		}(repository.GetAppApiKeyParams{
			KeyID: pgtype.UUID{
				Bytes: a.KeyId,
				Valid: true,
			},
			AppID: pgtype.UUID{
				Bytes: a.AppId,
				Valid: true,
			},
		}, key)
		return &pgtype.Text{
			String: key.String,
			Valid:  true,
		}, nil
	}
	key, err := q.GetAppApiKey(a.App.ReqCtx, repository.GetAppApiKeyParams{
		KeyID: pgtype.UUID{
			Bytes: a.KeyId,
			Valid: true,
		},
		AppID: pgtype.UUID{
			Bytes: a.AppId,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}
	go func(args repository.GetAppApiKeyParams, key pgtype.Text) {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(500*time.Millisecond))
		defer cancel()
		err := cache.Apps().Keys().Set(ctx, args.AppID.Bytes, args.KeyID.Bytes, key.String)
		if err != nil {

			return
		}
	}(repository.
		GetAppApiKeyParams{
		KeyID: pgtype.UUID{
			Bytes: a.KeyId,
			Valid: true,
		},
		AppID: pgtype.UUID{
			Bytes: a.AppId,
			Valid: true,
		},
	}, key)
	return &pgtype.Text{
		String: key.String,
		Valid:  true,
	}, nil
}
