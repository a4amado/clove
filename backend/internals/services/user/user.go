package userservice

import (
	repository "clove/internals/services/generatedRepo"
)

type InsertUserParams struct {
	Email           string
	Password        string
	EmailVerifycode string
}

func (us *Users) Insert(args InsertUserParams) (*repository.User, error) {

	user, err := us.DB.User_Insert(us.GetCtx(), args.Password)
	if err != nil {
		return nil, err
	}

	_, err = us.DB.UserEmail_Insert(us.GetCtx(), repository.UserEmail_InsertParams{
		Email:  args.Email,
		UserID: user.ID,
		Code:   args.EmailVerifycode,
	})
	if err != nil {

		return nil, err
	}

	return &user, nil
}
