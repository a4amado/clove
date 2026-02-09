package AppKeysHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
	"net/http"

	"github.com/google/uuid"
)

const (
	ERROR_DELETE_APP_API_KEY_INVALID_APP_ID = "ERROR_DELETE_APP_API_KEY_INVALID_APP_ID"
	ERROR_DELETE_APP_API_KEY_INVALID_KEY_ID = "ERROR_DELETE_APP_API_KEY_INVALID_KEY_ID"
	ERROR_DELETE_APP_FAILED_START_TX        = "ERROR_DELETE_APP_FAILED_START_TX"
	ERROR_DELETE_APP_FAILED_DELETE_APP      = "ERROR_DELETE_APP_FAILED_DELETE_APP"
	ERROR_DELETE_APP_FAILED_COMMIT          = "ERROR_DELETE_APP_FAILED_COMMIT"
)

func DeleteAppApiKey(w http.ResponseWriter, r *http.Request) {
	session, err := auth.ParseSessionFromRequest(r)
	if err != nil {
		auth.UnAuthResponse(w)
		return
	}
	if !session.Permessions.Can(auth.KEY, auth.DESTROY) {
		auth.UnAuthResponse(w)
		return
	}
	apId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_DELETE_APP_API_KEY_INVALID_APP_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	AppApiKey, err := uuid.Parse(r.PathValue("key_id"))
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_DELETE_APP_API_KEY_INVALID_KEY_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	srvs, tx, err := services.New(r.Context()).WithCache().WithTx()
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_DELETE_APP_FAILED_START_TX,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	err = srvs.Apps.Keys.Delete(appservice.DeleteKeyParams{
		AppID: apId,
		KeyID: AppApiKey,
	})

	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_DELETE_APP_FAILED_DELETE_APP,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
		tx.Rollback(r.Context())
		return
	}
	err = tx.Commit(r.Context())
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_DELETE_APP_FAILED_COMMIT,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
	}

}
