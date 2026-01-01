// Package meridian is the core routing coordination service for Clove.
//
// Meridian maintains the global routing registry through a multi-tier storage architecture:
//   - PostgreSQL: Source of truth for persistent routing state
//   - Valkey: Local cache for sub-millisecond route lookups
//   - RabbitMQ: Event stream for global state propagation
//
// Architecture:
//
//	Producer: Publishes local routing changes to RabbitMQ for global replication
//	Consumer: Ingests routing updates from RabbitMQ and materializes them in local Valkey
//	Query Interface: Serves routing lookups from Valkey with PostgreSQL fallback
//
// This design ensures eventual consistency across distributed Clove instances while
// maintaining low-latency access to routing data.

package meridian

import (
	"clove/internals/meridian/fanout"
	AppReplication "clove/internals/meridian/replication/app-replicatrion"
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

func (mer *Meridian) ReplicateApp() *AppReplication.AppReplication {
	return AppReplication.ReplicateApp()
}

func (mer *Meridian) ReplicateMessage() *MessageReplication.MessageReplication {
	return MessageReplication.ReplicateMessage()
}
