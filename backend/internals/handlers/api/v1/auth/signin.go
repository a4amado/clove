package AuthHandlersV1

import (
	"clove/internals/apperrors"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type SignInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	ERROR_SIGNIN_INVALID_BODY      = "ERROR_SIGNIN_INVALID_BODY"
	ERROR_SIGNIN_FAILD_TO_GRAB_TX  = "ERROR_SIGNIN_FAILD_TO_GRAB_TX"
	ERROR_SIGNIN_FAILD_TO_START_TX = "ERROR_SIGNIN_FAILD_TO_START_TX"
)

func SignIn(w http.ResponseWriter, r *http.Request) {

	body := SignInBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_SIGNIN_INVALID_BODY,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

}
