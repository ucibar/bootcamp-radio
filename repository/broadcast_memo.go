package repository

import (
	"github.com/uCibar/bootcamp-radio/stream"
	"sync"
)

type BroadcastRepository struct {
	mu         *sync.RWMutex
	broadcasts map[string]*stream.Broadcast
}

func NewBroadcastRepository() *BroadcastRepository {
	return &BroadcastRepository{mu: &sync.RWMutex{}, broadcasts: make(map[string]*stream.Broadcast)}
}

func (r *BroadcastRepository) Get(id string) (*stream.Broadcast, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.broadcasts[id], nil
}

func (r *BroadcastRepository) Add(b *stream.Broadcast) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.broadcasts[b.ID] = b
	return nil
}

func (r *BroadcastRepository) All() []*stream.Broadcast {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var broadcasts []*stream.Broadcast
	for _, b := range r.broadcasts {
		broadcasts = append(broadcasts, b)
	}

	return broadcasts
}
