package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/PentoHQ/tech-challenge-time/gateway/config"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var (
	responseTimeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gateway_response_time",
			Help:    "histogram of response times for gateway http handlers",
			Buckets: prometheus.ExponentialBuckets(0.5e-3, 2, 14), // 0.5ms to 4s
		},
		[]string{"path", "status_code"},
	)
)

func init() {
	prometheus.MustRegister(responseTimeHistogram)
}

func newGatewayHandler(ctx context.Context, cfg *config.Config, logger zerolog.Logger) (http.Handler, error) {
	// initialize middleware common to all handlers
	var mw []Middleware
	mw = append(mw, NewRecoverHandler())
	mw = append(mw, NewContextLog(logger)...)
	// we measure response time for all handlers
	mc := NewMetricsConfig().WithTimeHist(responseTimeHistogram)
	mw = append(mw, NewMetricsHandler(mc))

	router := mux.NewRouter()
	for _, url := range cfg.URLs {
		h, err := newHandler(ctx, url, logger)
		if err != nil {
			return nil, err
		}
		// relies on valid URL configuration; does not support query params
		router.Handle(url.Path, Use(h, mw...)).Methods(url.Method)
	}

	// fixme, support queries in route configuration
	// note, this is a temporary hack since the current implmentation
	// of config handling does not support query parameters. I would
	// not deploy this in a production environment but don't have time
	// to fix it during this coding challenge.
	const (
		DAY   = "day"
		WEEK  = "week"
		MONTH = "month"
	)
	u := config.URL{
		HTTP: config.HTTPConf{
			Host: "time-tracker:8080", // don't hardcode host value!
		},
	}
	h, err := newHandler(ctx, u, logger)
	if err != nil {
		return nil, err
	}
	router.Handle("/records", Use(h, mw...)).
		Methods("GET", "OPTIONS").
		Queries("user_id", "{id:[0-9]+}").
		Queries("tz", "{tz:[A-Za-z]+/[A-Za-z]+}").
		Queries("ts", "{ts:[0-9]+}").
		Queries("period", fmt.Sprintf("{period:(?:%s|%s|%s)}", DAY, WEEK, MONTH))

	router.Handle("/ready", &ReadinessHandler{})
	return router, nil
}

func newHandler(ctx context.Context, u config.URL, logger zerolog.Logger) (http.Handler, error) {
	p, err := u.Protocol()
	if err != nil {
		return nil, err
	}
	switch p {
	case config.HTTP:
		// in a real world scenario we would factor this out to perform more
		// sophisticated operations like rewriting headers for HTTPS connections.
		// we ignore Transfer-Encoding hop-by-hop header; expecting `chunked` to
		// be applied if required. returns http.StatusBadGateway if backend is
		// not reachable.
		// TODO: add circuit-breaker
		return sameHost(httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   u.HTTP.Host,
		})), nil
	default:
		return nil, fmt.Errorf("no handler found for %s", p)
	}
}

func sameHost(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = r.URL.Host
		handler.ServeHTTP(w, r)
	})
}

func addCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		handler.ServeHTTP(w, r)
	})
}

type Error struct {
	Err      string         `json:"error"`
	Response *http.Response `json:"-"` // Will not be marshalled
}

func (e Error) Error() string {
	if e.Response == nil {
		return e.Err
	}
	return fmt.Sprintf("%v %v: %d %v",
		e.Response.Request.Method,
		e.Response.Request.URL,
		e.Response.StatusCode,
		e.Err)
}

// HTTP errors
var (
	errInternal   = errors.New("internal_error")
	errBadRequest = errors.New("bad_request")
)

// WriteError writes an error to the http response in JSON format.
func WriteError(w http.ResponseWriter, r *http.Request, err error, code int) {
	// Prepare log.
	logger := LoggerFromRequest(r).With().
		Err(err).
		Int("status", code).
		Logger()

	// Hide error from client if it's internal.
	switch code {
	case http.StatusInternalServerError:
		logger.Error().Msg("unexpected http error")
		err = errInternal
	case http.StatusBadRequest:
		logger.Error().Msg("http error bad request")
		err = errBadRequest
	default:
		logger.Debug().Msg("http error")
	}
	EncodeJSON(w, r, &Error{Err: err.Error()}, code)
}

// EncodeJSON encodes v to w in JSON format.
func EncodeJSON(w http.ResponseWriter, r *http.Request, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		LoggerFromRequest(r).Error().Err(err).Interface("value", v).Msg("failed to encode value to http response")
	}
}

func LoggerFromRequest(r *http.Request) *zerolog.Logger {
	logger := hlog.FromRequest(r).With().
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Logger()
	return &logger
}
