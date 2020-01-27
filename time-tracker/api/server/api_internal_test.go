package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fgrimme/time-tracker/time-tracker/store"
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
	u  string // user id - used as the test case id in the mock store
	ts string // timestamp
	tz string // timezone
	p  string // period
}

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
	3: { // 500
		d: "expect store error to result in 500",
		e: errInternal,
		s: http.StatusInternalServerError,
		p: params{u: "3", ts: "0"}, // timzone and location can be empty
		b: []byte(fmt.Sprintf(`{"error":"%s"}`, errInternal.Error())),
	},
	// success
	4: { // 200
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

var startPeriodTests = []struct {
	d  string         // description of test case
	t  time.Time      // param t
	l  *time.Location // param loc
	p  string         // param period
	tr time.Time      // expected result
	e  error          // expected error
}{
	// error
	{
		d:  "expecting error due to invalid period",
		t:  time.Now(),
		l:  nil,
		p:  "invalid",
		tr: time.Now(),
		e:  errors.New("unknown period: invalid"),
	},
	// success - day
	{
		d:  "expecting same date when passing first day of week",
		t:  time.Date(2020, time.January, 01, 0, 0, 0, 0, time.UTC),
		l:  time.UTC,
		p:  "day",
		tr: time.Date(2020, time.January, 01, 0, 0, 0, 0, time.UTC),
	},
	{
		d:  "expecting start midnight when passing a later time on the same day",
		t:  time.Date(2020, time.January, 01, 03, 30, 10, 2, time.UTC),
		l:  time.UTC,
		p:  "day",
		tr: time.Date(2020, time.January, 01, 0, 0, 0, 0, time.UTC),
	},
	// success - week
	{
		d:  "expecting correct date when passing first day of week",
		t:  time.Date(2019, time.December, 30, 0, 0, 0, 0, time.UTC),
		l:  time.UTC,
		p:  "week",
		tr: time.Date(2019, time.December, 30, 0, 0, 0, 0, time.UTC),
	},
	{
		d:  "expecting correct date when passing third day of week",
		t:  time.Date(2020, time.January, 01, 0, 0, 0, 0, time.UTC),
		l:  time.UTC,
		p:  "week",
		tr: time.Date(2019, time.December, 30, 0, 0, 0, 0, time.UTC),
	},
	// success - month
	{
		d:  "expecting correct date when passing first day of month",
		t:  time.Date(2020, time.January, 01, 0, 0, 0, 0, time.UTC),
		l:  time.UTC,
		p:  "month",
		tr: time.Date(2020, time.January, 01, 0, 0, 0, 0, time.UTC),
	},
	{
		d:  "expecting correct date when passing last day of month",
		t:  time.Date(2020, time.January, 31, 0, 0, 0, 0, time.UTC),
		l:  time.UTC,
		p:  "month",
		tr: time.Date(2020, time.January, 01, 0, 0, 0, 0, time.UTC),
	},
}

func TestGetStartOfPeriod(t *testing.T) {
	for _, tc := range startPeriodTests {
		gotT, gotErr := getStartOfPeriod(tc.t, tc.l, tc.p)
		// unexpected errors
		if gotErr != nil && tc.e == nil {
			t.Fatalf("%s: unexpected err: %v", tc.d, gotErr)
		}
		// expected errors
		if gotErr == nil && tc.e != nil {
			t.Fatalf("%s: expected err: %v", tc.d, tc.e)
		}
		if gotErr != nil && tc.e != nil {
			if got, want := gotErr.Error(), tc.e.Error(); got != want {
				t.Errorf("%s:\nwant err\n%+v\ngot\n%+v", tc.d, want, got)
			}
			continue
		}
		if got, want := gotT, tc.tr; got != want {
			t.Errorf("%s:\nwant time\n%+v\ngot\n%+v", tc.d, want, got)
		}
	}
}
