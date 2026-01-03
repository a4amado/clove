package AppKeysHandlersV1

import (
	"clove/internals/apiguard"
	"clove/internals/apperrors"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/services"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type create_app_key_error string

const (
	ERROR_CREATE_APP_API_KEY_INVALID_APP_ID      = "ERROR_CREATE_APP_API_KEY_INVALID_APP_ID"
	ERROR_CREATE_APP_API_KEY_FAILED_START_TX     = "ERROR_CREATE_APP_API_KEY_FAILED_START_TX"
	ERROR_CREATE_APP_API_KEY_INVALID_BODY        = "ERROR_CREATE_APP_API_KEY_INVALID_BODY"
	ERROR_CREATE_APP_API_KEY_FAILED_GENERATE_KEY = "ERROR_CREATE_APP_API_KEY_FAILED_GENERATE_KEY"
	ERROR_CREATE_APP_API_KEY_FAILED_CREATE       = "ERROR_CREATE_APP_API_KEY_FAILED_CREATE"
	ERROR_CREATE_APP_API_KEY_FAILED_ENCODE       = "ERROR_CREATE_APP_API_KEY_FAILED_ENCODE"
	ERROR_CREATE_APP_API_KEY_FAILED_COMMIT       = "ERROR_CREATE_APP_API_KEY_FAILED_COMMIT"
)

type CreateAppApiTokenBody struct {
	Name string `json:"name"`
}

func CreateAppApiKey(w http.ResponseWriter, r *http.Request) {
	apId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_API_KEY_INVALID_APP_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	tx, err := postgresPool.NewTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_API_KEY_FAILED_START_TX,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}
	defer tx.Rollback(r.Context())

	body := CreateAppApiTokenBody{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_API_KEY_INVALID_BODY,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	randomKey, err := apiguard.RandomSecretKey()
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_API_KEY_FAILED_GENERATE_KEY,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	api, err := services.C(r.Context(), &tx, true).App(apId).Keys().Create(body.Name, randomKey)
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_API_KEY_FAILED_CREATE,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_API_KEY_FAILED_COMMIT,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(api); err != nil {
		// Can't write error after headers sent, just log it
		// In production, use proper logging here
		return
	}
}
