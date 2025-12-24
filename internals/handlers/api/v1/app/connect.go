package AppHandlersV1

import (
	appConsts "clove/internals/consts/app"
	dbPool "clove/internals/data/database/pool"
	redisPool "clove/internals/data/redispool"
	headers "clove/internals/handlers/api/response-utils/consts"
	"clove/internals/meridian"
	"clove/internals/repository"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	set "github.com/hashicorp/go-set"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
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
			allowed := set.From(app.AllowedOrigins)
			// this works under the assumption that app.AllowedOrigins are normalized on insert
			origin := strings.ToLower(r.Header.Get(headers.Origin))
			// make sure the request is comming from the authrized domain
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

type Connection struct {
	conn      *websocket.Conn
	writeLock sync.Mutex
	app       *repository.App
	channels  []string
	subs      []*redis.PubSub
	send      chan []byte
	done      chan struct{}
	closeOnce sync.Once
}

// NewConnection creates a Connection for the given WebSocket and app, initialized to subscribe to the provided channels.
// The returned Connection has a buffered send channel, a done channel for shutdown signaling, and a pre-sized subscription slice.
func NewConnection(conn *websocket.Conn, app *repository.App, channels []string) *Connection {
	return &Connection{
		conn:     conn,
		app:      app,
		channels: channels,
		subs:     make([]*redis.PubSub, 0, len(channels)),
		send:     make(chan []byte, 256),
		done:     make(chan struct{}),
	}
}

func (c *Connection) Close() error {
	var err error
	c.closeOnce.Do(func() {
		close(c.done)

		for _, sub := range c.subs {
			if sub != nil {
				if closeErr := sub.Close(); closeErr != nil {
					log.Printf("Error closing subscription: %v", closeErr)
				}
			}
		}

		close(c.send)

		c.writeLock.Lock()
		c.conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			time.Now().Add(writeWait))
		err = c.conn.Close()
		c.writeLock.Unlock()
	})
	return err
}

func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {

				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.Close()
				return
			}

			if err := c.writeMessage(websocket.BinaryMessage, message); err != nil {
				log.Printf("Error writing message: %v", err)

				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error writing ping: %v", err)
				return
			}

		case <-c.done:
			return
		}
	}
}

func (c *Connection) writeMessage(messageType int, data []byte) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	return c.conn.WriteMessage(messageType, data)
}

func (c *Connection) readPump() {
	defer c.Close()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			}
			break
		}

	}
}

func (c *Connection) subscribeToChannels(ctx context.Context, meridianClient *meridian.Meridian) error {
	appUuid, err := uuid.FromBytes(c.app.ID.Bytes[:])
	if err != nil {
		return err
	}

	for _, channel := range c.channels {
		pubSub := meridianClient.RedisFanOutConn.Subscribe(ctx, meridianClient.FormatChannelKey(appUuid, channel))
		c.subs = append(c.subs, pubSub)

		go func(ps *redis.PubSub, channelID string) {
			defer func() {
				if r := recover(); r != nil {
					c.Close()
					log.Printf("recovered from panic in subscription goroutine: %v", r)
				}
			}()

			ch := ps.Channel()
			for {
				select {
				case message, ok := <-ch:
					if !ok {
						log.Printf("Channel closed for %s", channelID)
						c.Close()
						return
					}

					select {
					case c.send <- []byte(message.Payload):
					case <-time.After(5 * time.Second):
						log.Printf("Send timeout for channel %s, client may be slow", channelID)

						c.Close()
						return
					case <-c.done:
						return
					}

				case <-c.done:
					return

				}
			}
		}(pubSub, channel)
	}

	return nil
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

	id := r.PathValue("app_id")
	channels := r.URL.Query().Get("channel")

	if len(channels) == 0 {
		http.Error(w, "No channels specified", http.StatusBadRequest)
		return
	}

	appUuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid app ID", http.StatusBadRequest)
		return
	}

	appID := pgtype.UUID{
		Bytes: appUuid,
		Valid: true,
	}

	// Try cache
	app, err := meridian.Client().FetchApp(ctx, appID)

	// Real cache error (not a miss) - fail fast
	if err != nil && !errors.Is(err, redisPool.ErrCacheMiss) {
		log.Printf("Cache error fetching app %s: %v", appID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Cache miss - fetch from database
	if app == nil {
		dbApp, err := repository.New(dbPool.Client()).FindAppById(ctx, appID)
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "App not found", http.StatusNotFound)
			return
		}
		if err != nil {
			log.Printf("Database error fetching app %s: %v", appID, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		app = &dbApp

		// Update cache asynchronously
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := meridian.Client().SaveApp(ctx, *app); err != nil {
				log.Printf("Failed to cache app %s: %v", appID, err)
			}
		}()
	}

	// Select upgrader based on app type
	var upgrader = newUpgrader(app)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed for app %s: %v", appID, err)
		return
	}

	connection := NewConnection(conn, app, []string{channels})

	if err := connection.subscribeToChannels(ctx, meridian.Client()); err != nil {
		log.Printf("Error subscribing to channels for app %s: %v", appID, err)
		connection.Close()
		return
	}

	go connection.writePump()
	go connection.readPump()
}
