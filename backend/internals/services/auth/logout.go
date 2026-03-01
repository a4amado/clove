package authservice

import "github.com/google/uuid"

// Logout deletes the credential associated with the given token.
// It is a no-op if the token does not exist.
func (s *AuthService) Logout(token string) error {
	cred, err := s.GetByToken(s.GetCtx(), token)
	if err != nil {
		return nil
	}
	return s.Delete(uuid.UUID(cred.ID.Bytes))
}
