package peerconn

import (
	"github.com/pion/webrtc/v3"
)

var Config = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
}

func NewPublisherPeer() (*webrtc.PeerConnection, error) {
	pc, err := webrtc.NewPeerConnection(Config)
	if err != nil {
		return nil, err
	}

	_, err = pc.AddTransceiverFromKind(
		webrtc.RTPCodecTypeAudio,
		webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly},
	)
	if err != nil {
		return nil, err
	}

	return pc, nil
}

func NewSubscriberPeer() (*webrtc.PeerConnection, error) {
	pc, err := webrtc.NewPeerConnection(Config)
	if err != nil {
		return nil, err
	}

	return pc, nil
}
