package MessageReplication

import (
	"clove/internals/repository"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

// InternalReplicatableDeliveryMsg
// the end user should always get the payload only
type InternalReplicatableDeliveryMsg struct {
	ChannelId string
	Payload   []byte
}

func (m *InternalReplicatableDeliveryMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

// getCurrentMachineRegion returns the current machine's region.
// Panics if Init() has not been called.
func (m *MessageReplication) getCurrentMachineRegion() repository.Region {
	if !m.currentMachineRegion.Valid() {
		panic("meridian.Init() must be called before using getCurrentMachineRegion()")
	}
	return m.currentMachineRegion
}

// PublishReplicatableAppMsgToKafka publishes an app replication message to Kafka
// for distribution to other regions. It first attempts to save locally, and on
// failure sends to all regions to ensure eventual consistency.
func (c *MessageReplication) PublishInternalReplicatableDeliveryMsgToKafka(ctx context.Context, msg InternalReplicatableDeliveryMsg, regions []repository.Region) []error {
	payload, err := msg.MarshalJSON()
	if err != nil {
		return []error{err}
	}
	errList := []error{}
	for _, region := range regions {
		err := c.regionKafkaWriters[region].WriteMessages(ctx, kafka.Message{
			Key:   []byte(msg.ChannelId),
			Value: payload,
		})
		if err != nil {
			errList = append(errList, err)
		}
	}
	return errList
}
