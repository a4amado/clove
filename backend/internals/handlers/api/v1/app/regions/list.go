package AppRegionsHandlersV1

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

type ListRegionsInput struct {
	AppID uuid.UUID `path:"app_id"`
}

type ListRegionsOutput struct {
	Regions []repository.Region `json:"regions"`
}

func ListAppRegions() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input ListRegionsInput, output *ListRegionsOutput) error {
		r := httpctx.Request(ctx)
		session, err := middleware.ParseSession(r, postgresPool.Client())
		if err != nil || !session.Permissions.Can(middleware.KEY, middleware.READ) {
			return &apperrors.AppError{
				StatusCode: http.StatusUnauthorized,
				Code:       "UNAUTHORIZED",
			}
		}

		appsrvs := services.New(ctx)
		regions, err := appsrvs.Apps.Regions.List(appservice.GetAppRegions{
			AppID: input.AppID,
		})
		if err != nil {
			return &apperrors.AppError{
				StatusCode: http.StatusInternalServerError,
			}
		}

		output.Regions = regions
		return nil
	})
	u.SetTitle("List App Regions")
	u.SetTags("Regions")
	return u
}
