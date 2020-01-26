package store_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/PentoHQ/tech-challenge-time/time-tracker/store"
)

var locs = make(map[string]*time.Location)

func initLoc(t *testing.T, loc string) {
	// get the start time in the given location
	l, err := time.LoadLocation(loc)
	if err != nil {
		t.Fatal(err)
	}
	locs[loc] = l
}

type unmarshalTest struct {
	d   string           // test case description
	in  []byte           // input as JSON string
	out store.TimeRecord // expected output
	e   error            // expected error
}

func TestUnmarshall(t *testing.T) {
	initLoc(t, "Europe/London")
	initLoc(t, "Asia/Tokyo")
	unmarhsalTests := []unmarshalTest{
		unmarshalTest{
			in: []byte(`{"user_id":3,"name":"foo","start_time":1577833200,"start_loc":"Europe/London","stop_time":1577836800,"stop_loc":"Europe/London", "duration":3600}`),
			out: store.TimeRecord{
				RecordID: 0,
				UserID:   3,
				Name:     "foo",
				Start:    time.Date(2019, time.December, 31, 23, 0, 0, 0, locs["Europe/London"]),
				StartLoc: "Europe/London",
				Stop:     time.Date(2020, time.January, 01, 0, 0, 0, 0, locs["Europe/London"]),
				StopLoc:  "Europe/London",
				Duration: 3600,
			},
		},
		unmarshalTest{
			in: []byte(`{"user_id":3,"name":"foo","start_time":1577833200,"start_loc":"Asia/Tokyo","stop_time":1577836800,"stop_loc":"Asia/Tokyo", "duration":3600}`),
			out: store.TimeRecord{
				RecordID: 0,
				UserID:   3,
				Name:     "foo",
				Start:    time.Date(2020, time.January, 1, 8, 0, 0, 0, locs["Asia/Tokyo"]),
				StartLoc: "Asia/Tokyo",
				Stop:     time.Date(2020, time.January, 1, 9, 0, 0, 0, locs["Asia/Tokyo"]),
				StopLoc:  "Asia/Tokyo",
				Duration: 3600,
			},
		},
	}
	for _, tc := range unmarhsalTests {
		var got store.TimeRecord
		gotErr := got.UnmarshalJSON(tc.in)
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
		if want := tc.out; !reflect.DeepEqual(got, want) {
			t.Errorf("%s:\nwant time\n%+v\ngot\n%+v", tc.d, want, got)
		}
	}
}

type marshalTest struct {
	d   string           // test case description
	in  store.TimeRecord // input
	out []byte           // expected output
	e   error            // expected error
}

func TestMarshall(t *testing.T) {
	marhsalTests := []marshalTest{
		marshalTest{
			out: []byte(`{"record_id":0,"user_id":3,"name":"foo","start_time":"31 Dec 2019 23:00:00","start_loc":"Europe/London","stop_time":"01 Jan 2020 00:00:00","stop_loc":"Europe/London","duration":"01:00:00"}`),
			in: store.TimeRecord{
				RecordID: 0,
				UserID:   3,
				Name:     "foo",
				Start:    time.Date(2019, time.December, 31, 23, 0, 0, 0, locs["Europe/London"]),
				StartLoc: "Europe/London",
				Stop:     time.Date(2020, time.January, 01, 0, 0, 0, 0, locs["Europe/London"]),
				StopLoc:  "Europe/London",
				Duration: 3600,
			},
		},
		marshalTest{
			out: []byte(`{"record_id":0,"user_id":3,"name":"foo","start_time":"01 Jan 2020 08:00:00","start_loc":"Asia/Tokyo","stop_time":"01 Jan 2020 09:00:00","stop_loc":"Asia/Tokyo","duration":"01:00:00"}`),
			in: store.TimeRecord{
				RecordID: 0,
				UserID:   3,
				Name:     "foo",
				Start:    time.Date(2020, time.January, 1, 8, 0, 0, 0, locs["Asia/Tokyo"]),
				StartLoc: "Asia/Tokyo",
				Stop:     time.Date(2020, time.January, 1, 9, 0, 0, 0, locs["Asia/Tokyo"]),
				StopLoc:  "Asia/Tokyo",
				Duration: 3600,
			},
		},
	}
	for _, tc := range marhsalTests {
		gotB, gotErr := tc.in.MarshalJSON()
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
		if got, want := string(gotB), string(tc.out); got != want {
			t.Errorf("%s:\nwant time\n%+v\ngot\n%+v", tc.d, want, got)
		}
	}
}
