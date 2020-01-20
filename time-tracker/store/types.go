package store

import (
	"encoding/json"
	"fmt"
	"time"
)

type TimeRecord struct {
	RecordID uint64    `json:"record_id"`
	UserID   uint64    `json:"user_id"`
	Name     string    `json:"name"`
	Start    time.Time `json:"start_time"`
	StartLoc string    `json:"start_loc"`
	Stop     time.Time `json:"stop_time"`
	StopLoc  string    `json:"stop_loc"`
}

// we support unmarshaling of timestamps
func (tr *TimeRecord) UnmarshalJSON(data []byte) error {
	var v = struct {
		RecordID uint64 `json:"record_id"`
		UserID   uint64 `json:"user_id"`
		Name     string `json:"name"`
		Start    int64  `json:"start_time"`
		StartLoc string `json:"start_loc"`
		Stop     int64  `json:"stop_time"`
		StopLoc  string `json:"stop_loc"`
	}{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	tr.RecordID = v.RecordID
	tr.UserID = v.UserID
	tr.Name = v.Name
	tr.Start = time.Unix(v.Start, 0)
	tr.StartLoc = v.StartLoc
	tr.Stop = time.Unix(v.Stop, 0)
	tr.StopLoc = v.StopLoc

	return nil
}

// we want to format the dates and duration and thus need to crate a custom marshal function.
func (tr *TimeRecord) MarshalJSON() ([]byte, error) {
	t := struct {
		RecordID uint64 `json:"record_id"`
		UserID   uint64 `json:"user_id"`
		Name     string `json:"name"`
		Start    string `json:"start_time"`
		StartLoc string `json:"start_loc"`
		Stop     string `json:"stop_time"`
		StopLoc  string `json:"stop_loc"`
		Duration string `json:"duration"`
	}{
		RecordID: tr.RecordID,
		UserID:   tr.UserID,
		Name:     tr.Name,
		Start:    tr.Start.Format("02 Jan 2006 15:04"),
		StartLoc: tr.StartLoc,
		Stop:     tr.Stop.Format("02 Jan 2006 15:04"),
		StopLoc:  tr.StopLoc,
		Duration: formatDuration(tr.Stop.Sub(tr.Start)),
	}
	return json.Marshal(t)
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
