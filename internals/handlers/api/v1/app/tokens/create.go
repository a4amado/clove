package AppTokensHandlersV1

import (
	"clove/internals/apiguard"
	"clove/internals/apperrors"
	envConsts "clove/internals/consts/env"
	"clove/internals/services"
	"clove/internals/tokenguard"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CreateAppOneTimeTokenBody struct {
	ChannelID string    `json:"channel_id"`
	ApiKeyId  uuid.UUID `json:"api_key_id"`
}

const (
	ERROR_CREATE_ONE_TIME_TOKEN_INVALID_ID        = "ERROR_CREATE_ONE_TIME_TOKEN_INVALID_ID"
	ERROR_CREATE_ONE_TIME_TOKEN_INVALID_BODY      = "ERROR_CREATE_ONE_TIME_TOKEN_INVALID_BODY"
	ERROR_CREATE_ONE_TIME_TOKEN_APP_NOT_FOUND     = "ERROR_CREATE_ONE_TIME_TOKEN_APP_NOT_FOUND"
	ERROR_CREATE_ONE_TIME_TOKEN_APP_QUERY_FAILED  = "ERROR_CREATE_ONE_TIME_TOKEN_APP_QUERY_FAILED"
	ERROR_CREATE_ONE_TIME_TOKEN_KEY_NOT_FOUND     = "ERROR_CREATE_ONE_TIME_TOKEN_APP_NOT_FOUND"
	ERROR_CREATE_ONE_TIME_TOKEN_KEY_QUERY_FAILED  = "ERROR_CREATE_ONE_TIME_TOKEN_APP_QUERY_FAILED"
	ERROR_ONE_TIME_TOKEN_KEY_ID_MISMATCH          = "ERROR_ONE_TIME_TOKEN_KEY_ID_MISMATCH"
	ERROR_ONE_TIME_TOKEN_FAILED_TP_GENERATE_TOKEN = "ERROR_ONE_TIME_TOKEN_FAILED_TP_GENERATE_TOKEN"
)

func CreateAppOneTimeToken(w http.ResponseWriter, r *http.Request) {
	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_CREATE_ONE_TIME_TOKEN_INVALID_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}
	body := CreateAppOneTimeTokenBody{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_CREATE_ONE_TIME_TOKEN_INVALID_BODY,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}
	app, err := services.C(r.Context(), nil, true).App(appId).Get()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			apperrors.WriteError(&w, &apperrors.AppError{
				ID:         uuid.New(),
				Code:       ERROR_CREATE_ONE_TIME_TOKEN_APP_NOT_FOUND,
				Message:    "",
				StatusCode: http.StatusBadRequest,
				Internal:   err,
				Request:    r,
			})
		} else {
			apperrors.WriteError(&w, &apperrors.AppError{
				ID:         uuid.New(),
				Code:       ERROR_CREATE_ONE_TIME_TOKEN_APP_QUERY_FAILED,
				Message:    "",
				StatusCode: http.StatusBadRequest,
				Internal:   err,
				Request:    r,
			})
		}
		return
	}

	key, err := services.C(r.Context(), nil, true).App(appId).Key(body.ApiKeyId).Get()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			apperrors.WriteError(&w, &apperrors.AppError{
				ID:         uuid.New(),
				Code:       ERROR_CREATE_ONE_TIME_TOKEN_KEY_NOT_FOUND,
				Message:    "",
				StatusCode: http.StatusBadRequest,
				Internal:   err,
				Request:    r,
			})
		} else {
			apperrors.WriteError(&w, &apperrors.AppError{
				ID:         uuid.New(),
				Code:       ERROR_CREATE_ONE_TIME_TOKEN_KEY_QUERY_FAILED,
				Message:    "",
				StatusCode: http.StatusBadRequest,
				Internal:   err,
				Request:    r,
			})
		}
		return
	}
	if key.String != apiguard.GetHeaderApi(r) {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_ONE_TIME_TOKEN_KEY_ID_MISMATCH,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}
	token, err := tokenguard.GenerateOneTimeToken(*app, body.ChannelID, body.ApiKeyId)
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_ONE_TIME_TOKEN_FAILED_TP_GENERATE_TOKEN,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}

	res := map[string]any{
		"token":  token,
		"region": envConsts.Region(),
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return
	}

}
