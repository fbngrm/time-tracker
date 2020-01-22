package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PentoHQ/tech-challenge-time/time-tracker/api"
	"github.com/PentoHQ/tech-challenge-time/time-tracker/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// HTTP errors
var (
	errInternal   = errors.New("internal_error")
	errNotFound   = errors.New("not_found")
	errBadRequest = errors.New("bad_request")
)

const (
	DAY   = "day"
	WEEK  = "week"
	MONTH = "month"
)

// newHandler creates an http handler that operates on time records.
func newHandler(ts timeRecordStore, timeout time.Duration, logger zerolog.Logger) (http.Handler, error) {
	var mw []middleware.Middleware
	mw = append(mw, middleware.NewRecoverHandler())
	mw = append(mw, middleware.NewContextLog(logger)...)
	mw = append(mw, middleware.NewCORSHandler())

	// services handle http requests and hold a store to operate on a database
	recordSrvc := middleware.Use(&timeRecordService{ts, timeout}, mw...)

	router := mux.NewRouter()
	router.Handle("/ready", &readinessHandler{}).Methods("GET")

	// time record store
	router.Handle("/record", recordSrvc).Methods("POST", "OPTIONS")
	router.Handle("/records", recordSrvc).
		Methods("GET").
		Queries("user_id", "{id:[0-9]+}").
		Queries("tz", "{tz:[A-Za-z]+/[A-Za-z]+}").
		Queries("ts", "{ts:[0-9]+}").
		Queries("period", fmt.Sprintf("{period:(?:%s|%s|%s)}", DAY, WEEK, MONTH))

	return router, nil
}

// encodeJSON encodes v to w in JSON format.
func encodeJSON(w http.ResponseWriter, r *http.Request, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		loggerFromRequest(r).Error().Err(err).Interface("value", v).Msg("failed to encode value to http response")
	}
}

func loggerFromRequest(r *http.Request) *zerolog.Logger {
	logger := hlog.FromRequest(r).With().
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Logger()
	return &logger
}

// writeError writes an error to the http response in JSON format.
func writeError(w http.ResponseWriter, r *http.Request, err error, code int) {
	// prepare log
	logger := loggerFromRequest(r).With().
		Err(err).
		Int("status", code).
		Logger()
	// hide error from client if it's internal
	if code == http.StatusInternalServerError {
		logger.Error().Msg("unexpected http error")
		err = errInternal
	} else if code == http.StatusBadRequest {
		logger.Error().Msg("bad request")
		err = errBadRequest
	} else {
		logger.Debug().Msg("http error")
	}
	encodeJSON(w, r, &api.Error{Err: err.Error()}, code)
}
