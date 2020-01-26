package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/PentoHQ/tech-challenge-time/time-tracker/store"
)

// timeRecordStore handles operations on time records.
type timeRecordStore interface {
	Create(ctx context.Context, r store.TimeRecord) (*store.TimeRecord, error)
	Get(ctx context.Context, userID uint64, t time.Time) ([]store.TimeRecord, error)
}

// recordService provides API methods to operate on time records.
type timeRecordService struct {
	timeRecordStore
	timeout time.Duration
}

// ServeHTTP serves requests to the time record enpoint.
func (rs *timeRecordService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle CORS preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), rs.timeout)
	defer cancel()
	// we attach the logger from the request to the context so we do not need
	// to pass it as an parameter
	ctx = loggerFromRequest(r).WithContext(ctx)

	_, route := path.Split(r.URL.Path)
	switch route {
	case "record":
		var tr store.TimeRecord
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // catch unwanted fields
		if err := decoder.Decode(&tr); err != nil {
			writeError(w, r, err, http.StatusInternalServerError)
			return
		}
		rs.createRecord(ctx, w, r, tr)
		return

	case "records":
		q := r.URL.Query()

		// get the user id from the requests params
		// if not supplied, we consider the request as malformed
		uid := q.Get("user_id")
		if len(uid) == 0 {
			writeError(w, r, errBadRequest, http.StatusBadRequest)
			return
		}
		userID, err := strconv.ParseUint(uid, 10, 64) // mux validates type
		if err != nil {
			writeError(w, r, errInternal, http.StatusInternalServerError)
			return
		}
		// get the timestamp from the requests params
		// if not supplied, we consider the request as malformed
		ts := q.Get("ts")
		if len(ts) == 0 {
			writeError(w, r, err, http.StatusBadRequest)
			return
		}
		timestamp, err := strconv.ParseInt(ts, 10, 64) // mux validates type
		if err != nil {
			writeError(w, r, errInternal, http.StatusInternalServerError)
			return
		}
		t := time.Unix(timestamp, 0)

		// get the tz-database zone name from the requests params
		// if not supplied, we assume UTC
		zone := q.Get("tz")
		loc, err := time.LoadLocation(zone)
		if err != nil {
			writeError(w, r, err, http.StatusBadRequest)
			return
		}

		// get the time period from the requests params
		// if not supplied, we assume DAY
		period := q.Get("period")
		if len(period) == 0 {
			period = DAY
		}

		rs.getRecords(ctx, w, r, userID, t.In(loc), loc, period)
		return
	}
	writeError(w, r, errNotFound, http.StatusNotFound)
	return
}

func (rs *timeRecordService) createRecord(ctx context.Context, w http.ResponseWriter, r *http.Request, tr store.TimeRecord) {
	rec, err := rs.Create(ctx, tr)
	if err != nil {
		writeError(w, r, err, http.StatusInternalServerError)
		return
	}
	encodeJSON(w, r, rec, http.StatusOK)
}

func (rs *timeRecordService) getRecords(ctx context.Context, w http.ResponseWriter, r *http.Request, userID uint64, t time.Time, loc *time.Location, period string) {
	day, err := getStartOfPeriod(t, loc, period)
	if err != nil {
		writeError(w, r, err, http.StatusInternalServerError)
		return
	}
	recs, err := rs.Get(ctx, userID, day)
	if err != nil {
		writeError(w, r, err, http.StatusInternalServerError)
		return
	}
	encodeJSON(w, r, recs, http.StatusOK)
}

// getStartOfPeriod returns the start or the first day in the given period.
func getStartOfPeriod(t time.Time, loc *time.Location, period string) (time.Time, error) {
	var day time.Time
	switch period {
	case DAY:
		// get the current day in the given location
		currentYear, currentMonth, today := t.Date()
		day = time.Date(currentYear, currentMonth, today, 0, 0, 0, 0, loc)
	case WEEK:
		// get the first day of the week in the given location
		isoYear, isoWeek := t.ISOWeek()
		day = firstDayOfISOWeek(isoYear, isoWeek, loc)
	case MONTH:
		// get the first day of the current month in the given location
		currentYear, currentMonth, _ := t.Date()
		day = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, loc)
	default:
		return time.Now(), fmt.Errorf("unknown period: %s", period)
	}
	return day, nil
}

// firstDayOfISOWeek returns the first day of the given week in the given year
// at the given location.
func firstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()
	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year { // iterate forward to the first day of the first week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoWeek < week { // iterate forward to the first day of the given week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}
