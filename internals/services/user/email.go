package userservice

import (
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (e *Email) Delete() error {
	return e.DB.UserEmail_Delete(e.GetCtx(), e.ToPgUUID(e.EmailID))
}

func (e *Emails) List() (*[]repository.UserEmail, error) {
	email, err := e.DB.UserEmails_List_ByUserID(e.GetCtx(), e.ToPgUUID(e.UserID))
	if err != nil {
		return nil, err
	}
	return &email, nil
}
func (e *Email) GetById() (repository.UserEmail, error) {
	return e.DB.UserEmail_Select(e.GetCtx(), e.ToPgUUID(e.EmailID))
}

func (e *Emails) GetByEmail(email string) (repository.UserEmail, error) {
	return e.DB.UserEmail_SelectByEmail(e.GetCtx(), email)
}

func (e *Email) Verify() error {
	return e.DB.UserEmails_Verify(e.GetCtx(), pgtype.UUID{
		Bytes: e.EmailID,
		Valid: true,
	})
}
