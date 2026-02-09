package AppRegionsHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
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
	PublicErrUpdateAppRegionsInvalidId   = "clove.io/app/region/update/invalid.app.id"
	PublicErrUpdateAppRegionsInvalidBody = "clove.io/app/region/update/invalid.req.body"
	PublicErrUpdateAppRegionsAreInvalid  = "clove.io/app/region/update/invalid.one.or.more.regions"
)

func UpdateAppRegions(w http.ResponseWriter, r *http.Request) {
	session, err := auth.ParseSessionFromRequest(r)
	if err != nil {
		auth.UnAuthResponse(w)
		return
	}
	if !session.Permessions.Can(auth.KEY, auth.UPDATE) {
		auth.UnAuthResponse(w)
		return
	}
	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       PublicErrUpdateAppRegionsInvalidId,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	body := UpdateAppRegionsBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       PublicErrUpdateAppRegionsInvalidBody,
			Message:    "",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	uniqueRegions := set.From(body.Regions)
	uniqueRegionsSlice := uniqueRegions.Slice()
	for _, region := range uniqueRegionsSlice {
		if !region.Valid() {
			apperrors.WriteError(w, &apperrors.AppError{
				ID:         uuid.New(),
				Code:       PublicErrUpdateAppRegionsAreInvalid,
				Message:    "",
				StatusCode: http.StatusBadRequest,
			})
			return
		}
	}
	if len(uniqueRegionsSlice) != len(body.Regions) {
		apperrors.WriteError(w, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       PublicErrUpdateAppRegionsAreInvalid,
			Message:    "Regions array shal not contain any doublicates",
			StatusCode: http.StatusBadRequest,
		})

		return
	}
	appsrvs := services.New(r.Context())
	err = appsrvs.Apps.Regions.Update(appservice.UpdateRegions{
		Regions: uniqueRegionsSlice,
		AppID:   appId,
	})
	if err != nil {
		http.Error(w, "Failed To Update App Regions", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(body.Regions); err != nil {
		http.Error(w, "Failed to Send Json", http.StatusInternalServerError)
		return
	}
}
