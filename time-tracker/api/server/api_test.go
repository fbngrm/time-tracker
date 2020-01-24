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
	return make([]store.TimeRecord, 1), getRecordTests[userID].e
}

// test cases indexed by user id
var createRecordTests = map[uint64]struct {
	d string // description of test case
	e error  // mock store error
	u string // route of test request
	p string // request payload
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

type params struct {
	u  string // user id
	ts string // timestamp
	tz string // timezone
	p  string // period
}

// var mockResults = map[string]map[string][]store.TimeRecord{
// 	"UTC":               {"day": make([]store.TimeRecord, 1)},
// 	"Europe/Berlin":     {"week": make([]store.TimeRecord, 1)},
// 	"Europe/Copenhagen": {"month": make([]store.TimeRecord, 1)},
// }

// test cases indexed by user id
var getRecordTests = map[uint64]struct {
	d string             // description of test case
	r []store.TimeRecord // mock store response
	e error              // mock store error
	s int                // expected http status code
	b []byte             // expected payload
	p params             // request params
}{
	// errors
	0: { // 400
		d: "expect missing user id to result in 400",
		s: http.StatusBadRequest,
		b: []byte(fmt.Sprintf(`{"error":"%s"}`, errBadRequest.Error())),
	},
	1: { // 400
		d: "expect missing timestamp to result in 400",
		s: http.StatusBadRequest,
		b: []byte(fmt.Sprintf(`{"error":"%s"}`, errBadRequest.Error())),
		p: params{u: "1"},
	},
	2: { // 500
		d: "expect wrong timezone to result in 500",
		s: http.StatusInternalServerError,
		b: []byte(fmt.Sprintf(`{"error":"%s"}`, errInternal.Error())),
		p: params{u: "1", ts: "invalid"},
	},
	// success
	3: { // 200
		d: "expect successfull request",
		s: http.StatusOK,
		p: params{u: "1", ts: "0"}, // timzone and location can be empty
		r: make([]store.TimeRecord, 1),
	},
}

func TestServeHTTPGet(t *testing.T) {
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

	for _, tc := range getRecordTests {
		tt := tc
		t.Run(tt.d, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/records", s.URL), nil)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			// set query string
			q := req.URL.Query()
			q.Add("user_id", tt.p.u)
			q.Add("ts", tt.p.ts)
			q.Add("tz", tt.p.tz)
			q.Add("period", tt.p.p)
			req.URL.RawQuery = q.Encode()

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
