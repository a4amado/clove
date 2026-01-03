package AppHandlersV1

import (
	"clove/internals/apiguard"
	"clove/internals/apperrors"
	"clove/internals/heartbeat/dogpile"
	"clove/internals/meridian"
	"clove/internals/meridian/fanout"
	"clove/internals/services"
	"clove/internals/tokenguard"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024
)

const (
	ERROR_USER_CONNECT_INVALID_APP_ID     = "ERROR_USER_CONNECT_INVALID_APP_ID"
	ERROR_USER_CONNECT_INVALID_TOKEN      = "ERROR_USER_CONNECT_INVALID_TOKEN"
	ERROR_USER_CONNECT_TOKEN_APP_MISMATCH = "ERROR_USER_CONNECT_TOKEN_APP_MISMATCH"
	ERROR_USER_CONNECT_APP_NOT_FOUND      = "ERROR_USER_CONNECT_APP_NOT_FOUND"
	ERROR_USER_CONNECT_WEBSOCKET_UPGRADE  = "ERROR_USER_CONNECT_WEBSOCKET_UPGRADE"
)

var dogpileInstance = dogpile.New()

// newUpgrader creates a websocket.Upgrader configured with buffer sizes based on the app's type and an origin check using the app's AllowedOrigins.
// The upgrader's CheckOrigin lowercases the request Origin header and allows the request only if it matches an entry in app.AllowedOrigins, which must be pre-normalized to lowercase.

type MessageToClient struct {
	Channel string `json:"channel"`
	Payload []byte `json:"payload"`
}

func (m *MessageToClient) Binary() ([]byte, error) {
	return json.Marshal(m)
}

// UserConnect upgrades the incoming HTTP request to a WebSocket for the specified app
// and subscribes the resulting connection to the requested channel(s).
func UserConnect(w http.ResponseWriter, r *http.Request) {
	lock := sync.Mutex{}
	ctx := r.Context()
	wsUpgrader := websocket.Upgrader{}
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		apperrors.WriteWsError(conn, &lock, &apperrors.AppError{
			ID:         uuid.New(),
			Request:    r,
			Code:       ERROR_USER_CONNECT_WEBSOCKET_UPGRADE,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
		})
		return
	}
	defer conn.Close()
	appUUID, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteWsError(conn, &lock, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_INVALID_APP_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	token := apiguard.GetHeaderApi(r)
	claims, err := tokenguard.ValidateOneTimeToken(token)
	if err != nil {
		apperrors.WriteWsError(conn, &lock, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_INVALID_TOKEN,
			Message:    "",
			StatusCode: http.StatusUnauthorized,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	if claims.App.ID.String() != appUUID.String() {
		apperrors.WriteWsError(conn, &lock, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_TOKEN_APP_MISMATCH,
			Message:    "",
			StatusCode: http.StatusForbidden,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	_, err = services.C(ctx, nil, true).App(appUUID).Get()
	if err != nil {
		apperrors.WriteWsError(conn, &lock, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_APP_NOT_FOUND,
			Message:    "",
			StatusCode: http.StatusForbidden,
			Internal:   err,
			ID:         uuid.New(),
			Request:    r,
		})
		return
	}

	fanoutClient := meridian.Client().Fanout()
	channelKey := fanoutClient.FormatChannelKey(fanout.ChannelKey{
		AppID:     appUUID,
		ChannelID: claims.ChannelID,
	})

	pubsub := fanout.Fanout().Subscribe(ctx, channelKey)
	defer pubsub.Close()

	messageChannel := pubsub.Channel()
	writeChan := make(chan []byte, 100)
	defer close(writeChan)

	dogpileInstance.Increase()
	defer dogpileInstance.Decrease()

	// Context that cancels when connection dies
	connCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	writeLock := sync.Mutex{}

	// Reader goroutine (processes messages from Valkey)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel() // Cancel on exit

		for {
			select {
			case msg, ok := <-messageChannel:
				if !ok {
					return
				}
				select {
				case writeChan <- []byte(msg.Payload):
				case <-connCtx.Done():
					return
				}
			case <-connCtx.Done():
				return
			}
		}
	}()

	// Writer goroutine (writes to WebSocket)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel() // Cancel on exit

		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		for {
			select {
			case data, ok := <-writeChan:
				if !ok {
					return
				}
				if err := writeToWebSocketWithLock(connCtx, conn, &writeLock, websocket.BinaryMessage, data); err != nil {
					return
				}
			case <-ticker.C:
				if err := writeToWebSocketWithLock(connCtx, conn, &writeLock, websocket.PingMessage, nil); err != nil {
					return
				}
			case <-connCtx.Done():
				return
			}
		}
	}()

	// Read pump (handles pongs and close messages)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel() // Cancel on exit

		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		conn.SetReadLimit(maxMessageSize)

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	wg.Wait()
}

func writeToWebSocketWithLock(ctx context.Context, conn *websocket.Conn, lock *sync.Mutex, msgType int, data []byte) error {
	lock.Lock()
	defer lock.Unlock()

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.WriteMessage(msgType, data)
}
