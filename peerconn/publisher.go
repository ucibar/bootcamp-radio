package peerconn

import (
	"errors"
	"github.com/pion/webrtc/v3"
	"io"
	"log"
	"sync"
)

type Publisher struct {
	mu *sync.Mutex

	ID     string
	pc     *webrtc.PeerConnection
	tracks map[string]webrtc.TrackLocal
}

func NewPublisher(id string, pc *webrtc.PeerConnection) *Publisher {
	p := &Publisher{mu: &sync.Mutex{}, ID: id, pc: pc, tracks: make(map[string]webrtc.TrackLocal)}
	p.onTrack()
	return p
}

func (p *Publisher) Tracks() []webrtc.TrackLocal {
	p.mu.Lock()
	defer p.mu.Unlock()

	var tracks []webrtc.TrackLocal
	for _, track := range p.tracks {
		tracks = append(tracks, track)
	}

	return tracks
}

func (p *Publisher) addTrack(lt webrtc.TrackLocal) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.tracks[lt.ID()] = lt
}

func (p *Publisher) removeTrack(lt webrtc.TrackLocal) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.tracks, lt.ID())
}

func (p *Publisher) onTrack() {
	p.pc.OnTrack(func(rt *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		if rt.Kind() != webrtc.RTPCodecTypeAudio {
			return
		}

		log.Printf("new remote track: Publisher= %s, Track=%s\n", p.ID, rt.ID())

		lt, err := webrtc.NewTrackLocalStaticRTP(rt.Codec().RTPCodecCapability, rt.ID(), rt.StreamID())
		if err != nil {
			log.Printf("new local track error: Publisher= %s, Error: %v\n", p.ID, err)
			return
		}

		p.addTrack(lt)

		buffer := make([]byte, 1024)
		for {
			i, _, readErr := rt.Read(buffer)
			if readErr != nil {
				log.Printf("remote track read error: Publisher= %s, Track: %s, Error: %v\n", p.ID, rt.ID(), readErr)
				p.removeTrack(lt)
				return
			}

			if _, writeErr := lt.Write(buffer[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
				log.Printf("local track read error: Publisher= %s, Track: %s, Error: %v\n", p.ID, lt.ID(), writeErr)
				p.removeTrack(lt)
				return
			}
		}
	})
}
