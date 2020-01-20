package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PentoHQ/tech-challenge-time/gateway/config"
	"github.com/PentoHQ/tech-challenge-time/gateway/metrics"
	"github.com/PentoHQ/tech-challenge-time/gateway/server"
	"github.com/rs/zerolog"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "unkown"

	cfgPath     = kingpin.Flag("cfg-file", "path to config file").Envar("CFG_PATH").Required().String()
	httpAddr    = kingpin.Flag("http-addr", "address of HTTP server").Envar("HTTP_ADDR").Required().String()
	metricsAddr = kingpin.Flag("metrics-addr", "address of metrics server").Envar("METRICS_ADDR").Required().String()
	service     = kingpin.Flag("service", "service name").Envar("SERVICE").Default("gateway").String()

	// should be greater than prometheus scrape interval (default 30s); decreased in coding challenge
	shutdownDelay = kingpin.Flag("shutdown-delay", "shutdown delay in ms").Envar("SHUTDOWN_DELAY").Default("5000").Int()
)

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit
		cancel()
	}()

	cfg, err := config.FromFile(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", *service, err)
		os.Exit(2)
	}
	logger := newLogger(*service, version)
	httpSrv, err := server.New(ctx, *httpAddr, cfg, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", *service, err)
		os.Exit(2)
	}

	runServer(ctx, httpSrv, metrics.New(*metricsAddr, logger), *shutdownDelay)
}

// Using default log level debug and write to stderr.
// Note, we log in (inefficient) human friendly format to console here since it
// is a coding challenge. In a production environment we would prefer structured,
// machine parsable format. So we could make use of automated log analysis e.g.
// error reporting.
func newLogger(service, version string) zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// replace standard log
	log.SetFlags(0)
	log.SetOutput(logger)
	return logger.With().
		Interface("service", service).
		Interface("version", version).
		Logger()
}

func runServer(ctx context.Context, httpSrv *server.HTTPServer, metricsSrv *metrics.MetricsServer, shutdownDelay int) {
	go httpSrv.Run()
	go metricsSrv.Run()

	<-ctx.Done()

	server.HealthCheckShutDown()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownDelay)*time.Millisecond)
	defer cancel()

	// when shutting down, we first gracefully shutting down the main http
	// server, waiting for it to finish processing all the running requests,
	// then we shut down the metrics server, which includes waiting for
	// prometheus to scrape the metrics one more time, to avoid loosing any data.
	httpSrv.Shutdown(ctx)
	metricsSrv.Shutdown(ctx)
}
