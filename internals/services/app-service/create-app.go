package AppService

import (
	repository "clove/internals/services/generatedRepo"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (appService *AppService) CreateUser(ctx context.Context, p repository.InsertUserParams) (*repository.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	user, err := appService.Queries.InsertUser(ctx, p)
	if err != nil {
		var pgErrpr *pgconn.PgError
		if errors.As(err, &pgErrpr) {
			if pgErrpr.Code == pgerrcode.UniqueViolation {
				return nil, fmt.Errorf("email '%s' is already in use", p.Email)
			}
		}
		return nil, err
	}
	return &user, nil

}
