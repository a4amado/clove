package AppHandlersV1

import (
	"clove/internals/apperrors"
	"clove/internals/middleware"
	"clove/internals/heartbeat/dogpile"
	"clove/internals/meridian"
	"clove/internals/meridian/fanout"
	"clove/internals/services"
	appservice "clove/internals/services/apps"
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
	session, ok := middleware.SessionFromContext(r.Context())
	if !ok || !session.Permissions.Can(middleware.DELIVERY, middleware.READ) {
		middleware.UnAuthResponse(w)
		return
	}

	lock := sync.Mutex{}
	ctx := r.Context()
	wsUpgrader := websocket.Upgrader{}
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		apperrors.WriteWsError(conn, &lock, &apperrors.AppError{
			ID:         uuid.New(),
			Code:       ERROR_USER_CONNECT_WEBSOCKET_UPGRADE,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	defer conn.Close()

	appUUID, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		apperrors.WriteWsError(conn, &lock, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_INVALID_APP_ID,
			StatusCode: http.StatusBadRequest,
			ID:         uuid.New(),
		})
		return
	}

	if session.AppID.Bytes != appUUID {
		apperrors.WriteWsError(conn, &lock, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_TOKEN_APP_MISMATCH,
			StatusCode: http.StatusForbidden,
			ID:         uuid.New(),
		})
		return
	}

	srvs := services.New(r.Context())
	_, err = srvs.Apps.Get(appservice.GetParams{AppID: appUUID})
	if err != nil {
		apperrors.WriteWsError(conn, &lock, &apperrors.AppError{
			Code:       ERROR_USER_CONNECT_APP_NOT_FOUND,
			StatusCode: http.StatusForbidden,
			ID:         uuid.New(),
		})
		return
	}

	fanoutClient := meridian.Client().Fanout()
	channelKey := fanoutClient.FormatChannelKey(fanout.ChannelKey{
		AppID:     appUUID,
		ChannelID: session.ChannelID,
	})

	pubsub := fanout.Fanout().Subscribe(ctx, channelKey)
	defer pubsub.Close()

	messageChannel := pubsub.Channel()
	writeChan := make(chan []byte, 100)
	defer close(writeChan)

	dogpileInstance.Increase()
	defer dogpileInstance.Decrease()

	connCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	writeLock := sync.Mutex{}

	wg.Go(func() {
		defer cancel()
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
	})

	wg.Add(1)
	wg.Go(func() {
		defer wg.Done()
		defer cancel()

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
	})

	wg.Add(1)
	wg.Go(func() {
		defer cancel()

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
	})

	wg.Wait()
}

func writeToWebSocketWithLock(ctx context.Context, conn *websocket.Conn, lock *sync.Mutex, msgType int, data []byte) error {
	lock.Lock()
	defer lock.Unlock()

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.WriteMessage(msgType, data)
}
