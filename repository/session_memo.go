package repository

import (
	"github.com/uCibar/bootcamp-radio/peerconn"
	"sync"
)

type SessionRepository struct {
	mu      *sync.RWMutex
	streams map[string]*peerconn.Session
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{mu: &sync.RWMutex{}, streams: make(map[string]*peerconn.Session)}
}

func (r *SessionRepository) Get(id string) (*peerconn.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.streams[id], nil
}

func (r *SessionRepository) Add(b *peerconn.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.streams[b.ID] = b
	return nil
}

func (r *SessionRepository) All() []*peerconn.Session {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var broadcasts []*peerconn.Session
	for _, b := range r.streams {
		broadcasts = append(broadcasts, b)
	}

	return broadcasts
}
