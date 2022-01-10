package peerconn

import (
	"github.com/pion/webrtc/v3"
	"github.com/uCibar/bootcamp-radio/entity"
)

var Config = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
}

func NewPublisherPeer(id string, u *entity.User) (*Peer, error) {
	peer, err := NewPeer(id, u)
	if err != nil {
		return nil, err
	}

	_, err = peer.pc.AddTransceiverFromKind(
		webrtc.RTPCodecTypeAudio,
		webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		},
	)
	if err != nil {
		return nil, err
	}

	peer.onTrack()

	return peer, nil
}

func NewSubscriberPeer(id string, u *entity.User) (*Peer, error) {
	peer, err := NewPeer(id, u)
	if err != nil {
		return nil, err
	}

	return peer, nil
}
