package AppRegionsHandlersV1

import (
	"clove/internals/services"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func ListAppRegions(w http.ResponseWriter, r *http.Request) {
	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		http.Error(w, "Invalid App ID", http.StatusBadRequest)
		return
	}
	appSrvs := services.C(r.Context(), nil, true)
	regions, err := appSrvs.App(appId).Regions().List()
	if err != nil {
		http.Error(w, "Failed To Fetch App Regions", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(regions); err != nil {
		http.Error(w, "Failed to Send Json", http.StatusInternalServerError)
		return
	}
}
