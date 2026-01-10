package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"fmt"
)

func (s *AppsService) Create(arg repository.App_InsertParams) (*repository.App, error) {
	app, err := s.Q().App_Insert(s.CTX(), arg)
	if err != nil {
		return nil, fmt.Errorf("failed to insert app: %w", err)
	}
	return &app, nil
}
