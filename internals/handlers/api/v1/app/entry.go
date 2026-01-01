package AppHandlersV1

import (
	"clove/internals/meridian"
	MessageReplication "clove/internals/meridian/replication/message-replication"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type MessageEntryBody struct {
	ChannelID string `json:"channel_id"`
	Payload   string `json:"payload"`
}

func MessageEntry(w http.ResponseWriter, r *http.Request) {

	appIdStr := r.PathValue("app_id")

	appId, err := uuid.Parse(appIdStr)
	if err != nil {
		http.Error(w, "Invalid App ID", http.StatusUnauthorized)
		return
	}
	// !: disabled for dev
	// token := r.Header.Get("Authorization")
	// splitToken := strings.Split(token, " ")
	// if len(splitToken) != 2 {
	// 	http.Error(w, "Invalid Token", http.StatusUnauthorized)
	// 	return
	// }

	// claims, err := tokenguard.ValidateOneTimeToken(token)
	// if err != nil {
	// 	http.Error(w, "Invalid Token", http.StatusForbidden)
	// 	return
	// }

	// // Wrap the body with MaxBytesReader to enforce size limit
	// r.Body = http.MaxBytesReader(w, r.Body, plans.GetPlanMessageSizeLimit(claims.App.AppType))
	// defer r.Body.Close()
	// if err != nil {
	// 	http.Error(w, "Failed to read body", http.StatusBadRequest)
	// 	return
	// }
	// if plans.DoesMessageSizeExceedsLimit(claims.App.AppType, int64(len(body))) {
	// 	http.Error(w, "Message size exceeds limit", http.StatusRequestEntityTooLarge)
	// 	return
	// }

	body := MessageEntryBody{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad ssssssssss", http.StatusBadRequest)
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
