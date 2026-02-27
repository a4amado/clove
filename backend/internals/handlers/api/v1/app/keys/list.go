package AppKeysHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/middleware"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/handlers/api/httpctx"
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
		r := httpctx.Request(ctx)
		session, err := middleware.ParseSession(r, postgresPool.Client())
		if err != nil || !session.Permissions.Can(middleware.KEY, middleware.READ) {
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
	u.SetTitle("List API Keys")
	u.SetTags("Keys")
	return u
}
