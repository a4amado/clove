package AuthHandlersV1

import (
	"clove/internals/apperrors"
	postgresPool "clove/internals/data/postgres/pool"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	tx, err := postgresPool.NewTx(r.Context(), pgx.TxOptions{})

	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_SIGNIN_FAILD_TO_GRAB_TX,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	tx, err = tx.Begin(r.Context())
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_SIGNIN_FAILD_TO_START_TX,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

}
