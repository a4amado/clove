package AuthHandlersV1

import (
	"context"

	"github.com/swaggest/usecase"
)

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInOutput struct{}

func SignIn() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input SignInInput, output *SignInOutput) error {
		return nil
	})
	u.SetTitle("User Sign In")
	u.SetTags("Auth")
	return u
}
