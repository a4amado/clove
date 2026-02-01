package fanout

import (
	"clove/internals/data/valkeyPool"
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

type FanOut struct {
	adapter valkeycompat.Cmdable
}

var fanoutOnce = sync.Once{}
var fanoutInstance *FanOut

func Fanout() *FanOut {
	fanoutOnce.Do(func() {
		fanoutInstance = &FanOut{
			adapter: valkeycompat.NewAdapter(valkeyPool.Client(valkeyPool.ValkeyFanout)),
		}
	})
	return fanoutInstance
}

func (f *FanOut) Publish(ctx context.Context, channel string, message any) error {
	var msgStr string
	switch v := message.(type) {
	case string:
		msgStr = v
	case []byte:
		msgStr = valkeycompat.BytesToString(v)
	default:
		msgStr = fmt.Sprintf("%v", v)
	}
	return f.adapter.Publish(ctx, channel, msgStr).Err()
}
func (f *FanOut) Subscribe(ctx context.Context, channels ...string) valkeycompat.PubSub {
	return f.adapter.Subscribe(ctx, channels...)
}

type ChannelKey struct {
	AppID     uuid.UUID
	ChannelID string
}

func (f *FanOut) FormatChannelKey(key ChannelKey) string {
	return fmt.Sprintf("%s:%s", key.AppID.String(), key.ChannelID)
}
