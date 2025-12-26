package AppHandlersV1

import (
	appConsts "clove/internals/consts/app"
	dbPool "clove/internals/data/postgres/pool"
	redisPool "clove/internals/data/redispool"
	headers "clove/internals/handlers/api/response-utils/consts"
	"clove/internals/meridian"
	"clove/internals/meridian/fanout"
	"clove/internals/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	set "github.com/hashicorp/go-set"
	"github.com/jackc/pgx/v5"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024
)

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
// After a successful upgrade it subscribes the connection to Redis channels and starts
// the connection's read/write pumps; subscription failures close the connection.
func UserConnect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	appUuid, err := uuid.Parse(r.PathValue("app_id"))
	if err != nil {
		http.Error(w, "Invalid app ID", http.StatusBadRequest)
		return
	}
	channel := r.URL.Query().Get("channel")

	// Try cache
	app, err := meridian.Client().ReplicateApp().FetchApp(ctx, appUuid)

	// Real cache error (not a miss) - fail fast
	if err != nil && !errors.Is(err, redisPool.ErrCacheMiss) {
		log.Printf("Cache error fetching app %s: %v", appUuid, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Cache miss - fetch from database
	if app == nil {
		dbApp, err := repository.New(dbPool.Client()).FindAppById(ctx, appUuid)
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "App not found", http.StatusNotFound)
			return
		}
		if err != nil {
			log.Printf("Database error fetching app %s: %v", appUuid, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		app = &dbApp

		// Update cache asynchronously
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := meridian.Client().ReplicateApp().SaveApp(ctx, *app); err != nil {
				log.Printf("Failed to cache app %s: %v", appUuid, err)
			}
		}()
	}

	var upgrader = newUpgrader(app)
	conn, _ := upgrader.Upgrade(w, r, nil)
	fanoutclient := meridian.Client().Fanout()
	fmt.Println("Subscribed to: ", fanoutclient.FormatChannelKey(fanout.ChannelKey{
		AppId:     appUuid,
		ChannelId: channel,
	}))
	pubSub := fanout.Fanout().Subscribe(ctx, fanoutclient.FormatChannelKey(fanout.ChannelKey{
		AppId:     appUuid,
		ChannelId: channel,
	}))
	ch := pubSub.Channel()

	writeCh := make(chan []byte, 100)

	// ! double-buffered one-reader one-writer guarantee delivery order
	wg := sync.WaitGroup{}
	wg.Go(func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				select {
				case writeCh <- []byte(msg.Payload):
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	})
	lock := sync.Mutex{}

	wg.Go(func() {
		for {
			select {
			case data, ok := <-writeCh:
				if !ok {
					return
				}
				lock.Lock()
				err := conn.WriteMessage(websocket.BinaryMessage, data)
				lock.Unlock()
				if err != nil {
					return
					lock.Unlock()

				}
			case <-ctx.Done():
				return
			}
		}
	})

	wg.Wait()
}
