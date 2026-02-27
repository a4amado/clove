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
	"github.com/hashicorp/go-set"
	"github.com/swaggest/usecase"
)

type UpdateRegionsInput struct {
	AppID   uuid.UUID           `path:"app_id"`
	Regions []repository.Region `json:"regions"`
}

type UpdateRegionsOutput struct {
	Regions []repository.Region `json:"regions"`
}

func UpdateAppRegions() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input UpdateRegionsInput, output *UpdateRegionsOutput) error {
		r := httpctx.Request(ctx)
		session, err := middleware.ParseSession(r, postgresPool.Client())
		if err != nil || !session.Permissions.Can(middleware.KEY, middleware.UPDATE) {
			return &apperrors.AppError{
				StatusCode: http.StatusUnauthorized,
				Code:       "UNAUTHORIZED",
			}
		}

		uniqueRegions := set.From(input.Regions)
		uniqueRegionsSlice := uniqueRegions.Slice()
		for _, region := range uniqueRegionsSlice {
			if !region.Valid() {
				return &apperrors.AppError{
					Code:       "clove.io/app/region/update/invalid.one.or.more.regions",
					StatusCode: http.StatusBadRequest,
				}
			}
		}
		if len(uniqueRegionsSlice) != len(input.Regions) {
			return &apperrors.AppError{
				Code:       "clove.io/app/region/update/invalid.one.or.more.regions",
				Message:    "Regions array shall not contain any duplicates",
				StatusCode: http.StatusBadRequest,
			}
		}

		appsrvs := services.New(ctx)
		err = appsrvs.Apps.Regions.Update(appservice.UpdateRegions{
			Regions: uniqueRegionsSlice,
			AppID:   input.AppID,
		})
		if err != nil {
			return &apperrors.AppError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update app regions",
			}
		}

		output.Regions = input.Regions
		return nil
	})
	u.SetTitle("Update App Regions")
	u.SetTags("Regions")
	return u
}
