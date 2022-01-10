package stream

import (
	"github.com/pion/webrtc/v3"
	"sync"
)

type Stream struct {
	mu *sync.RWMutex

	ID string

	tracks      map[string]webrtc.TrackLocal
	subscribers map[string]Subscriber
}

func NewStream(tracks []webrtc.TrackLocal) *Stream {
	ts := make(map[string]webrtc.TrackLocal)
	for _, t := range tracks {
		ts[t.ID()] = t
	}

	return &Stream{mu: &sync.RWMutex{}, tracks: ts, subscribers: make(map[string]Subscriber)}
}

func (s *Stream) AddTrack(t webrtc.TrackLocal) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tracks[t.ID()] = t
}

func (s *Stream) GetTracks() []webrtc.TrackLocal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tracks []webrtc.TrackLocal
	for _, track := range s.tracks {
		tracks = append(tracks, track)
	}

	return tracks
}

func (s *Stream) AddSubscriber(subscriber Subscriber) error {
	s.mu.Lock()
	s.subscribers[subscriber.ID()] = subscriber
	s.mu.Unlock()

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, track := range s.tracks {
		err := subscriber.SubscribeTrack(track)
		if err != nil {
			return err
		}
	}

	return nil
}
