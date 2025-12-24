// this module tracks the number of open connections to this server
// this will be used as a part of the heartbeat module used for load balancing connections between servers
package dogpile

import (
	"sync"
)

type DogPile struct {
	Count uint64
	Mux   sync.Mutex
}

func (d *DogPile) Increse() {
	d.Mux.Lock()
	defer d.Mux.Unlock()
	d.Count = d.Count + 1
}
func (d *DogPile) Decrese() {
	d.Mux.Lock()
	defer d.Mux.Unlock()
	d.Count = d.Count - 1
}

func (d *DogPile) GetNumberOfConnections() uint64 {
	return d.Count
}

func New() DogPile {
	return DogPile{
		Count: 0,
		Mux:   sync.Mutex{},
	}
}
