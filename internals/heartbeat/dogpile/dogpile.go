// this module tracks the number of open connections to this server
// this will be used as a part of the heartbeat module used for load balancing connections between servers
package dogpile

import (
	"sync"
)

type DogPile struct {
	Count uint64
	mux   sync.Mutex
}

func (d *DogPile) Increase() {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.Count = d.Count + 1
}
func (d *DogPile) Decrease() {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.Count = d.Count - 1
}
func (d *DogPile) GetCurrentCount() uint64 {
	return d.Count
}

// New creates a DogPile with Count set to 0 and an unlocked mutex.
func New() DogPile {
	return DogPile{
		Count: 0,
		mux:   sync.Mutex{},
	}

}
