package userservice

import (
	repository "clove/internals/services/generatedRepo"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (us *UserService) CreateUser(params repository.InsertUserParams) (*repository.User, error) {
	us.EnsureTransactionAndContext()
	queries := us.Queries.WithTx(*us.tx)

	user, err := queries.InsertUser(us.ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, fmt.Errorf("email '%s' is already in use", params.Email)
			}
		}
		return nil, err
	}
	return &user, nil

}
