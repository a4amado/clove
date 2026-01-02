package AppHandlersV1

import (
	"clove/internals/apiguard"
	"clove/internals/apperrors"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/services"
	repository "clove/internals/services/generatedRepo"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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
	ERROR_INVALID_CREATE_APP_BODY    = "ERROR_INVALID_CREATE_APP_BODY"
	ERROR_FAILED_START_CREATE_APP_TX = "ERROR_FAILED_START_CREATE_APP_TX"
	ERROR_FAILED_INSERT_APP_DB       = "ERROR_FAILED_INSERT_APP_DB"
	ERROR_FAILED_GENERATE_API_KEY    = "ERROR_FAILED_GENERATE_API_KEY"
	ERROR_FAILED_CREATE_API_KEY      = "ERROR_FAILED_CREATE_API_KEY"
)

func CreateApp(w http.ResponseWriter, r *http.Request) {
	body := CreateAppStruct{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_INVALID_CREATE_APP_BODY,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			ID:         uuid.New(),
		})
		return
	}

	tx, err := postgresPool.NewTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_FAILED_START_CREATE_APP_TX,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
			Request:    r,
			ID:         uuid.New(),
		})
		return
	}
	defer tx.Rollback(r.Context()) // Always rollback on early return

	app, err := services.C(r.Context(), &tx, true).Apps().Create(repository.InsertAppParams{
		AppSlug:        fmt.Sprintf("%s:%s", uuid.NewString(), body.AppSlug),
		Regions:        body.Regions,
		AppType:        repository.AppTypePro,
		UserID:         body.UserId,
		AllowedOrigins: body.AllowedOrigins,
	})
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_FAILED_INSERT_APP_DB,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
			Request:    r,
			ID:         uuid.New(),
		})
		return
	}

	key, err := apiguard.RandomSecretKey()
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_FAILED_GENERATE_API_KEY,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
			Request:    r,
			ID:         uuid.New(),
		})
		return
	}

	appApiKey, err := services.C(r.Context(), &tx, true).App(app.App.ID.Bytes).Keys().Create(repository.CreateAppApiKeyParams{
		AppID: app.App.ID,
		Key: pgtype.Text{
			String: key,
			Valid:  true,
		},
		Name: pgtype.Text{
			String: "Clove: Initial Auto Generated Key",
		},
	})
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_FAILED_CREATE_API_KEY,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
			Request:    r,
			ID:         uuid.New(),
		})
		return
	}

	res := CreateAppReponse{
		App:  app.App,
		Keys: []repository.AppApiKey{*appApiKey},
	}

	if err := tx.Commit(r.Context()); err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       "ERROR_FAILED_COMMIT_TX",
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
			Request:    r,
			ID:         uuid.New(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
