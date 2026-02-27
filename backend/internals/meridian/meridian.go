// Package meridian handles routing of user-generated events across Clove regions.
//
// Meridian uses RabbitMQ exclusively for cross-region delivery of user events.
// Global routing state (app registry, channel membership) is managed by the
// Valkey cluster and is not replicated through this package.
//
// Architecture:
//
//	Producer: Publishes user-generated messages to RabbitMQ for cross-region fanout
//	Consumer: Ingests messages from RabbitMQ and publishes them to local Valkey pub/sub
//	Fanout:   Valkey pub/sub delivers messages to connected WebSocket clients

package meridian

import (
	"clove/internals/meridian/fanout"
	MessageReplication "clove/internals/meridian/replication/message-replication"
	"sync"
)

type Meridian struct {
}

var meridianInstance *Meridian
var meridianOnce = sync.Once{}

// Client returns the singleton Meridian instance with role-specific Valkey connections
// (store, fan-out, heartbeat) initialized from the package valkeyPool.
// Initialization is performed exactly once and is safe for concurrent use.
func Client() *Meridian {
	meridianOnce.Do(func() {
		meridianInstance = &Meridian{}
	})
	return meridianInstance
}

func (mer *Meridian) Fanout() *fanout.FanOut {
	return fanout.Fanout()
}

func (mer *Meridian) ReplicateMessage() *MessageReplication.MessageReplication {
	return MessageReplication.ReplicateMessage()
}
