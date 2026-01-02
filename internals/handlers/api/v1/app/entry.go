package AppHandlersV1

import (
	"clove/internals/apiguard"
	"clove/internals/meridian"
	MessageReplication "clove/internals/meridian/replication/message-replication"
	"clove/internals/services"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func MessageEntry(w http.ResponseWriter, r *http.Request) {

	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		http.Error(w, "Invalid App ID", http.StatusUnauthorized)
		return
	}
	app_key_id, err := uuid.Parse(r.URL.Query().Get("app_key_id"))
	if err != nil {
		http.Error(w, "Invalid key ID", http.StatusBadRequest)
		return
	}
	channel_id := r.URL.Query().Get("channel_id")
	if channel_id == "" {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}
	apiHeadersKey := apiguard.GetHeaderApi(r)

	// Wrap the body with MaxBytesReader to enforce size limit
	r.Body = http.MaxBytesReader(w, r.Body, 1_000_000)
	defer r.Body.Close()

	appSrvs := services.C(r.Context(), nil, true)
	Apikey, err := appSrvs.App(appId).Key(app_key_id).Get()
	if err != nil {
		http.Error(w, "Invalid ApiKey", http.StatusBadRequest)
		return
	}
	if Apikey.String != apiHeadersKey {
		http.Error(w, "Invalid ApiKey", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}
	errList := meridian.Client().ReplicateMessage().PublishInternalReplicatableDeliveryMsgToLocalRabbitMQ(r.Context(), MessageReplication.InternalReplicatableDeliveryMsg{
		ChannelID: channel_id,
		AppID:     appId,
		Payload:   body,
	})
	if errList != nil {
		http.Error(w, "Failed to replicate"+errList.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
