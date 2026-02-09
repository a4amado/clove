package AuthHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"

	"clove/internals/services"
	userservice "clove/internals/services/user"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type User_Signup_Body struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

const (
	ERROR_USER_SIGNUP_BAD_BODY      = "ERROR_USER_SIGNUP_BAD_BODY"
	ERROR_USER_EMAIL_ALREADY_EXISTS = "ERROR_USER_EMAIL_ALREADY_EXISTS"
)

func User_Signup(w http.ResponseWriter, r *http.Request) {
	req_id := uuid.New()

	body := User_Signup_Body{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {

		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_USER_SIGNUP_BAD_BODY,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			ID:         req_id,
		})
		return
	}

	srvs, tx, err := services.New(r.Context()).WithTx()
	if err != nil {

		apperrors.WriteError(w, &apperrors.AppError{
			StatusCode: http.StatusInternalServerError,
			ID:         req_id,
		})
		return
	}
	code, _ := auth.GenRandKey(10)
	user, err := srvs.Users.Insert(userservice.InsertUserParams{
		Email:           body.Email,
		Password:        body.Password,
		EmailVerifycode: code,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		tx.Rollback(r.Context())

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {

			apperrors.WriteError(w, &apperrors.AppError{
				Code:       ERROR_USER_EMAIL_ALREADY_EXISTS,
				Message:    "An account with this email already exists",
				StatusCode: http.StatusConflict,
			})
			return
		}

		apperrors.WriteError(w, &apperrors.AppError{
			Message:    "",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	err = tx.Commit(r.Context())
	if err != nil {

		apperrors.WriteError(w, &apperrors.AppError{
			Message:    "",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	session_token, _ := auth.GenerateSessionToken(*user)
	r.AddCookie(&http.Cookie{
		Name:     "token",
		Value:    session_token,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})
	_ = json.NewEncoder(w).Encode(user)
}
