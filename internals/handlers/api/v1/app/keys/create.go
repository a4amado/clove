package AppKeysHandlersV1

import (
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/services"
	repository "clove/internals/services/generatedRepo"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateAppApiTokenBody struct {
	Name string `json:"name"`
}

func CreateAppApiKey(w http.ResponseWriter, r *http.Request) {
	apId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	tx, _ := postgresPool.NewTx(r.Context(), pgx.TxOptions{})

	body := CreateAppApiTokenBody{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Failed To Parse Body", http.StatusBadRequest)
		tx.Rollback(r.Context())
		return
	}
	api, err := services.App(r.Context(), &tx, true, apId).Keys().Generate(repository.CreateAppApiKeyParams{
		AppID: pgtype.UUID{Bytes: apId, Valid: true},
		Key: pgtype.Text{
			String: uuid.NewString(),
			Valid:  true,
		},
	})
	if err != nil {
		tx.Rollback(r.Context())
		http.Error(w, "Failed To Create Token", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(api); err != nil {
		tx.Rollback(r.Context())
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}
	tx.Commit(r.Context())

}
