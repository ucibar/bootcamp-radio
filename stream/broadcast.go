package stream

import (
	"sync"
)

type Broadcast struct {
	mu *sync.Mutex

	ID          string
	publisher   Publisher
	subscribers []Subscriber

	Title    string
	Username string
}

func NewBroadcast(id string, publisher Publisher, title, username string) *Broadcast {
	return &Broadcast{mu: &sync.Mutex{}, ID: id, publisher: publisher, subscribers: make([]Subscriber, 0), Title: title, Username: username}
}

func (b *Broadcast) AddSubscriber(subscriber Subscriber) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers = append(b.subscribers, subscriber)

	for _, track := range b.publisher.Tracks() {
		err := subscriber.AddTrack(track)
		if err != nil {
			return err
		}
	}

	return nil
}
