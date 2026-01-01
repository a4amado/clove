package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateAppResponse struct {
	App repository.App       `json:"app"`
	Key repository.AppApiKey `json:"key"`
}

func (as *AppsCtx) Create(arg repository.InsertAppParams) (*CreateAppResponse, error) {
	queries := as.BaseCtx.Queries
	if as.BaseCtx.Tx != nil {
		queries = queries.WithTx(*as.BaseCtx.Tx)
	}
	app, err := queries.InsertApp(as.BaseCtx.ReqCtx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to insert app: %w", err)
	}

	key, err := queries.CreateAppApiKey(as.BaseCtx.ReqCtx, repository.CreateAppApiKeyParams{
		AppID: app.ID,
		Key: pgtype.Text{
			String: uuid.NewString(),
			Valid:  true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return &CreateAppResponse{
		App: app,
		Key: key,
	}, nil
}
