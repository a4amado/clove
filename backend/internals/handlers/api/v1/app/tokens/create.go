package AppTokensHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	envConsts "clove/internals/consts/env"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CreateAppOneTimeTokenBody struct {
	ChannelID string `json:"channel_id"`
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
	session, err := auth.ParseSessionFromRequest(r)
	if err != nil {
		auth.UnAuthResponse(w)
		return
	}
	if !session.Permessions.Can(auth.OneTimeToken, auth.CREATE) {
		auth.UnAuthResponse(w)
		return
	}
	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_CREATE_ONE_TIME_TOKEN_INVALID_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	body := CreateAppOneTimeTokenBody{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_CREATE_ONE_TIME_TOKEN_INVALID_BODY,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	appsrvs := services.New(r.Context())
	app, err := appsrvs.Apps.Get(appservice.GetParams{
		AppID: appId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			apperrors.WriteError(w, &apperrors.AppError{
				ID:         uuid.New(),
				Code:       ERROR_CREATE_ONE_TIME_TOKEN_APP_NOT_FOUND,
				Message:    "",
				StatusCode: http.StatusBadRequest,
			})
		} else {
			apperrors.WriteError(w, &apperrors.AppError{
				ID:         uuid.New(),
				Code:       ERROR_CREATE_ONE_TIME_TOKEN_APP_QUERY_FAILED,
				Message:    "",
				StatusCode: http.StatusBadRequest,
			})
		}
		return
	}
	if !session.Permessions.Can(auth.OneTimeToken, auth.CREATE) {
	}

	token, err := auth.GenerateOneTimeToken(*app, body.ChannelID)
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_ONE_TIME_TOKEN_FAILED_TP_GENERATE_TOKEN,
			Message:    "",
			StatusCode: http.StatusBadRequest,
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
