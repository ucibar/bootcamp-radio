package peerconn

import (
	"fmt"
	"github.com/pion/webrtc/v3"
	"github.com/uCibar/bootcamp-radio/entity"
	"github.com/uCibar/bootcamp-radio/stream"
	"log"
	"sync"
)

type Session struct {
	mu           *sync.RWMutex
	ID           string
	Title        string
	owner        *Peer
	participants []*Peer

	stream *stream.Stream
}

func NewSession(ID string, title string, owner *Peer) *Session {
	s := stream.NewStream(nil)
	owner.OnTrack(func(t webrtc.TrackLocal) {
		s.AddTrack(t)
	})

	session := &Session{mu: &sync.RWMutex{}, ID: ID, Title: title, owner: owner, participants: make([]*Peer, 0), stream: s}
	return session
}

func (s *Session) Owner() *entity.User {
	return s.owner.User()
}

func (s *Session) AddParticipant(p *Peer) error {
	s.mu.Lock()
	s.participants = append(s.participants, p)
	s.mu.Unlock()

	err := s.stream.AddSubscriber(p)
	if err != nil {
		return nil
	}

	err = s.messageToEveryone([]byte(fmt.Sprintf("%s joined", p.User().Username)))
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (s *Session) messageToEveryone(msg []byte) error {
	err := s.messageToOwner(msg)
	if err != nil {
		return err
	}

	s.messageToParticipants(msg)
	return nil
}

func (s *Session) messageToOwner(msg []byte) error {
	return s.owner.SendMessage(msg)
}

func (s *Session) messageToParticipants(msg []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, peer := range s.participants {
		if err := peer.SendMessage(msg); err != nil {
			log.Println(err)
		}
	}
}
