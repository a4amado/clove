package MessageReplication

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

// Fanout message from RabbitMQ to the users
func (mr *MessageReplication) PublishFanoutMessageToWebSocket(ctx context.Context, msg InternalReplicatableDeliveryMsg) error {
	appID, err := uuid.FromBytes([]byte(msg.AppID.String()))
	if err != nil {
		return err
	}
	adapter := valkeycompat.NewAdapter(mr.conn)
	return adapter.Publish(ctx, mr.FormatChannelKey(appID, msg.ChannelID), valkeycompat.BytesToString(msg.Payload)).Err()
}
func (mr *MessageReplication) SubscribeFanoutMessage(ctx context.Context, channels []string) valkeycompat.PubSub {
	adapter := valkeycompat.NewAdapter(mr.conn)
	return adapter.Subscribe(ctx, channels...)
}
func (mr *MessageReplication) FormatChannelKey(appID uuid.UUID, channelID string) string {
	return fmt.Sprintf("%s:%s", appID.String(), channelID)
}
