package store

import (
	"context"
	"time"

	"github.com/fgrimme/time-tracker/time-tracker/database"
)

type TimeRecordStore struct {
	db *database.DB
}

func New(db *database.DB) *TimeRecordStore {
	return &TimeRecordStore{
		db: db,
	}
}

// Create inserts a new time record to the datastore. The record id is not inserted
// and must be created by the datastore. Returns the newly created record with
// the generated id.
func (ts *TimeRecordStore) Create(ctx context.Context, r TimeRecord) (*TimeRecord, error) {
	query := `
  INSERT INTO time_records(
    user_id,
	name,
	start_time,
	start_time_loc,
	stop_time,
	stop_time_loc,
	duration)
  VALUES($1,$2,$3,$4,$5,$6,$7)
  RETURNING
    id,
    user_id,
	name,
	start_time AT TIME ZONE start_time_loc,
	start_time_loc,
	stop_time AT TIME ZONE stop_time_loc,
	stop_time_loc,
	duration
  `
	db := ts.db.GetDB()
	ctx, cancel := ts.db.RequestContext(ctx)
	defer cancel()

	var tr TimeRecord
	err := db.QueryRowContext(ctx, query,
		r.UserID,
		r.Name,
		r.Start,
		r.StartLoc,
		r.Stop,
		r.StopLoc,
		r.Duration).
		Scan(
			&tr.RecordID,
			&tr.UserID,
			&tr.Name,
			&tr.Start,
			&tr.StartLoc,
			&tr.Stop,
			&tr.StopLoc,
			&tr.Duration)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (ts *TimeRecordStore) Get(ctx context.Context, userID uint64, t time.Time) ([]TimeRecord, error) {
	query := `
  SELECT
  	id,
	name,
	start_time AT TIME ZONE tr.start_time_loc,
	start_time_loc,
	stop_time AT TIME ZONE tr.stop_time_loc,
	stop_time_loc,
	duration
  FROM time_records
  AS tr
  WHERE tr.user_id = $1
  AND
  tr.stop_time >= $2
  ORDER BY start_time DESC;
  `

	db := ts.db.GetDB()
	ctx, cancel := ts.db.RequestContext(ctx)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, userID, t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recs := make([]TimeRecord, 0)

	var id uint64
	var name string
	var start, stop time.Time
	var startLoc, stopLoc string
	var duration int64
	for rows.Next() {
		if err := rows.Scan(
			&id,
			&name,
			&start,
			&startLoc,
			&stop,
			&stopLoc,
			&duration); err != nil {
			return nil, err
		}
		rec := TimeRecord{
			RecordID: id,
			UserID:   userID,
			Name:     name,
			Start:    start,
			StartLoc: startLoc,
			Stop:     stop,
			StopLoc:  stopLoc,
			Duration: duration,
		}
		recs = append(recs, rec)
	}

	return recs, rows.Err()
}
