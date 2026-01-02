package AppHandlersV1

import (
	"clove/internals/apiguard"
	postgresPool "clove/internals/data/postgres/pool"
	"clove/internals/services"
	repository "clove/internals/services/generatedRepo"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateAppStruct struct {
	AppSlug        string              `json:"app_slug"`
	Regions        []repository.Region `json:"regions"`
	UserId         pgtype.UUID         `json:"user_id"`
	AllowedOrigins []string            `json:"allowed_origins"`
}
type CreateAppReponse struct {
	repository.App `json:"app"`
	Keys           []repository.AppApiKey `json:"keys"`
}

func CreateApp(w http.ResponseWriter, r *http.Request) {

	body := CreateAppStruct{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {

		http.Error(w, "Failed to parse the body "+err.Error(), http.StatusBadRequest)
		return
	}

	tx, _ := postgresPool.NewTx(r.Context(), pgx.TxOptions{})

	app, err := services.C(r.Context(), &tx, true).Apps().Create(repository.InsertAppParams{
		AppSlug:        fmt.Sprintf("%s:%s", uuid.NewString(), body.AppSlug),
		Regions:        body.Regions,
		AppType:        repository.AppTypePro,
		UserID:         body.UserId,
		AllowedOrigins: body.AllowedOrigins,
	})
	if err != nil {
		http.Error(w, "Failed create and app", http.StatusBadRequest)
		tx.Rollback(r.Context())
		return
	}

	key, err := apiguard.RandomSecretKey()
	if err != nil {
		http.Error(w, "Failed generate initial Key", http.StatusInternalServerError)
		tx.Rollback(r.Context())
		return
	}

	appApiKey, err := services.C(r.Context(), &tx, true).App(app.App.ID.Bytes).Keys().Create(repository.CreateAppApiKeyParams{
		AppID: app.App.ID,
		Key: pgtype.Text{
			String: key,
			Valid:  true,
		},
		Name: pgtype.Text{
			String: "Clove: Initial Auto Generated Key",
		},
	})
	if err != nil {
		http.Error(w, "Failed to create app key", http.StatusBadRequest)
		tx.Rollback(r.Context())
		return
	}
	res := CreateAppReponse{
		App:  app.App,
		Keys: []repository.AppApiKey{*appApiKey},
	}
	json.NewEncoder(w).Encode(res)
	tx.Commit(r.Context())
}
