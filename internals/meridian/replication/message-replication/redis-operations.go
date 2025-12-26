package MessageReplication

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

// Fanout message from kafka to the users
func (c *MessageReplication) PublishFanoutMessageToWebsocket(ctx context.Context, msg InternalReplicatableDeliveryMsg) error {
	adapter := valkeycompat.NewAdapter(c.conn)
	return adapter.Publish(ctx, c.FormatChannelKey(msg.AppID, msg.ChannelId), valkeycompat.BytesToString(msg.Payload)).Err()
}
func (c *MessageReplication) SubscribeFanoutMessage(ctx context.Context, channels []string) valkeycompat.PubSub {
	adapter := valkeycompat.NewAdapter(c.conn)
	return adapter.Subscribe(ctx, channels...)
}
func (c *MessageReplication) FormatChannelKey(app uuid.UUID, channel string) string {
	return fmt.Sprintf("%s:%s", app.String(), channel)
}
