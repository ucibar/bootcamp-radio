package peerconn

import (
	"errors"
	"github.com/pion/webrtc/v3"
	"github.com/uCibar/bootcamp-radio/entity"
	"io"
	"log"
)

type Peer struct {
	pc   *webrtc.PeerConnection
	id   string
	user *entity.User
	dc   *webrtc.DataChannel

	onTrackHandler func(track webrtc.TrackLocal)
}

func NewPeer(id string, user *entity.User) (*Peer, error) {
	pc, err := webrtc.NewPeerConnection(Config)
	if err != nil {
		return nil, err
	}

	dc, err := pc.CreateDataChannel("message-channel", nil)
	if err != nil {
		return nil, err
	}

	return &Peer{pc: pc, id: id, user: user, dc: dc}, nil
}

func (p *Peer) ID() string {
	return p.id
}

func (p *Peer) PC() *webrtc.PeerConnection {
	return p.pc
}

func (p *Peer) User() *entity.User {
	return p.user
}

func (p *Peer) SendMessage(msg []byte) error {
	return p.dc.Send(msg)
}

func (p *Peer) OnMessage(cb func([]byte)) {
	p.dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		cb(msg.Data)
	})
}

func (p *Peer) SubscribeTrack(track webrtc.TrackLocal) error {
	_, err := p.pc.AddTrack(track)
	if err != nil {
		return err
	}

	return nil
}

func (p *Peer) OnTrack(cb func(t webrtc.TrackLocal)) {
	p.onTrackHandler = cb
}

func (p *Peer) onTrack() {
	p.pc.OnTrack(func(rt *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		if rt.Kind() != webrtc.RTPCodecTypeAudio {
			return
		}

		log.Printf("new remote track: Publisher= %s, Track=%s\n", p.id, rt.ID())

		lt, err := webrtc.NewTrackLocalStaticRTP(rt.Codec().RTPCodecCapability, rt.ID(), rt.StreamID())
		if err != nil {
			log.Printf("new local track error: Publisher= %s, Error: %v\n", p.id, err)
			return
		}

		go p.onTrackHandler(lt)

		buffer := make([]byte, 1024)
		for {
			i, _, readErr := rt.Read(buffer)
			if readErr != nil {
				log.Printf("remote track read error: Publisher= %s, Track: %s, Error: %v\n", p.id, rt.ID(), readErr)
				return
			}

			if _, writeErr := lt.Write(buffer[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
				log.Printf("local track read error: Publisher= %s, Track: %s, Error: %v\n", p.id, lt.ID(), writeErr)
				return
			}
		}
	})
}
