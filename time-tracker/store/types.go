package store

import "time"

type TimeRecord struct {
	RecordID uint64        `json:"record_id"`
	UserID   uint64        `json:"user_id"`
	Name     string        `json:"name"`
	Start    time.Time     `json:"start_time"`
	StartLoc string        `json:"start_loc"`
	Stop     time.Time     `json:"stop_time"`
	StopLoc  string        `json:"stop_loc"`
	Duration time.Duration `json:"duration"`
}
