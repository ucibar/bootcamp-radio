package stream

import "github.com/pion/webrtc/v3"

type Publisher interface {
	Tracks() []webrtc.TrackLocal
}

type Subscriber interface {
	AddTrack(webrtc.TrackLocal) error
}
