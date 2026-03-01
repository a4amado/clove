package AuthHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/handlers/api/httpctx"
	"clove/internals/middleware"
	authservice "clove/internals/services/auth"
	"clove/internals/services"
	"context"
	"net/http"

	"github.com/swaggest/usecase"
)

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInOutput struct{}

func SignIn() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input SignInInput, output *SignInOutput) error {
		srvs := services.New(ctx)

		result, err := srvs.Auth.Login(authservice.LoginParams{
			Email:    input.Email,
			Password: input.Password,
		})
		if err != nil {
			return &apperrors.AppError{
				StatusCode: http.StatusUnauthorized,
				Code:       "UNAUTHORIZED",
			}
		}

		if w := httpctx.ResponseWriter(ctx); w != nil {
			middleware.SetSessionCookie(w, result.Token, result.ExpiresAt)
		}

		return nil
	})
	u.SetExpectedErrors(apperrors.ErrUnauthorized)
	u.SetTitle("User Sign In")
	u.SetTags("Auth")
	return u
}
