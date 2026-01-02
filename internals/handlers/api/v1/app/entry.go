package AppHandlersV1

import (
	"clove/internals/apiguard"
	"clove/internals/apperrors"
	"clove/internals/meridian"
	MessageReplication "clove/internals/meridian/replication/message-replication"
	"clove/internals/services"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	ERROR_MESSAGE_ENTRY_INVALID_APP_ID       = "ERROR_MESSAGE_ENTRY_INVALID_APP_ID"
	ERROR_MESSAGE_ENTRY_INVALID_APP_KEY_ID   = "ERROR_MESSAGE_ENTRY_INVALID_APP_KEY_ID"
	ERROR_MESSAGE_ENTRY_MISSING_CHANNEL_ID   = "ERROR_MESSAGE_ENTRY_MISSING_CHANNEL_ID"
	ERROR_MESSAGE_ENTRY_APP_KEY_NOT_FOUND    = "ERROR_MESSAGE_ENTRY_APP_KEY_NOT_FOUND"
	ERROR_MESSAGE_ENTRY_FAILED_FETCH_APP_KEY = "ERROR_MESSAGE_ENTRY_FAILED_FETCH_APP_KEY"
	ERROR_MESSAGE_ENTRY_UNAUTHORIZED_API_KEY = "ERROR_MESSAGE_ENTRY_UNAUTHORIZED_API_KEY"
	ERROR_MESSAGE_ENTRY_INVALID_REQUEST_BODY = "ERROR_MESSAGE_ENTRY_INVALID_REQUEST_BODY"
	ERROR_MESSAGE_ENTRY_REPLICATION_FAILED   = "ERROR_MESSAGE_ENTRY_REPLICATION_FAILED"
)

func MessageEntry(w http.ResponseWriter, r *http.Request) {
	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_INVALID_APP_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			Request:    r,
			ID:         uuid.New(),
		})
		return
	}

	app_key_id, err := uuid.Parse(r.URL.Query().Get("app_key_id"))
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_INVALID_APP_KEY_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	channel_id := r.URL.Query().Get("channel_id")
	if channel_id == "" {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_MISSING_CHANNEL_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	apiHeadersKey := apiguard.GetHeaderApi(r)

	r.Body = http.MaxBytesReader(w, r.Body, 1_000_000)
	defer r.Body.Close()

	appSrvs := services.C(r.Context(), nil, true)
	Apikey, err := appSrvs.App(appId).Key(app_key_id).Get()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			apperrors.WriteError(w, &apperrors.AppError{
				Code:       ERROR_MESSAGE_ENTRY_APP_KEY_NOT_FOUND,
				Message:    "",
				StatusCode: http.StatusNotFound,
				Internal:   err,
				ID:         uuid.New(),
				Request:    r,
			})
			return
		}
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_FAILED_FETCH_APP_KEY,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
		})
		return
	}

	if Apikey.String != apiHeadersKey {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_UNAUTHORIZED_API_KEY,
			Message:    "",
			StatusCode: http.StatusUnauthorized,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_INVALID_REQUEST_BODY,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	errList := meridian.Client().ReplicateMessage().PublishInternalReplicatableDeliveryMsgToLocalRabbitMQ(r.Context(), MessageReplication.InternalReplicatableDeliveryMsg{
		ChannelID: channel_id,
		AppID:     appId,
		Payload:   body,
	})
	if errList != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_REPLICATION_FAILED,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   errList,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
