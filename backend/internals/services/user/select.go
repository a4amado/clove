package userservice

import (
	repository "clove/internals/services/generatedRepo"
	"fmt"

	"github.com/google/uuid"
)

func (us *Users) Get(id uuid.UUID) (*repository.User, error) {
	user, err := us.DB.User_Select(us.GetCtx(), us.ToPgUUID(id))
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}
