package userservice

import (
	"clove/internals/services/types"
)

type Users struct {
	*types.BaseService
}

type Emails struct {
	*Users
}
