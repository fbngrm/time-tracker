package store

import (
	"context"
	"time"

	"github.com/PentoHQ/tech-challenge-time/time-tracker/database"
)

// TimeRecordStore table of the encapsulated datastore.
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
  INSERT INTO players(user_id,name,start_time,start_time_loc,stop_time,stop_time_loc)
  VALUES($1,$2,$3,$4,$5,$5)
  RETURNING *
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
		r.StopLoc).
		Scan(
			&tr.RecordID,
			&tr.UserID,
			&tr.Name,
			&tr.Start,
			&tr.StartLoc,
			&tr.Stop,
			&tr.StopLoc)
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
	start_time,
	start_time_loc,
	stop_time,
	stop_time_loc
  FROM time_records
  WHERE user_id = $1
  AND
  stop_time >= $2;
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
	for rows.Next() {
		if err := rows.Scan(
			&id,
			&name,
			&start,
			&startLoc,
			&stop,
			&stopLoc,
		); err != nil {
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
		}
		recs = append(recs, rec)
	}

	return recs, rows.Err()
}
