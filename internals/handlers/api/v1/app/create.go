package AppHandlersV1

import (
	"clove/internals/apperrors"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/services"
	repository "clove/internals/services/generatedRepo"
	"clove/internals/services/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	set "github.com/hashicorp/go-set"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateAppStruct struct {
	AppSlug        string              `json:"app_slug"`
	Regions        []repository.Region `json:"regions"`
	UserId         pgtype.UUID         `json:"user_id"`
	AllowedOrigins []string            `json:"allowed_origins"`
}

type CreateAppReponse struct {
	repository.App `json:"app"`
	Keys           []repository.AppApiKey `json:"keys"`
}

const (
	ERROR_CREATE_APP_INVALID_BODY                 = "ERROR_CREATE_APP_INVALID_BODY"
	ERROR_CREATE_APP_FAILED_START_TX              = "ERROR_CREATE_APP_FAILED_START_TX"
	ERROR_CREATE_APP_FAILED_INSERT_APP_DB         = "ERROR_CREATE_APP_FAILED_INSERT_APP_DB"
	ERROR_CREATE_APP_FAILED_GENERATE_API_KEY      = "ERROR_CREATE_APP_FAILED_GENERATE_API_KEY"
	ERROR_CREATE_API_KEY_FAILED                   = "ERROR_FAILED_CREATE_API_KEY"
	ERROR_CREATE_API_KEY_SOME_REGIONS_ARE_INVALID = "ERROR_CREATE_API_KEY_SOME_REGIONS_ARE_INVALID"
)

func CreateApp(w http.ResponseWriter, r *http.Request) {

	body := CreateAppStruct{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_INVALID_BODY,
			Message:    "",
			StatusCode: http.StatusBadRequest,

			ID: uuid.New(),
		})
		return
	}
	uniqueRegions := set.From(body.Regions)
	uniqueRegionsSlice := uniqueRegions.Slice()
	for _, region := range uniqueRegionsSlice {
		if !region.Valid() {
			apperrors.WriteError(w, &apperrors.AppError{
				ID:         uuid.New(),
				Code:       ERROR_CREATE_API_KEY_SOME_REGIONS_ARE_INVALID,
				Message:    "",
				StatusCode: http.StatusBadRequest,
			})
			return
		}
	}
	tx, err := postgresPool.NewTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_FAILED_START_TX,
			Message:    "",
			StatusCode: http.StatusInternalServerError,

			ID: uuid.New(),
		})
		return
	}
	defer tx.Rollback(r.Context()) // Always rollback on early return

	srvs := services.New(types.ServiceParams{
		Ctx:      r.Context(),
		Tx:       &tx,
		UseCache: true,
	})
	app, err := srvs.Apps().Create(repository.App_InsertParams{
		AppSlug:        fmt.Sprintf("%s:%s", uuid.NewString(), body.AppSlug),
		Regions:        body.Regions,
		AppType:        repository.AppTypePro,
		UserID:         body.UserId,
		AllowedOrigins: body.AllowedOrigins,
	})
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_FAILED_INSERT_APP_DB,
			Message:    "",
			StatusCode: http.StatusInternalServerError,

			ID: uuid.New(),
		})
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       "ERROR_FAILED_COMMIT_TX",
			Message:    "",
			StatusCode: http.StatusInternalServerError,

			ID: uuid.New(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}
