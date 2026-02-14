package AppTokensHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	envConsts "clove/internals/consts/env"
	"clove/internals/handlers/api/httpctx"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swaggest/usecase"
)

type CreateTokenInput struct {
	AppID     uuid.UUID `path:"app_id"`
	ChannelID string    `json:"channel_id"`
}

type CreateTokenOutput struct {
	Token  string `json:"token"`
	Region string `json:"region"`
}

func CreateAppOneTimeToken() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input CreateTokenInput, output *CreateTokenOutput) error {
		r := httpctx.Request(ctx)
		session, err := auth.ParseSessionFromRequest(r)
		if err != nil || !session.Permessions.Can(auth.OneTimeToken, auth.CREATE) {
			return &apperrors.AppError{
				StatusCode: http.StatusUnauthorized,
				Code:       "UNAUTHORIZED",
			}
		}

		appsrvs := services.New(ctx)
		app, err := appsrvs.Apps.Get(appservice.GetParams{
			AppID: input.AppID,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &apperrors.AppError{
					Code:       "ERROR_CREATE_ONE_TIME_TOKEN_APP_NOT_FOUND",
					StatusCode: http.StatusNotFound,
				}
			}
			return &apperrors.AppError{
				StatusCode: http.StatusInternalServerError,
			}
		}

		token, err := auth.GenerateOneTimeToken(*app, input.ChannelID)
		if err != nil {
			return &apperrors.AppError{
				StatusCode: http.StatusInternalServerError,
			}
		}

		output.Token = token
		output.Region = string(envConsts.Region())
		return nil
	})
	u.SetTitle("Create One-Time Token")
	u.SetTags("Tokens")
	return u
}
