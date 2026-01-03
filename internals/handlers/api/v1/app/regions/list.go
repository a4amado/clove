package AppRegionsHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/services"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

const (
	ERROR_LIST_APP_REGIONS_INVALID_APP_ID  = "ERROR_LIST_APP_REGIONS_INVALID_APP_ID"
	ERROR_LIST_APP_REGIONS_FAILED_DB_QUERY = "ERROR_LIST_APP_REGIONS_FAILED_DB_QUERY"
)

func ListAppRegions(w http.ResponseWriter, r *http.Request) {
	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_LIST_APP_REGIONS_INVALID_APP_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}
	appSrvs := services.C(r.Context(), nil, true)
	regions, err := appSrvs.App(appId).Regions().List()
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_LIST_APP_REGIONS_FAILED_DB_QUERY,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}
	if err := json.NewEncoder(w).Encode(regions); err != nil {
		apperrors.WriteError(nil, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_LIST_APP_REGIONS_FAILED_DB_QUERY,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}
}
