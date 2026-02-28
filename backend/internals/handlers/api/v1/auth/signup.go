package AuthHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	"clove/internals/handlers/api/httpctx"
	"clove/internals/services"
	credentialservice "clove/internals/services/credential"
	repository "clove/internals/services/generatedRepo"
	userservice "clove/internals/services/user"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/swaggest/usecase"
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
			tx.Rollback(ctx)
			var pgErr *pgconn.PgError
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

		token, err := auth.GenRandKey(32)
		if err != nil {
			tx.Rollback(ctx)
			return &apperrors.AppError{StatusCode: http.StatusInternalServerError}
		}

		perms := auth.NewPermissionsBuilder()
		perms.Allow(auth.RESOURCES_ALL, auth.OPERATIONS_ALL)
		permsStr, _ := perms.String()

		expires := time.Now().Add(time.Hour * 24 * 7)
		_, err = srvs.Credentials.Create(credentialservice.InsertParams{
			Token:       token,
			UserID:      uuid.UUID(user.ID.Bytes),
			Type:        repository.CredentialTypeSession,
			Permissions: permsStr,
			ExpiresAt:   expires,
		})
		if err != nil {
			tx.Rollback(ctx)
			return &apperrors.AppError{StatusCode: http.StatusInternalServerError}
		}

		if err = tx.Commit(ctx); err != nil {
			return &apperrors.AppError{StatusCode: http.StatusInternalServerError}
		}

		if w := httpctx.ResponseWriter(ctx); w != nil {
			auth.SetSessionCookie(w, token, expires)
		}

		output.User = *user
		return nil
	})
	u.SetExpectedErrors(apperrors.ErrConflict, apperrors.ErrInternalServerError)
	u.SetTitle("User Signup")
	u.SetTags("Auth")
	return u
}
