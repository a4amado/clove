package AuthHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/middleware"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/handlers/api/httpctx"
	"clove/internals/services"
	repository "clove/internals/services/generatedRepo"
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/swaggest/usecase"
)

type MeInput struct{}

type MeOutput struct {
	repository.User
}

func Me() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input MeInput, output *MeOutput) error {
		r := httpctx.Request(ctx)
		session, err := middleware.ParseSession(r, postgresPool.Client())
		if err != nil || session.SessionType != middleware.RegularSession {
			return &apperrors.AppError{
				StatusCode: http.StatusUnauthorized,
				Code:       "UNAUTHORIZED",
			}
		}

		srvs := services.New(ctx)
		user, err := srvs.Users.Get(uuid.UUID(session.UserID.Bytes))
		if err != nil {
			return &apperrors.AppError{
				StatusCode: http.StatusUnauthorized,
				Code:       "UNAUTHORIZED",
			}
		}

		output.User = *user
		return nil
	})
	u.SetTitle("Get Current User")
	u.SetTags("Auth")
	return u
}
