package AppKeysHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	"clove/internals/handlers/api/httpctx"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/swaggest/usecase"
)

type DeleteKeyInput struct {
	AppID uuid.UUID `path:"app_id"`
	KeyID uuid.UUID `path:"key_id"`
}

type DeleteKeyOutput struct{}

func DeleteAppApiKey() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input DeleteKeyInput, output *DeleteKeyOutput) error {
		r := httpctx.Request(ctx)
		session, err := auth.ParseSessionFromRequest(r)
		if err != nil || !session.Permessions.Can(auth.KEY, auth.DESTROY) {
			return &apperrors.AppError{
				StatusCode: http.StatusUnauthorized,
				Code:       "UNAUTHORIZED",
			}
		}

		srvs, tx, err := services.New(ctx).WithCache().WithTx()
		if err != nil {
			return &apperrors.AppError{
				StatusCode: http.StatusInternalServerError,
			}
		}
		err = srvs.Apps.Keys.Delete(appservice.DeleteKeyParams{
			AppID: input.AppID,
			KeyID: input.KeyID,
		})
		if err != nil {
			tx.Rollback(ctx)
			return &apperrors.AppError{
				StatusCode: http.StatusBadRequest,
			}
		}

		if err = tx.Commit(ctx); err != nil {
			return &apperrors.AppError{
				StatusCode: http.StatusInternalServerError,
			}
		}

		return nil
	})
	u.SetTitle("Delete API Key")
	u.SetTags("Keys")
	return u
}
