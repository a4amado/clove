package userservice

import (
	"clove/internals/services/types"

	"github.com/google/uuid"
)

type User struct {
	*types.BaseService
	UserID uuid.UUID
}

func (u *User) SetUserId(id uuid.UUID) {
	u.UserID = id
}

type Users struct {
	*types.BaseService
}

type Email struct {
	*User
	EmailID uuid.UUID
}

type Emails struct {
	*User
}
