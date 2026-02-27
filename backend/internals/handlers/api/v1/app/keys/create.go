package AppKeysHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/middleware"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/handlers/api/httpctx"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
	credentialservice "clove/internals/services/credential"
	repository "clove/internals/services/generatedRepo"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/swaggest/usecase"
)

type CreateKeyInput struct {
	AppID uuid.UUID `path:"app_id"`
	Name  string    `json:"name"`
}

type CreateKeyOutput struct {
	repository.AppApiKey
	Token string `json:"token"`
}

func CreateAppApiKey() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input CreateKeyInput, output *CreateKeyOutput) error {
		r := httpctx.Request(ctx)
		session, err := middleware.ParseSession(r, postgresPool.Client())
		if err != nil || !session.Permissions.Can(middleware.KEY, middleware.CREATE) {
			return &apperrors.AppError{
				StatusCode: http.StatusUnauthorized,
				Code:       "UNAUTHORIZED",
			}
		}

		token, err := middleware.GenRandKey(32)
		if err != nil {
			return &apperrors.AppError{StatusCode: http.StatusInternalServerError}
		}

		perms := middleware.NewPermissionsBuilder()
		perms.Allow(middleware.DELIVERY, middleware.CREATE)
		perms.Allow(middleware.OneTimeToken, middleware.CREATE)
		permsStr, _ := perms.String()

		srvs, tx, err := services.New(ctx).WithCache().WithTx()
		if err != nil {
			return &apperrors.AppError{StatusCode: http.StatusInternalServerError}
		}

		_, err = srvs.Credentials.Create(credentialservice.InsertParams{
			Token:       token,
			AppID:       input.AppID,
			Type:        repository.CredentialTypeSdk,
			Permissions: permsStr,
			ExpiresAt:   time.Now().Add(time.Hour * 24 * 365 * 100),
		})
		if err != nil {
			tx.Rollback(ctx)
			return &apperrors.AppError{StatusCode: http.StatusInternalServerError}
		}

		key, err := srvs.Apps.Keys.Create(appservice.CreateKeyParams{
			AppID:   input.AppID,
			KeyName: input.Name,
			Prefix:  token[:min(5, len(token)-1)],
			Suffix:  token[max(0, len(token)-5):],
		})
		if err != nil {
			tx.Rollback(ctx)
			return &apperrors.AppError{StatusCode: http.StatusInternalServerError}
		}

		if err := tx.Commit(ctx); err != nil {
			return &apperrors.AppError{StatusCode: http.StatusInternalServerError}
		}

		output.AppApiKey = *key
		output.Token = token
		return nil
	})
	u.SetTitle("Create API Key")
	u.SetTags("Keys")
	return u
}
