package userservice

import (
	repository "clove/internals/services/generatedRepo"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type DeleteEmailParams struct {
	EmailID uuid.UUID
}

func (e *Emails) Delete(args DeleteEmailParams) error {
	return e.DB.UserEmail_Delete(e.GetCtx(), e.ToPgUUID(args.EmailID))
}

type ListUserEmailParams struct {
	UserID uuid.UUID
}

func (e *Emails) List(args ListUserEmailParams) (*[]repository.UserEmail, error) {
	email, err := e.DB.UserEmails_List_ByUserID(e.GetCtx(), e.ToPgUUID(args.UserID))
	if err != nil {
		return nil, err
	}
	return &email, nil
}

type GetEmailParams struct {
	EmailID uuid.UUID
}

func (e *Emails) Get(args GetEmailParams) (repository.UserEmail, error) {
	return e.DB.UserEmail_Select(e.GetCtx(), e.ToPgUUID(args.EmailID))
}

type GetEmailByEmailParams struct {
	Email string
}

func (e *Emails) GetByEmail(args GetEmailByEmailParams) (repository.UserEmail, error) {
	return e.DB.UserEmail_SelectByEmail(e.GetCtx(), args.Email)
}

type VerifyEmail struct {
	EmailID uuid.UUID
}

func (e *Emails) Verify(args VerifyEmail) error {
	return e.DB.UserEmails_Verify(e.GetCtx(), pgtype.UUID{
		Bytes: args.EmailID,
		Valid: true,
	})
}
