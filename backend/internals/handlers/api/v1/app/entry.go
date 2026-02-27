package AppHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/middleware"
	"clove/internals/meridian"
	MessageReplication "clove/internals/meridian/replication/message-replication"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
)

const (
	ERROR_MESSAGE_ENTRY_INVALID_APP_ID       = "ERROR_MESSAGE_ENTRY_INVALID_APP_ID"
	ERROR_MESSAGE_ENTRY_MISSING_CHANNEL_ID   = "ERROR_MESSAGE_ENTRY_MISSING_CHANNEL_ID"
	ERROR_MESSAGE_ENTRY_APP_KEY_NOT_FOUND    = "ERROR_MESSAGE_ENTRY_APP_KEY_NOT_FOUND"
	ERROR_MESSAGE_ENTRY_FAILED_FETCH_APP_KEY = "ERROR_MESSAGE_ENTRY_FAILED_FETCH_APP_KEY"
	ERROR_MESSAGE_ENTRY_UNAUTHORIZED_API_KEY = "ERROR_MESSAGE_ENTRY_UNAUTHORIZED_API_KEY"
	ERROR_MESSAGE_ENTRY_INVALID_REQUEST_BODY = "ERROR_MESSAGE_ENTRY_INVALID_REQUEST_BODY"
	ERROR_MESSAGE_ENTRY_REPLICATION_FAILED   = "ERROR_MESSAGE_ENTRY_REPLICATION_FAILED"
	ERROR_MESSAGE_ENTRY_BODY_TOO_LARGE       = "ERROR_MESSAGE_ENTRY_BODY_TOO_LARGE"
	ERROR_MESSAGE_ENTRY_BAD_BODY             = "ERROR_MESSAGE_ENTRY_BAD_BODY"
)

func WSMessageEntry(w http.ResponseWriter, r *http.Request) {
	session, ok := middleware.SessionFromContext(r.Context())
	if !ok || !session.Permissions.Can(middleware.DELIVERY, middleware.CREATE) {
		middleware.UnAuthResponse(w)
		return
	}

	appID, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_INVALID_APP_ID,
			StatusCode: http.StatusBadRequest,
			ID:         uuid.New(),
		})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 50*1024)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			apperrors.WriteError(w, &apperrors.AppError{
				Code:       ERROR_MESSAGE_ENTRY_BODY_TOO_LARGE,
				StatusCode: http.StatusBadRequest,
				ID:         uuid.New(),
			})
		} else {
			apperrors.WriteError(w, &apperrors.AppError{
				Code:       ERROR_MESSAGE_ENTRY_BAD_BODY,
				StatusCode: http.StatusBadRequest,
				ID:         uuid.New(),
			})
		}
		return
	}

	channelID := r.URL.Query().Get("channel_id")
	errList := meridian.Client().ReplicateMessage().PublishInternalReplicatableDeliveryMsgToLocalRabbitMQ(r.Context(), MessageReplication.InternalReplicatableDeliveryMsg{
		ChannelID: channelID,
		AppID:     appID,
		Payload:   body,
	})
	if errList != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_REPLICATION_FAILED,
			StatusCode: http.StatusInternalServerError,
			ID:         uuid.New(),
		})
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
