package MessageReplication

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Fanout message from kafka to the users
func (c *MessageReplication) PublishFanoutMessageToWebsocket(ctx context.Context, msg InternalReplicatableDeliveryMsg) error {
	res := c.conn.Publish(ctx, c.FormatChannelKey(msg.AppID, msg.ChannelId), msg.Payload)
	return res.Err()
}
func (c *MessageReplication) SubscribeFanoutMessage(ctx context.Context, channels []string) *redis.PubSub {
	return c.conn.Subscribe(ctx, channels...)

}
func (c *MessageReplication) FormatChannelKey(app uuid.UUID, channel string) string {
	return fmt.Sprintf("%s:%s", app.String(), channel)
}
