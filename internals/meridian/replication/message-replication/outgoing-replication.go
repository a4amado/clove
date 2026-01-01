package MessageReplication

import (
	repository "clove/internals/services/generatedRepo"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// InternalReplicatableDeliveryMsg
// the end user should always get the payload only
type InternalReplicatableDeliveryMsg struct {
	AppID     uuid.UUID
	ChannelID string
	Payload   []byte
}

// PublishInternalReplicatableDeliveryMsgToRabbitMQ publishes a message replication message to RabbitMQ
// for distribution to other regions.
func (mr *MessageReplication) PublishInternalReplicatableDeliveryMsgToRabbitMQ(ctx context.Context, msg InternalReplicatableDeliveryMsg, regions []repository.Region) []error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return []error{err}
	}

	errList := []error{}

	// Publish to local region
	routingKey := fmt.Sprintf("%s.msg-replication", mr.localRegion)
	err = mr.localRabbitMQChannel.PublishWithContext(
		ctx,
		mr.exchangeName, // exchange
		routingKey,      // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			MessageId:    msg.ChannelID,
		},
	)
	if err != nil {
		errList = append(errList, err)
	}

	// Publish to cross-regions
	for _, region := range regions {
		ch := mr.crossRegionChannels[region]
		if ch == nil {
			continue
		}
		routingKey := fmt.Sprintf("%s.msg-replication", region)
		err := ch.PublishWithContext(
			ctx,
			mr.exchangeName, // exchange
			routingKey,      // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         payload,
				DeliveryMode: amqp.Persistent,
				Timestamp:    time.Now(),
				MessageId:    msg.ChannelID,
			},
		)
		if err != nil {
			errList = append(errList, err)
		}
	}
	return errList
}

func (mr *MessageReplication) PublishInternalReplicatableDeliveryMsgToRabbitMQGlobally(ctx context.Context, msg InternalReplicatableDeliveryMsg) []error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return []error{err}
	}
	errList := []error{}
	// Publish to local region
	routingKey := fmt.Sprintf("%s.msg-replication", mr.localRegion)
	err = mr.localRabbitMQChannel.PublishWithContext(
		ctx,
		mr.exchangeName, // exchange
		routingKey,      // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			MessageId:    msg.ChannelID,
		},
	)
	if err != nil {
		errList = append(errList, err)
	}

	// Publish to all cross-regions
	for region, ch := range mr.crossRegionChannels {
		if ch == nil {
			continue
		}
		routingKey := fmt.Sprintf("%s.msg-replication", region)
		err := ch.PublishWithContext(
			ctx,
			mr.exchangeName, // exchange
			routingKey,      // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         payload,
				DeliveryMode: amqp.Persistent,
				Timestamp:    time.Now(),
				MessageId:    msg.ChannelID,
			},
		)
		if err != nil {
			errList = append(errList, err)
		}
	}
	return errList
}

func (mr *MessageReplication) PublishInternalReplicatableDeliveryMsgToLocalRabbitMQ(ctx context.Context, msg InternalReplicatableDeliveryMsg) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	routingKey := fmt.Sprintf("%s.msg-replication", mr.localRegion)
	err = mr.localRabbitMQChannel.PublishWithContext(
		ctx,
		mr.exchangeName, // exchange
		routingKey,      // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			MessageId:    msg.ChannelID,
		},
	)
	return err
}
