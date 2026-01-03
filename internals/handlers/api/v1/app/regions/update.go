package AppRegionsHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/services"
	repository "clove/internals/services/generatedRepo"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/go-set"
)

type UpdateAppRegionsBody struct {
	Regions []repository.Region `json:"regions"`
}

const (
	ERROR_UPDATE_APP_REGIONS_INVALID_ID               = "clove.io/app/region/update/invalid.app.id"
	ERROR_UPDATE_APP_REGIONS_INVALID_BODY             = "clove.io/app/region/update/invalid.req.body"
	ERROR_UPDATE_APP_REGIONS_SOME_REGIONS_ARE_INVALID = "clove.io/app/region/update/invalid.one.or.more.regions"
)

func UpdateAppRegions(w http.ResponseWriter, r *http.Request) {
	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_UPDATE_APP_REGIONS_INVALID_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}
	body := UpdateAppRegionsBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(&w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_UPDATE_APP_REGIONS_INVALID_BODY,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
		})
		return
	}
	uniqueRegions := set.From(body.Regions)
	uniqueRegionsSlice := uniqueRegions.Slice()
	for _, region := range uniqueRegionsSlice {
		if !region.Valid() {
			apperrors.WriteError(&w, &apperrors.AppError{
				ID:         uuid.New(),
				Code:       ERROR_UPDATE_APP_REGIONS_SOME_REGIONS_ARE_INVALID,
				Message:    "",
				StatusCode: http.StatusBadRequest,
				Internal:   err,
				Request:    r,
			})
			return
		}
	}
	if len(uniqueRegionsSlice) != len(body.Regions) {
		http.Error(w, "Regions array shal not contain any doublicates", http.StatusBadRequest)
		return
	}

	err = services.C(r.Context(), nil, true).App(appId).Regions().Update(body.Regions)
	if err != nil {
		http.Error(w, "Failed To Update App Regions", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(body.Regions); err != nil {
		http.Error(w, "Failed to Send Json", http.StatusInternalServerError)
		return
	}
}
