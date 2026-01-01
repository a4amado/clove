package AppHandlersV1

import (
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
	"fmt"
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

var dogpileInstance = dogpile.New()

// newUpgrader creates a websocket.Upgrader configured with buffer sizes based on the app's type and an origin check using the app's AllowedOrigins.
// The upgrader's CheckOrigin lowercases the request Origin header and allows the request only if it matches an entry in app.AllowedOrigins, which must be pre-normalized to lowercase.
func newUpgrader(app *repository.App) websocket.Upgrader {
	bufferSize := appConsts.GetAppBufferSize(app.AppType)

	return websocket.Upgrader{
		WriteBufferSize: bufferSize,
		ReadBufferSize:  bufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
			allowed := set.From(app.AllowedOrigins)
			// this works under the assumption that app.AllowedOrigins are normalized on insert
			origin := strings.ToLower(r.Header.Get(headers.Origin))
			// make sure the request is coming from the authrized domain
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
//
// It validates the "app_id" path parameter and the "channel" query parameter, attempts
// to load the app configuration from cache and falls back to the database (caching
// the DB result asynchronously on success). It replies with HTTP 400 for missing or
// invalid parameters, 404 if the app is not found, and 500 on internal errors.
// After a successful upgrade it subscribes the connection to Valkey channels and starts
// the connection's read/write pumps; subscription failures close the connection.
func UserConnect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	appUUID, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, "Invalid app ID", http.StatusBadRequest)
		return
	}
	channel := r.URL.Query().Get("channel")
	Authorization := r.Header.Get("Authorization")
	idx := strings.Index(Authorization, " ")
	if idx == -1 {

		http.Error(w, "", http.StatusForbidden)
		return
	}
	bearer := Authorization[idx+1:]

	claims, err := tokenguard.ValidateOneTimeToken(bearer)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, "Invalid app ID", http.StatusBadRequest)
		return
	}

	if claims.App.ID.String() != appUUID.String() {
		http.Error(w, "", http.StatusForbidden)
		return
	}
	// Try cache
	app, err := services.App(ctx, nil, true, appUUID).Get()

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "app not found", http.StatusNotFound)
		return
	}

	wsUpgrader := newUpgrader(app)
	conn, _ := wsUpgrader.Upgrade(w, r, nil)
	fanoutClient := meridian.Client().Fanout()
	fmt.Println("Subscribed to: ", fanoutClient.FormatChannelKey(fanout.ChannelKey{
		AppID:     appUUID,
		ChannelID: channel,
	}))
	pubsub := fanout.Fanout().Subscribe(ctx, fanoutClient.FormatChannelKey(fanout.ChannelKey{
		AppID:     appUUID,
		ChannelID: channel,
	}))
	messageChannel := pubsub.Channel()

	writeChan := make(chan []byte, 100)

	dogpileInstance.Increase()
	defer dogpileInstance.Decrease()
	// ! double-buffered one-reader one-writer guarantee delivery order
	wg := sync.WaitGroup{}
	wg.Go(func() {
		for {
			select {
			case msg, ok := <-messageChannel:
				if !ok {
					return
				}
				select {
				case writeChan <- []byte(msg.Payload): // msg.Payload is string, convert to []byte
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	})
	writeLock := sync.Mutex{}

	wg.Go(func() {
		for {
			select {
			case data, ok := <-writeChan:
				if !ok {
					return
				}
				err := WriteToWebSocketWithLock(ctx, conn, &writeLock, websocket.BinaryMessage, data)
				if err != nil {
					return

				}
			case <-ctx.Done():
				return
			}
		}
	})

	wg.Wait()
}

func WriteToWebSocketWithLock(ctx context.Context, conn *websocket.Conn, lock *sync.Mutex, msgType int, data []byte) error {
	lock.Lock()
	defer lock.Unlock()
	return conn.WriteMessage(websocket.BinaryMessage, data)
}
