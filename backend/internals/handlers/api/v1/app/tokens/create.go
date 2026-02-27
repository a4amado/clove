package AppTokensHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/middleware"
	postgresPool "clove/internals/data/postgres/pool"
	envConsts "clove/internals/consts/env"
	"clove/internals/handlers/api/httpctx"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
	credentialservice "clove/internals/services/credential"
	repository "clove/internals/services/generatedRepo"
	"context"
	"errors"
	"net/http"
	"time"

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
		session, err := middleware.ParseSession(r, postgresPool.Client())
		if err != nil || !session.Permissions.Can(middleware.OneTimeToken, middleware.CREATE) {
			return &apperrors.AppError{
				StatusCode: http.StatusUnauthorized,
				Code:       "UNAUTHORIZED",
			}
		}

		srvs := services.New(ctx)
		_, err = srvs.Apps.Get(appservice.GetParams{AppID: input.AppID})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &apperrors.AppError{
					Code:       "ERROR_CREATE_ONE_TIME_TOKEN_APP_NOT_FOUND",
					StatusCode: http.StatusNotFound,
				}
			}
			return &apperrors.AppError{StatusCode: http.StatusInternalServerError}
		}

		token, err := middleware.GenRandKey(32)
		if err != nil {
			return &apperrors.AppError{StatusCode: http.StatusInternalServerError}
		}

		perms := middleware.NewPermissionsBuilder()
		perms.Allow(middleware.DELIVERY, middleware.READ)
		permsStr, _ := perms.String()

		_, err = srvs.Credentials.Create(credentialservice.InsertParams{
			Token:       token,
			AppID:       input.AppID,
			ChannelID:   input.ChannelID,
			Type:        repository.CredentialTypeOtt,
			Permissions: permsStr,
			ExpiresAt:   time.Now().Add(time.Minute),
		})
		if err != nil {
			return &apperrors.AppError{StatusCode: http.StatusInternalServerError}
		}

		output.Token = token
		output.Region = string(envConsts.Region())
		return nil
	})
	u.SetTitle("Create One-Time Token")
	u.SetTags("Tokens")
	return u
}
