package AuthHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	"clove/internals/handlers/api/httpctx"
	"clove/internals/services"
	userservice "clove/internals/services/user"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/swaggest/usecase"

	repository "clove/internals/services/generatedRepo"
)

type SignupInput struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type SignupOutput struct {
	repository.User
}

const (
	ERROR_USER_SIGNUP_BAD_BODY      = "ERROR_USER_SIGNUP_BAD_BODY"
	ERROR_USER_EMAIL_ALREADY_EXISTS = "ERROR_USER_EMAIL_ALREADY_EXISTS"
)

func Signup() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input SignupInput, output *SignupOutput) error {
		srvs, tx, err := services.New(ctx).WithTx()
		if err != nil {
			return &apperrors.AppError{
				StatusCode: http.StatusInternalServerError,
			}
		}
		code, _ := auth.GenRandKey(10)
		user, err := srvs.Users.Insert(userservice.InsertUserParams{
			Email:           input.Email,
			Password:        input.Password,
			EmailVerifycode: code,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			tx.Rollback(ctx)

			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				return &apperrors.AppError{
					Code:       ERROR_USER_EMAIL_ALREADY_EXISTS,
					Message:    "An account with this email already exists",
					StatusCode: http.StatusConflict,
				}
			}

			return &apperrors.AppError{
				StatusCode: http.StatusInternalServerError,
			}
		}

		if err = tx.Commit(ctx); err != nil {
			return &apperrors.AppError{
				StatusCode: http.StatusInternalServerError,
			}
		}

		sessionToken, _ := auth.GenerateSessionToken(*user)
		if r := httpctx.Request(ctx); r != nil {
			http.SetCookie(httpctx.ResponseWriter(ctx), &http.Cookie{
				Name:     "token",
				Value:    sessionToken,
				Expires:  time.Now().Add(time.Hour * 24 * 7),
				HttpOnly: true,
				SameSite: http.SameSiteDefaultMode,
			})
		}

		output.User = *user
		return nil
	})
	u.SetTitle("User Signup")
	u.SetTags("Auth")
	return u
}
