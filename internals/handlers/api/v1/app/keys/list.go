package AppKeysHandlersV1

import (
	"clove/internals/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ListAppApiKeys(w http.ResponseWriter, r *http.Request) {
	appId, err := uuid.Parse(r.PathValue("app_id"))

	if err != nil || appId == uuid.Nil {
		http.Error(w, "Invalid App ID", http.StatusBadRequest)
		return
	}
	pageIdxStr := r.URL.Query().Get("page_idx")
	if pageIdxStr == "" {
		pageIdxStr = "0"
	}
	page_idx, err := strconv.ParseInt(pageIdxStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid page_idx search params '?page_idx=int'", http.StatusBadRequest)
		return
	}
	keys, err := services.C(r.Context(), nil, true).App(appId).Keys().List(int32(page_idx))
	for idx, key := range keys {
		key.Key = pgtype.Text{
			String: "[Redacted]",
			Valid:  true,
		}
		keys[idx] = key

	}
	if err != nil {
		http.Error(w, "Insternal server error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(keys)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
