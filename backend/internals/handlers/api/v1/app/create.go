package AppHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	"clove/internals/handlers/api/httpctx"
	"clove/internals/services"
	repository "clove/internals/services/generatedRepo"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	set "github.com/hashicorp/go-set"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/swaggest/usecase"
)

type CreateAppInput struct {
	AppSlug        string              `json:"app_slug"`
	Regions        []repository.Region `json:"regions"`
	UserId         uuid.UUID           `json:"user_id"`
	AllowedOrigins []string            `json:"allowed_origins"`
	Plan           repository.AppType  `json:"plan" enum:"free,pro,enterprise"`
}

type CreateAppOutput struct {
	repository.App
}

const (
	ERROR_CREATE_APP_INVALID_BODY                 = "ERROR_CREATE_APP_INVALID_BODY"
	ERROR_CREATE_APP_FAILED_START_TX              = "ERROR_CREATE_APP_FAILED_START_TX"
	ERROR_CREATE_APP_FAILED_INSERT_APP_DB         = "ERROR_CREATE_APP_FAILED_INSERT_APP_DB"
	ERROR_CREATE_APP_FAILED_GENERATE_API_KEY      = "ERROR_CREATE_APP_FAILED_GENERATE_API_KEY"
	ERROR_CREATE_API_KEY_FAILED                   = "ERROR_FAILED_CREATE_API_KEY"
	ERROR_CREATE_API_KEY_SOME_REGIONS_ARE_INVALID = "ERROR_CREATE_API_KEY_SOME_REGIONS_ARE_INVALID"
)

func CreateApp() usecase.Interactor {

	u := usecase.NewInteractor(func(ctx context.Context, input CreateAppInput, output *CreateAppOutput) error {
		r := httpctx.Request(ctx)
		session, err := auth.ParseSessionFromRequest(r)
		if err != nil || !session.Permessions.Can(auth.APP, auth.CREATE) {
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
					Code:       ERROR_CREATE_API_KEY_SOME_REGIONS_ARE_INVALID,
					StatusCode: http.StatusBadRequest,
				}
			}
		}

		srvs, tx, err := services.New(ctx).WithTx()
		if err != nil {
			return &apperrors.AppError{
				Code:       ERROR_CREATE_APP_FAILED_START_TX,
				StatusCode: http.StatusInternalServerError,
			}
		}
		app, err := srvs.Apps.Create(repository.App_InsertParams{
			AppSlug: fmt.Sprintf("%s:%s", uuid.NewString(), input.AppSlug),
			Regions: input.Regions,
			AppType: repository.AppTypePro,
			UserID: pgtype.UUID{
				Bytes: input.UserId,
				Valid: true,
			},
			AllowedOrigins: input.AllowedOrigins,
		})
		if err != nil {
			return &apperrors.AppError{
				Code:       ERROR_CREATE_APP_FAILED_INSERT_APP_DB,
				StatusCode: http.StatusInternalServerError,
			}
		}

		if err := tx.Commit(ctx); err != nil {
			return &apperrors.AppError{
				Code:       "ERROR_FAILED_COMMIT_TX",
				StatusCode: http.StatusInternalServerError,
			}
		}

		output.App = *app
		return nil
	})
	u.SetTitle("Create App")
	u.SetTags("Apps")
	return u
}
