package AppHandlersV1

import (
	"clove/internals/apiguard"
	"clove/internals/apperrors"
	appConsts "clove/internals/consts/app"
	headers "clove/internals/handlers/api/response-utils/consts"
	"clove/internals/heartbeat/dogpile"
	"clove/internals/meridian"
	"clove/internals/meridian/fanout"
	"clove/internals/services"
	repository "clove/internals/services/generatedRepo"
	"clove/internals/tokenguard"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	set "github.com/hashicorp/go-set"
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
func newUpgrader(app *repository.App) websocket.Upgrader {
	bufferSize := appConsts.GetAppBufferSize(app.AppType)

	return websocket.Upgrader{
		WriteBufferSize: bufferSize,
		ReadBufferSize:  bufferSize,
		CheckOrigin: func(r *http.Request) bool {
			allowed := set.From(app.AllowedOrigins)
			// this works under the assumption that app.AllowedOrigins are normalized on insert
			origin := strings.ToLower(r.Header.Get(headers.Origin))
			// make sure the request is coming from the authorized domain
			return allowed.Contains(origin)
		},
	}
}

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
	ctx := r.Context()

	appUUID, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_INVALID_APP_ID,
			Message:    "",
			StatusCode: http.StatusBadRequest,
			Internal:   err,
		})
		return
	}

	token := apiguard.GetHeaderApi(r)
	claims, err := tokenguard.ValidateOneTimeToken(token)
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_INVALID_TOKEN,
			Message:    "",
			StatusCode: http.StatusUnauthorized,
			Internal:   err,
		})
		return
	}

	if claims.App.ID.String() != appUUID.String() {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_TOKEN_APP_MISMATCH,
			Message:    "",
			StatusCode: http.StatusForbidden,
			Internal:   nil,
		})
		return
	}

	app, err := services.C(ctx, nil, true).App(appUUID).Get()
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_APP_NOT_FOUND,
			Message:    "",
			StatusCode: http.StatusNotFound,
			Internal:   err,
		})
		return
	}

	wsUpgrader := newUpgrader(app)
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		apperrors.WriteError(w, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_WEBSOCKET_UPGRADE,
			Message:    "",
			StatusCode: http.StatusInternalServerError,
			Internal:   err,
		})
		return
	}
	defer conn.Close()

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
