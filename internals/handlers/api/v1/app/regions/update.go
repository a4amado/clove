package AppRegionsHandlersV1

import (
	"clove/internals/services"
	repository "clove/internals/services/generatedRepo"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/go-set"
)

type UpdateAppRegionsBody struct {
	Regions []repository.Region `json:"regions"`
}

func UpdateAppRegions(w http.ResponseWriter, r *http.Request) {
	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		http.Error(w, "Invalid App ID", http.StatusBadRequest)
		return
	}
	body := UpdateAppRegionsBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	uniqueRegions := set.From(body.Regions)
	uniqueRegionsSlice := uniqueRegions.Slice()
	for _, region := range uniqueRegionsSlice {
		if !region.Valid() {
			http.Error(w, fmt.Sprintf("%v is not a valid region", region), http.StatusBadRequest)
			return
		}
	}
	if len(uniqueRegionsSlice) != len(body.Regions) {
		http.Error(w, "Regions array shal not contain any doublicates", http.StatusBadRequest)
		return
	}

	err = services.App(r.Context(), nil, true, appId).Regions(&body.Regions).Update()
	if err != nil {
		http.Error(w, "Failed To Update App Regions", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(body.Regions); err != nil {
		http.Error(w, "Failed to Send Json", http.StatusInternalServerError)
		return
	}
}
