package userservice

import (
	repository "clove/internals/services/generatedRepo"

	"github.com/jackc/pgx/v5/pgtype"
)

func (e *Email) Delete() error {
	return e.Q().UserEmail_Delete(e.CTX(), e.ToPgUUID(e.EmailID))
}

func (e *Emails) List() (*[]repository.UserEmail, error) {
	email, err := e.Q().UserEmails_List_ByUserID(e.CTX(), e.ToPgUUID(e.UserID))
	if err != nil {
		return nil, err
	}
	return &email, nil
}
func (e *Email) GetById() (repository.UserEmail, error) {
	return e.Q().UserEmail_Select(e.CTX(), e.ToPgUUID(e.EmailID))
}

func (e *Emails) GetByEmail(email string) (repository.UserEmail, error) {
	return e.Q().UserEmail_SelectByEmail(e.CTX(), email)
}

func (e *Email) Verify() error {
	return e.Q().UserEmails_Verify(e.CTX(), pgtype.UUID{
		Bytes: e.EmailID,
		Valid: true,
	})
}
