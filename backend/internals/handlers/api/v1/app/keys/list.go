package AppKeysHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
	repository "clove/internals/services/generatedRepo"
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/swaggest/usecase"
)

type ListKeysInput struct {
	AppID   uuid.UUID `path:"app_id"`
	PageIdx int64     `query:"page_idx"`
}

type ListKeysOutput struct {
	Keys []repository.AppApiKey `json:"keys"`
}

func ListAppApiKeys() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input ListKeysInput, output *ListKeysOutput) error {
		session, ok := auth.SessionFromContext(ctx)
		if !ok || !session.Permissions.Can(auth.KEY, auth.READ) {
			return &apperrors.AppError{
				StatusCode: http.StatusUnauthorized,
				Code:       "UNAUTHORIZED",
			}
		}

		srvs := services.New(ctx).WithCache()
		keys, err := srvs.Apps.Keys.List(appservice.ListKeysParams{
			AppId: input.AppID,
			Page:  int32(input.PageIdx),
		})
		if err != nil {
			return &apperrors.AppError{
				StatusCode: http.StatusInternalServerError,
			}
		}

		output.Keys = keys
		return nil
	})
	u.SetExpectedErrors(apperrors.ErrUnauthorized, apperrors.ErrInternalServerError)
	u.SetTitle("List API Keys")
	u.SetTags("Keys")
	return u
}
