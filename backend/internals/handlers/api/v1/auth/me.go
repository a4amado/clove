package AuthHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	"clove/internals/handlers/api/httpctx"
	repository "clove/internals/services/generatedRepo"
	"context"
	"net/http"

	"github.com/swaggest/usecase"
)

type MeInput struct{}

type MeOutput struct {
	repository.User
}

func Me() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input MeInput, output *MeOutput) error {
		r := httpctx.Request(ctx)
		session, err := auth.ParseSessionFromRequest(r)
		if err != nil {
			return &apperrors.AppError{
				StatusCode: http.StatusUnauthorized,
				Code:       "UNAUTHORIZED",
			}
		}
		output.User = session.User
		return nil
	})
	u.SetTitle("Get Current User")
	u.SetTags("Auth")
	return u
}
