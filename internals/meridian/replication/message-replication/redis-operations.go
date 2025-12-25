package MessageReplication

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Fanout message from kafka to the users
func (c *MessageReplication) PublishFanoutMessage(ctx context.Context, msg InternalReplicatableDeliveryMsg) error {
	res := c.conn.Publish(ctx, msg.ChannelId, msg.Payload)
	return res.Err()
}
func (c *MessageReplication) SubscribeFanoutMessage(ctx context.Context, channels []string) *redis.PubSub {
	return c.conn.Subscribe(ctx, channels...)

}
