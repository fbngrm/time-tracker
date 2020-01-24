package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/PentoHQ/tech-challenge-time/time-tracker/store"
)

var (
	startCreateTest time.Time
	stopCreateTest  time.Time
)

func init() {
}

// uses the user id to get the test data.
type mockTimeRecordStore struct{}

func (rs *mockTimeRecordStore) Create(ctx context.Context, r store.TimeRecord) (*store.TimeRecord, error) {
	var tr store.TimeRecord
	decoder := json.NewDecoder(strings.NewReader(createRecordTests[r.UserID].p))
	decoder.DisallowUnknownFields() // catch unwanted fields
	if err := decoder.Decode(&tr); err != nil {
		return nil, err
	}
	tr.RecordID = r.UserID
	return &tr, createRecordTests[r.UserID].e
}
func (rs *mockTimeRecordStore) Get(ctx context.Context, userID uint64, t time.Time) ([]store.TimeRecord, error) {
	// return getRecordTests[userID].r, getRecordTests[userID].e
	return nil, nil
}

// test cases indexed by id
var createRecordTests = map[uint64]struct {
	d string // description of test case
	e error  // mock store error
	p string // request payload
	u string // route of the test request
	s int    // expected http status code
	b []byte // expected payload
}{
	// errors
	0: {
		d: "expect missing path segment to result in 404",
		u: "/",
		s: http.StatusNotFound,
		b: []byte("404 page not found"),
	},
	1: { // 500
		d: "expect mal-formed JSON payload to result in 500",
		e: errInternal,
		u: "record",
		p: `{"user_id":2`, // missing closing brace
		s: http.StatusInternalServerError,
		b: []byte(fmt.Sprintf(`{"error":"%s"}`, errInternal.Error())),
	},
	2: { // 500
		d: "expect store error to result in 500",
		e: errInternal,
		u: "record",
		p: `{"user_id":2}`, // user_id is the testcase-id used by the mock store
		s: http.StatusInternalServerError,
		b: []byte(fmt.Sprintf(`{"error":"%s"}`, errInternal.Error())),
	},
	// success
	3: {
		d: "expect to successfully create and return a time record",
		u: "record",
		p: `{"user_id":3,"name":"foo","start_time":1577833200,"start_loc":"Europe/Berlin","stop":1577836800,"stop_loc":"Europe/Berlin", "duration":3600}`,
		s: http.StatusOK,
		b: []byte(`{"record_id":3,"user_id":3,"name":"foo","start_time":"01 Jan 2020 00:00:00","start_loc":"Europe/Berlin","stop":"01 Jan 2020 01:00:00","stop_loc":"Europe/Berlin", "duration":"01:00:00"}`),
	},
}

func getDateFromString(t *testing.T, date string) time.Time {
	ti, err := time.Parse(time.RFC3339, date)
	if err != nil {
		t.Fatal(err)
	}
	return ti
}

func TestServeHTTPCreate(t *testing.T) {
	// service initialized with a mock store to
	// control the data and errors we return
	rs := &timeRecordService{
		&mockTimeRecordStore{},
		200 * time.Millisecond,
	}
	// test server
	s := httptest.NewServer(rs)
	defer s.Close()
	c := s.Client()

	for _, tc := range createRecordTests {
		tt := tc
		t.Run(tt.d, func(t *testing.T) {
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", s.URL, tt.u), strings.NewReader(tt.p))
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			resp, err := c.Do(req)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			// expected result
			if want, got := tt.s, resp.StatusCode; want != got {
				t.Errorf("want status code %d got %d", want, got)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			resp.Body.Close()
			if want, got := tt.b, body; bytes.Compare(want, got) == 1 {
				t.Errorf("want response\n%+s\ngot\n%+s", want, got)
			}
		})
	}
}

// // test cases indexed by player id
// var updateTests = map[uint64]struct {
// 	d string        // description of test case
// 	r *store.Player // mock store response
// 	e error         // mock store error
// 	u string        // request url path
// 	p string        // request payload
// 	s int           // expected http status code
// 	b []byte        // expected payload
// }{
// 	// url path errors
// 	0: { // 500
// 		d: "expect malformed JSON payload to result in 500 when updating player",
// 		u: "players/update",
// 		p: `{"player_id":1`,
// 		s: http.StatusInternalServerError,
// 		b: []byte(fmt.Sprintf(`{"error":"%s"}`, errInternal.Error())),
// 	},
// 	// store errors
// 	1: { // 500
// 		d: "expect store error to result in 500 when updating player",
// 		e: errInternal,
// 		u: "players/update",
// 		p: `{"player_id":1}`,
// 		s: http.StatusInternalServerError,
// 		b: []byte(fmt.Sprintf(`{"error":"%s"}`, errInternal.Error())),
// 	},
// 	// success
// 	2: { // 200
// 		d: "expect player's store to get updated and status get set to benched",
// 		u: "players/update",
// 		r: &store.Player{
// 			PlayerID:  2,
// 			RosterID:  0,
// 			FirstName: "foo",
// 			LastName:  "bar",
// 			Alias:     "foobar",
// 			Status:    "benched",
// 		},
// 		p: `{"player_id":2,"roster_id":1,"first_name":"foo","last_name":"bar","alias":"foobar","status":"active"}`,
// 		s: http.StatusOK,
// 		b: []byte(`{"player_id":2,"roster_id":0,"first_name":"foo","last_name":"bar","alias":"foobar","status":"benched"}`),
// 	},
// }

// func TestUpdate(t *testing.T) {
// 	// service initialized with a mock store to
// 	// control the data and errors we return
// 	ps := &playerService{
// 		&mockPlayerStore{},
// 		200 * time.Millisecond,
// 	}

// 	router := mux.NewRouter()
// 	router.Handle("/players/update", ps).Methods("PATCH")

// 	s := httptest.NewServer(router)
// 	defer s.Close()
// 	c := s.Client()

// 	for _, tc := range updateTests {
// 		tt := tc
// 		t.Run(tt.d, func(t *testing.T) {
// 			req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", s.URL, tt.u), strings.NewReader(tt.p))
// 			if err != nil {
// 				t.Fatalf("unexpected err: %v", err)
// 			}
// 			resp, err := c.Do(req)
// 			if err != nil {
// 				t.Fatalf("unexpected err: %v", err)
// 			}
// 			// expected result
// 			if want, got := tt.s, resp.StatusCode; want != got {
// 				t.Errorf("want status code %d got %d", want, got)
// 			}
// 			body, err := ioutil.ReadAll(resp.Body)
// 			if err != nil {
// 				t.Fatalf("unexpected err: %v", err)
// 			}
// 			resp.Body.Close()
// 			if want, got := tt.b, body; bytes.Compare(want, got) == 1 {
// 				t.Errorf("want response\n%+s\ngot\n%+s", want, got)
// 			}
// 		})
// 	}
// }

// // test cases indexed by player id
// var changeTests = map[uint64]struct {
// 	d string              // description of test case
// 	r *store.PlayerChange // mock store response
// 	e error               // mock store error
// 	u string              // request url path
// 	p string              // request payload
// 	s int                 // expected HTTP status code
// 	b []byte              // expected payload
// }{
// 	// url path errors
// 	0: { // 500
// 		d: "expect malformed JSON payload to result in 500 when updating player",
// 		u: "players/change",
// 		p: `{"player_id":1`,
// 		s: http.StatusInternalServerError,
// 		b: []byte(fmt.Sprintf(`{"error":"%s"}`, errInternal.Error())),
// 	},
// 	// store errors
// 	1: { // 500
// 		d: "expect store error to result in 500 when updating player",
// 		e: errInternal,
// 		u: "players/change",
// 		s: http.StatusInternalServerError,
// 		b: []byte(fmt.Sprintf(`{"error":"%s"}`, errInternal.Error())),
// 	},
// 	// success
// 	2: { // 200
// 		d: "expect players statuses to get swapped",
// 		u: "players/change",
// 		r: &store.PlayerChange{
// 			Active: store.Player{
// 				PlayerID: 2,
// 				RosterID: 0,
// 				Status:   "benched",
// 			},
// 			Benched: store.Player{
// 				PlayerID: 3,
// 				RosterID: 0,
// 				Status:   "active",
// 			},
// 		},
// 		p: `{"active":{"player_id":2,"roster_id":0,"status":"active"},"benched":{"player_id":3,"roster_id":0,"status":"benched"}}`,
// 		s: http.StatusOK,
// 		b: []byte(`{"active":{"player_id":2,"roster_id":0,"first_name":"","last_name":"","alias":"","status":"benched"},"benched":{"player_id":3,"roster_id":0,"first_name":"","last_name":"","alias":"","status":"acti`),
// 	},
// }

// func TestChange(t *testing.T) {
// 	// service initialized with a mock store to
// 	// control the data and errors we return
// 	ps := &playerService{
// 		&mockPlayerStore{},
// 		200 * time.Millisecond,
// 	}

// 	router := mux.NewRouter()
// 	router.Handle("/players/change", ps).Methods("PATCH")

// 	s := httptest.NewServer(router)
// 	defer s.Close()
// 	c := s.Client()

// 	for _, tc := range changeTests {
// 		tt := tc
// 		t.Run(tt.d, func(t *testing.T) {
// 			req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", s.URL, tt.u), strings.NewReader(tt.p))
// 			if err != nil {
// 				t.Fatalf("unexpected err: %v", err)
// 			}
// 			resp, err := c.Do(req)
// 			if err != nil {
// 				t.Fatalf("unexpected err: %v", err)
// 			}
// 			// expected result
// 			if want, got := tt.s, resp.StatusCode; want != got {
// 				t.Errorf("want status code %d got %d", want, got)
// 			}
// 			body, err := ioutil.ReadAll(resp.Body)
// 			if err != nil {
// 				t.Fatalf("unexpected err: %v", err)
// 			}
// 			resp.Body.Close()
// 			if want, got := tt.b, body; bytes.Compare(want, got) == 1 {
// 				t.Errorf("want response\n%+s\ngot\n%+s", want, got)
// 			}
// 		})
// 	}
// }
