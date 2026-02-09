package AppKeysHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

	session, err := auth.ParseSessionFromRequest(r)
	if err != nil {
		auth.UnAuthResponse(w)
		return
	}

	if !session.Permessions.Can(auth.KEY, auth.CREATE) {
		auth.UnAuthResponse(w)
		return
	}

	apId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_API_KEY_INVALID_APP_ID,
			Message:    fmt.Sprintf("%s is not a valid id", apId.String()),
			StatusCode: http.StatusBadRequest,

			ID: uuid.New(),
		})
		return
	}

	body := CreateAppApiTokenBody{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_API_KEY_INVALID_BODY,
			Message:    "",
			StatusCode: http.StatusBadRequest,

			ID: uuid.New(),
		})
		return
	}

	token, err := auth.GenerateSDKToken(apId)
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_API_KEY_FAILED_GENERATE_KEY,
			Message:    "",
			StatusCode: http.StatusInternalServerError,

			ID: uuid.New(),
		})
		return
	}
	srvs, tx, err := services.New(r.Context()).WithCache().WithTx()
	key, err := srvs.Apps.Keys.Create(appservice.CreateKeyParams{
		AppID:   apId,
		KeyName: body.Name,
		Prefix:  token[:min(5, len(token)-1)],
		Suffix:  token[max(0, len(token)-5):],
	})
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_API_KEY_FAILED_CREATE,
			Message:    "",
			StatusCode: http.StatusInternalServerError,

			ID: uuid.New(),
		})
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_CREATE_APP_API_KEY_FAILED_COMMIT,
			Message:    "",
			StatusCode: http.StatusInternalServerError,

			ID: uuid.New(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(key); err != nil {
		// Can't write error after headers sent, just log it
		// In production, use proper logging here
		return
	}
}
