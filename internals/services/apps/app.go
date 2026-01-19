package appservice

import (
	repository "clove/internals/services/generatedRepo"
	"fmt"

	"github.com/google/uuid"
)

func (s *AppsService) Create(arg repository.App_InsertParams) (*repository.App, error) {
	app, err := s.DB.App_Insert(s.GetCtx(), arg)
	if err != nil {
		return nil, fmt.Errorf("failed to insert app: %w", err)
	}
	return &app, nil
}

type GetParams struct {
	AppID uuid.UUID
}

func (s *AppsService) Get(args GetParams) (*repository.App, error) {

	app, err := s.DB.App_Select(s.GetCtx(), s.ToPgUUID(args.AppID))
	if err != nil {
		return nil, fmt.Errorf("failed to insert app: %w", err)
	}
	return &app, nil
}

func (s *AppsService) Update() {

}
