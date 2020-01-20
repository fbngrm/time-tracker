package server

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/PentoHQ/tech-challenge-time/gateway/config"
	"github.com/rs/zerolog"
)

// urls for test-handlers
var gatewayConf = config.Config{
	URLs: []config.URL{
		config.URL{ // proxy
			Path:   "/drivers/{id:[0-9]+}",
			Method: "GET",
			HTTP: config.HTTPConf{
				Host: "zombie-driver", // will be overwritten by test server host and port
			},
		},
	},
}

// test data by driver-ID
var gatewayTests = map[string]struct {
	d string // description of test case
	z string // response of zombie-driver srevice mock
	p string // request path
	r string // expected response data
	s int    // expected response status code
}{
	"0": {
		d: "expect StatusBadGateway when failing to reach the backend",
		p: "/drivers/0", // 0 => test hijacked requests; backend not reachable
		s: http.StatusBadGateway,
	},
	"1": {
		d: "expect successful proxying; #1",
		z: `{"id":1,"zombie":true}`,
		p: "/drivers/1",
		r: `{"id":1,"zombie":true}`,
		s: http.StatusOK,
	},
	"2": {
		d: "expect successful proxying; #2",
		z: `{"id":2,"zombie":false}`,
		p: "/drivers/2",
		r: `{"id":2,"zombie":false}`,
		s: http.StatusOK,
	},
	"3": {
		d: "expect StatusNotFound for invalid URL",
		p: "/drivers",
		r: "404 page not found",
		s: http.StatusNotFound,
	},
	"4": {
		d: "expect StatusNotFound for unknown driver",
		p: "/drivers/404",
		s: http.StatusNotFound,
	},
}

func TestProxy(t *testing.T) {
	// mute logger in tests
	logger := zerolog.New(ioutil.Discard)
	log.SetFlags(0)
	log.SetOutput(logger)

	// create a mock zombie-service for the reverse proxy
	zombieService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// driver id
		segments := strings.Split(r.URL.Path, "/")
		if len(segments) != 3 {
			t.Fatalf("expect 3 path segments but got %d", len(segments))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id := segments[2]

		// backend not reachable
		if id == "0" {
			c, _, _ := w.(http.Hijacker).Hijack()
			c.Close()
			return
		}

		// we ignore other hop-by-hop headers for now
		if r.Header.Get("X-Forwarded-For") == "" {
			t.Error("expect X-Forwarded-For header")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// send mock data
		if p, ok := gatewayTests[id]; ok {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(p.z))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
				return
			}
			return
		}
		// driver ID unknown
		w.WriteHeader(http.StatusNotFound)
	}))
	defer zombieService.Close()

	// Note, we need to overwrite the URL host in the test config with the
	// address of the test zombieService. The httptest.Server uses a local
	// Listener initialized to listen on a random port. Using a custom Listener
	// and providing a port would require supporting `serveFlag` and IPv6.
	// For more info see:
	// https://golang.org/src/net/http/httptest/server.go?s=477:1449#L72
	u, err := url.Parse(zombieService.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gatewayConf.URLs[1].HTTP.Host = u.Host

	// handler to test
	h, err := newGatewayHandler(context.Background(), &gatewayConf, logger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// serve test handler
	gatewayService := httptest.NewServer(h)
	defer gatewayService.Close()
	gatewayClient := gatewayService.Client()

	t.Run("zombie-service", func(t *testing.T) {
		for id := range gatewayTests {
			tt := gatewayTests[id]
			t.Run(tt.d, func(t *testing.T) {
				t.Parallel()
				req, err := http.NewRequest("GET", gatewayService.URL+tt.p, nil)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				req.Close = true
				req.Header.Set("Connection", "close")

				res, err := gatewayClient.Do(req)
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}

				if w, g := tt.s, res.StatusCode; w != g {
					t.Errorf("want status code %d got %d", w, g)
				}
				data, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("failed to read response %v", err)
				}
				if w, g := tt.r, strings.TrimSpace(string(data)); w != g {
					t.Errorf("want response %s got %s", w, g)
				}
			})
		}
	})
}
