package store

import (
	"encoding/json"
	"fmt"
	"time"
)

// User input must always be a timestamp of seconds since UNIX epoch and a
// tz-database location name. We do not trust the user's timezone database and
// convert provided timestamps to the user's location before operating on them.
// This is due to the possiblity that a user's timezone database is not up-to-date
// and thus delievers dates with a wrong UTC-offset.
//
// E.g. if the European Union decides to abondon daylight saving time, the databse
// on the user's system may not be updated in time.
// By doing all time conversions and zone calculations in a central location,
// we can control the timezone database and guarantee correct conversions.
//
// For more information about the tz-database/zoneinfo on UNIX systems, see:
// https://www.iana.org/time-zones

// TimeRecord is a timezone aware representation of a time record.
// In other words, it contains the start and stop time in an UTC-offset aware
// format after conversion from the user's input.
type TimeRecord struct {
	RecordID uint64
	UserID   uint64
	Name     string
	Start    time.Time // time in the user's location
	StartLoc string
	Stop     time.Time // time in the user's location
	StopLoc  string
	Duration int64
}

// TimeStamp is a timezone naive representation of a time record.
// In other words, it contains the start and stop time in an UTC-offset naive
// format (UNIX timestamp) from the user's input before conversion to the user's
// local time.
type TimeStamp struct {
	RecordID uint64 `json:"record_id"`
	UserID   uint64 `json:"user_id"`
	Name     string `json:"name"`
	Start    int64  `json:"start_time"` // seconds since UNIX epoch
	StartLoc string `json:"start_loc"`
	Stop     int64  `json:"stop_time"` // seconds since UNIX epoch
	StopLoc  string `json:"stop_loc"`
}

// UnmarshalJSON unmarshals an offset naive timestamp with start and stop time
// as UNIX timestamps to an offset aware time record with the start and stop
// time in the user's location.
func (tr *TimeRecord) UnmarshalJSON(data []byte) error {
	var ts TimeStamp
	if err := json.Unmarshal(data, &ts); err != nil {
		return err
	}

	// get the start time in the users location
	loc, err := time.LoadLocation(ts.StartLoc)
	if err != nil {
		return err
	}
	startInLoc := time.Unix(ts.Start, 0).In(loc)

	// get the stop time in the users location
	loc, err = time.LoadLocation(ts.StopLoc)
	if err != nil {
		return err
	}
	stopInLoc := time.Unix(ts.Stop, 0).In(loc)

	tr.UserID = ts.UserID
	tr.Name = ts.Name
	tr.Start = startInLoc
	tr.StartLoc = ts.StartLoc
	tr.Stop = stopInLoc
	tr.StopLoc = ts.StopLoc
	tr.Duration = ts.Stop - ts.Start

	return nil
}

// MarshalJSON formats the dates and duration.
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
		Start:    tr.Start.Format("02 Jan 2006 15:04:05"),
		StartLoc: tr.StartLoc,
		Stop:     tr.Stop.Format("02 Jan 2006 15:04:05"),
		StopLoc:  tr.StopLoc,
		Duration: formatDuration(time.Second * time.Duration(tr.Duration)),
	}
	return json.Marshal(t)
}

func formatDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
