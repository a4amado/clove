package AppHandlersV1

import (
	"clove/internals/meridian"
	MessageReplication "clove/internals/meridian/replication/message-replication"
	"io"
	"net/http"

	"github.com/google/uuid"
)

// MessageEntry handles an HTTP request that publishes a message to the websocket fanout for an app channel.
// 
// It parses `app_id` and `channel_id` from the request path, enforces a 1 MB request body limit, reads the
// request payload, and forwards it to the replication service as an InternalReplicatableDeliveryMsg.
// On success it responds with HTTP 202 Accepted. It responds with HTTP 400 for invalid UUIDs or body read
// failures, and HTTP 500 if publishing to the replication service fails.
func MessageEntry(w http.ResponseWriter, r *http.Request) {
	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		http.Error(w, "Invalid appId", http.StatusBadRequest)
		return
	}

	channelID := r.PathValue("channel_id")

	// Wrap the body with MaxBytesReader to enforce size limit
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1000)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	err = meridian.Client().ReplicateMessage().PublishFanoutMessageToWebsocket(r.Context(), MessageReplication.InternalReplicatableDeliveryMsg{
		ChannelId: channelID,
		AppID:     appId,
		Payload:   body,
	})
	if err != nil {
		http.Error(w, "Operation Failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}