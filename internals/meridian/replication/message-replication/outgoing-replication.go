package MessageReplication

import (
	"clove/internals/repository"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// InternalReplicatableDeliveryMsg
// the end user should always get the payload only
type InternalReplicatableDeliveryMsg struct {
	AppID     uuid.UUID
	ChannelId string
	Payload   []byte
}

// PublishReplicatableAppMsgToKafka publishes an app replication message to Kafka
// for distribution to other regions. It first attempts to save locally, and on
// failure sends to all regions to ensure eventual consistency.
func (c *MessageReplication) PublishInternalReplicatableDeliveryMsgToKafka(ctx context.Context, msg InternalReplicatableDeliveryMsg, regions []repository.Region) []error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return []error{err}
	}
	err = c.localKafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.ChannelId),
		Value: payload,
	})
	if err != nil {
		return []error{err}
	}
	errList := []error{}
	for _, region := range regions {
		err := c.crossRegionWriters[region].WriteMessages(ctx, kafka.Message{
			Key:   []byte(msg.ChannelId),
			Value: payload,
		})
		if err != nil {
			errList = append(errList, err)
		}
	}
	return errList
}

func (c *MessageReplication) PublishInternalReplicatableDeliveryMsgToKafkaGlobaly(ctx context.Context, msg InternalReplicatableDeliveryMsg) []error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return []error{err}
	}
	err = c.localKafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.ChannelId),
		Value: payload,
	})
	if err != nil {
		return []error{err}
	}
	errList := []error{}
	for _, region := range c.crossRegionWriters {
		err := region.WriteMessages(ctx, kafka.Message{
			Key:   []byte(msg.ChannelId),
			Value: payload,
		})
		if err != nil {
			errList = append(errList, err)
		}
	}
	return errList
}
func (c *MessageReplication) PublishInternalReplicatableDeliveryMsgToLocalKafka(ctx context.Context, msg InternalReplicatableDeliveryMsg) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = c.localKafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.ChannelId),
		Value: payload,
	})
	return err
}
