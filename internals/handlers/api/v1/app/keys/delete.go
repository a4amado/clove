package AppKeysHandlersV1

import (
	"clove/internals/apperrors"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/services"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	ERROR_DELETE_APP_API_KEY_INVALID_APP_ID = "ERROR_DELETE_APP_API_KEY_INVALID_APP_ID"
	ERROR_DELETE_APP_API_KEY_INVALID_KEY_ID = "ERROR_DELETE_APP_API_KEY_INVALID_KEY_ID"
	ERROR_DELETE_APP_FAILED_START_TX        = "ERROR_DELETE_APP_FAILED_START_TX"
	ERROR_DELETE_APP_FAILED_DELETE_APP      = "ERROR_DELETE_APP_FAILED_DELETE_APP"
	ERROR_DELETE_APP_FAILED_COMMIT          = "ERROR_DELETE_APP_FAILED_COMMIT"
)

func DeleteAppApiKey(w http.ResponseWriter, r *http.Request) {
	apId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_DELETE_APP_API_KEY_INVALID_APP_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}
	AppApiKey, err := uuid.Parse(r.PathValue("key_id"))
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_DELETE_APP_API_KEY_INVALID_KEY_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}

	tx, err := postgresPool.NewTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_DELETE_APP_FAILED_START_TX,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}
	tx.Begin(r.Context())
	n, err := services.C(r.Context(), &tx, true).App(apId).Key(AppApiKey).Delete()

	if n == 0 || err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_DELETE_APP_FAILED_DELETE_APP,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		tx.Rollback(r.Context())
		return
	}
	err = tx.Commit(r.Context())
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_DELETE_APP_FAILED_COMMIT,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
	}

}
