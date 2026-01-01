package AppKeysHandlersV1

import (
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/services"
	"fmt"
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
	AppApiKey, err := uuid.Parse(r.PathValue("key_id"))
	if err != nil {
		http.Error(w, "Invalid Api KeyID", http.StatusBadRequest)
		return
	}

	tx, _ := postgresPool.NewTx(r.Context(), pgx.TxOptions{})
	tx.Begin(r.Context())
	n, err := services.App(r.Context(), &tx, true, apId).Key(AppApiKey).Delete()

	if err != nil {
		http.Error(w, "Failed to delete api key", http.StatusBadRequest)
		tx.Rollback(r.Context())
		return
	}
	if n == 0 {
		http.Error(w, fmt.Sprintf("Failed to delete api key: %d", n), http.StatusBadRequest)
		tx.Rollback(r.Context())
		return
	}
	tx.Commit(r.Context())

}
