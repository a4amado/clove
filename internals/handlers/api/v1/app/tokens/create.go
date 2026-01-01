package AppTokensHandlersV1

import (
	envConsts "clove/internals/consts/env"
	"clove/internals/services"
	"clove/internals/tokenguard"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type CreateAppOneTimeTokenBody struct {
	ChannelID string
}

func CreateAppOneTimeToken(w http.ResponseWriter, r *http.Request) {
	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		http.Error(w, "Invalid App ID", http.StatusUnauthorized)
		return
	}
	body := CreateAppOneTimeTokenBody{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}
	app, err := services.App(r.Context(), nil, true, appId).Get()
	if err != nil {
		http.Error(w, "App Not Found", http.StatusNotFound)
		return
	}
	token, err := tokenguard.GenerateOneTimeToken(*app, body.ChannelID)
	if err != nil {
		http.Error(w, "Failed To Create Token"+err.Error(), http.StatusInternalServerError)
		return
	}

	res := map[string]any{
		"token":  token,
		"region": envConsts.Region(),
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
