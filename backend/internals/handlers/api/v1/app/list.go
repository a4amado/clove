package AppHandlersV1

import (
	"clove/internals/middleware"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
	repository "clove/internals/services/generatedRepo"
	"context"
	"errors"

	"github.com/swaggest/usecase"
)

func ListApps() usecase.Interactor {
	Interactor := usecase.NewInteractor(func(ctx context.Context, input interface{}, output *[]repository.App) error {
		session, ok := middleware.SessionFromContext(ctx)
		if !ok {
			return errors.ErrUnsupported
		}
		srvc := services.New(ctx)
		apps, err := srvc.Apps.List(appservice.ListParams{
			UserID: session.UserID.Bytes,
		})
		if err != nil {
			return err
		}
		*output = apps
		return nil
	})
	return Interactor
}
