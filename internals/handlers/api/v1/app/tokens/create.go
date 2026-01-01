package AppTokensHandlersV1

import (
	"clove/internals/apiguard"
	envConsts "clove/internals/consts/env"
	"clove/internals/services"
	"clove/internals/tokenguard"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CreateAppOneTimeTokenBody struct {
	ChannelID string    `json:"channel_id"`
	ApiKeyId  uuid.UUID `json:"api_key_id"`
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
	app, err := services.C(r.Context(), nil, true).App(appId).Get()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
	key, err := services.C(r.Context(), nil, true).App(appId).Key(body.ApiKeyId).Get()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
	if key.String != apiguard.GetHeaderApi(r) {
		http.Error(w, "Incorrect Api Key", http.StatusForbidden)
		return
	}
	token, err := tokenguard.GenerateOneTimeToken(*app, body.ChannelID, body.ApiKeyId)
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
