package peerconn

import (
	"github.com/pion/webrtc/v3"
)

type Subscriber struct {
	ID string
	pc *webrtc.PeerConnection
}

func NewSubscriber(id string, pc *webrtc.PeerConnection) *Subscriber {
	p := &Subscriber{ID: id, pc: pc}
	return p
}

func (s *Subscriber) AddTrack(track webrtc.TrackLocal) error {
	_, err := s.pc.AddTrack(track)
	if err != nil {
		return err
	}

	return nil
}
