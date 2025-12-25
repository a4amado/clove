package fanout

import (
	redisPool "clove/internals/data/redispool"
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type FanOut struct {
	conn *redis.Client
}

var fanoutOnce = sync.Once{}
var fanout *FanOut

func Fanout() *FanOut {
	fanoutOnce.Do(func() {
		fanout = &FanOut{
			conn: redisPool.Client(redisPool.RedisFanout),
		}
	})
	return fanout
}

func (f *FanOut) Publish(ctx context.Context, channel string, message any) *redis.IntCmd {
	return f.conn.Publish(ctx, channel, message)
}
func (f *FanOut) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return f.conn.Subscribe(ctx, channels...)
}

type ChannelKey struct {
	AppId     uuid.UUID
	ChannelId string
}

func (c *FanOut) FormatChannelKey(key ChannelKey) string {
	return fmt.Sprintf("%s:%s", key.AppId.String(), key.ChannelId)
}
