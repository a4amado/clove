package AppKeysHandlersV1

import (
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/services"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func DeleteAppApiKey(w http.ResponseWriter, r *http.Request) {
	apId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		http.Error(w, "Invalid App ID", http.StatusBadRequest)
		return
	}
	AppApiKey, err := uuid.Parse(r.PathValue("api_token_id"))
	if err != nil {
		http.Error(w, "Invalid Api KeyID", http.StatusBadRequest)
		return
	}

	tx, _ := postgresPool.NewTx(r.Context(), pgx.TxOptions{})
	err = services.App(r.Context(), &tx, true, apId).Key(AppApiKey).Delete()

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Failed to delete api key", http.StatusBadRequest)
		}
		tx.Rollback(r.Context())
		return
	}
	w.WriteHeader(http.StatusAccepted)
	tx.Commit(r.Context())

}
