package AuthHandlersV1

import (
	"clove/internals/apiguard"
	"clove/internals/apperrors"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/logger"
	"clove/internals/services"
	"clove/internals/services/types"
	userservice "clove/internals/services/user"
	"clove/internals/tokenguard"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
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
	log := logger.NewJSONLogger(logger.JSONLoggerParams{
		RequestID: req_id,
		Req:       r,
	})

	body := User_Signup_Body{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("failed to decode request body",
			slog.String("error", err.Error()),
		)
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_USER_SIGNUP_BAD_BODY,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			ID:         req_id,
		})
		return
	}

	tx, err := postgresPool.NewTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		log.Error("failed to start database transaction",
			slog.String("error", err.Error()),
		)
		apperrors.WriteError(w, &apperrors.AppError{
			StatusCode: http.StatusInternalServerError,
			ID:         req_id,
		})
		return
	}
	defer tx.Rollback(r.Context())

	srvs := services.New(types.ServiceParams{
		Ctx:      r.Context(),
		Tx:       &tx,
		UseCache: false,
	})

	code, _ := apiguard.RandomSecretKey()
	user, err := srvs.Users().Insert(userservice.InsertUserParams{
		Email:           body.Email,
		Password:        body.Password,
		EmailVerifycode: code,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		tx.Rollback(r.Context())

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.Warn("user signup failed: email already exists",
				slog.String("email", body.Email),
				slog.String("error", pgErr.Error()),
			)
			apperrors.WriteError(w, &apperrors.AppError{
				Code:       ERROR_USER_EMAIL_ALREADY_EXISTS,
				Message:    "An account with this email already exists",
				StatusCode: http.StatusConflict,
			})
			return
		}

		log.Error("failed to insert user",
			slog.String("error", err.Error()),
			slog.String("email", body.Email),
		)
		apperrors.WriteError(w, &apperrors.AppError{
			Message:    "",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	err = tx.Commit(r.Context())
	if err != nil {
		log.Error("failed to commit transaction",
			slog.String("error", err.Error()),
			slog.String("user_id", user.ID.String()),
		)
		apperrors.WriteError(w, &apperrors.AppError{
			Message:    "",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	log.Info("user signup successful",
		slog.String("user_id", user.ID.String()),
		slog.String("email", body.Email),
	)

	session_token, _ := tokenguard.GenerateSessionToken(*user)
	r.AddCookie(&http.Cookie{
		Name:     "token",
		Value:    session_token,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})
	_ = json.NewEncoder(w).Encode(user)
}
