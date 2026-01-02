package AppHandlersV1

import (
	"clove/internals/apiguard"
	"clove/internals/meridian"
	MessageReplication "clove/internals/meridian/replication/message-replication"
	"clove/internals/services"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type MessageEntryBody struct {
	ChannelID string    `json:"channel_id"`
	Payload   string    `json:"payload"`
	AppKeyId  uuid.UUID `json:"app_key_id"`
}

func MessageEntry(w http.ResponseWriter, r *http.Request) {

	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		http.Error(w, "Invalid App ID", http.StatusUnauthorized)
		return
	}

	apiHeadersKey := apiguard.GetHeaderApi(r)

	// Wrap the body with MaxBytesReader to enforce size limit
	r.Body = http.MaxBytesReader(w, r.Body, 1_000_000)
	defer r.Body.Close()

	body := MessageEntryBody{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad Body", http.StatusBadRequest)
		return
	}

	appSrvs := services.C(r.Context(), nil, true)
	Apikey, err := appSrvs.App(appId).Key(body.AppKeyId).Get()
	if err != nil {
		http.Error(w, "Invalid ApiKey", http.StatusBadRequest)
		return
	}
	if Apikey.String != apiHeadersKey {
		http.Error(w, "Invalid ApiKey", http.StatusBadRequest)
		return
	}
	errList := meridian.Client().ReplicateMessage().PublishInternalReplicatableDeliveryMsgToLocalRabbitMQ(r.Context(), MessageReplication.InternalReplicatableDeliveryMsg{
		ChannelID: body.ChannelID,
		AppID:     appId,
		Payload:   []byte(body.Payload),
	})
	if errList != nil {
		http.Error(w, "Failed to replicate"+errList.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
