// this endpoint is
package AppHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/auth/apiguard"
	"clove/internals/meridian"
	MessageReplication "clove/internals/meridian/replication/message-replication"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
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
	ERROR_MESSAGE_ENTRY_BODY_TOO_LARGE       = "ERROR_MESSAGE_ENTRY_BODY_TOO_LARGE"
	ERROR_MESSAGE_ENTRY_BAD_BODY             = "ERROR_MESSAGE_ENTRY_BAD_BODY"
)

func WSMessageEntry(w http.ResponseWriter, r *http.Request) {
	appId, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_INVALID_APP_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,

			ID: uuid.New(),
		})
		return
	}

	app_key_id, err := uuid.Parse(r.URL.Query().Get("app_key_id"))
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_INVALID_APP_KEY_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,

			ID: uuid.New(),
		})
		return
	}

	channel_id := r.URL.Query().Get("channel_id")
	if channel_id == "" {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_MISSING_CHANNEL_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,

			ID: uuid.New(),
		})
		return
	}

	apiHeadersKey := apiguard.GetHeaderApi(r)

	r.Body = http.MaxBytesReader(w, r.Body, 50*1024)
	defer r.Body.Close()

	srvs := services.New(r.Context()).WithCache()

	Apikey, err := srvs.Apps.Keys.Get(appservice.GetKeyParams{
		KeyID: app_key_id,
		AppID: appId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			apperrors.WriteError(w, &apperrors.AppError{
				Code:       ERROR_MESSAGE_ENTRY_APP_KEY_NOT_FOUND,
				Message:    "",
				StatusCode: http.StatusNotFound,

				ID: uuid.New(),
			})
			return
		}
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_FAILED_FETCH_APP_KEY,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	if Apikey.Key.String != apiHeadersKey {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_UNAUTHORIZED_API_KEY,
			Message:    "",
			StatusCode: http.StatusUnauthorized,

			ID: uuid.New(),
		})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError

		if errors.As(err, &maxBytesErr) {
			apperrors.WriteError(w, &apperrors.AppError{
				Code:       ERROR_MESSAGE_ENTRY_BODY_TOO_LARGE,
				Message:    "",
				StatusCode: http.StatusBadRequest,

				ID: uuid.New(),
			})
		} else {

			apperrors.WriteError(w, &apperrors.AppError{
				Code:       ERROR_MESSAGE_ENTRY_BAD_BODY,
				Message:    "",
				StatusCode: http.StatusBadRequest,

				ID: uuid.New(),
			})
		}

		return
	}

	if len(body) > 32*1024 {

		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_MESSAGE_ENTRY_BODY_TOO_LARGE,
			Message:    "",
			StatusCode: http.StatusBadRequest,

			ID: uuid.New(),
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
			ID:         uuid.New(),
		})
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
