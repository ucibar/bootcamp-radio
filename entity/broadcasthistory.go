package entity

import "time"

type BroadcastHistory struct {
	ID         int64
	UserID     int64
	Title      string
	MaxViewers int64
	BeginAt    time.Time
	EndAt      time.Time
}
