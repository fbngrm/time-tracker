# Pento tech challenge
The program implements the required functionality.

Further it is possible for a user to stop and continue a session.
Sessions can be started in one timezone and stopped in another.
The program will always show the local wall clock time of the timezone as well as the zone names, the session was started and stopped in.
In other words, it is possible for a user to start a session in e.g. Copenhagen (UTC+1), travel to and stop the session in London (UTC) and view the session in another timezone but still see the relevant timezone information.

Note, sessions must be saved to survive the reboot of the system.

## Storing time records
User input is provided as a timestamp of seconds since UNIX epoch and a tz-database location name.
The program does not rely on the user's timezone database and converts provided timestamps to the user's location before operating on or persisting them.
This is due to the possibility that a user's timezone database is not up-to-date and thus delivers dates with a wrong UTC-offset.

E.g. the European Union decides to abandon daylight saving time, the database on the user's system may not be updated in time when the user starts or stops a session and reports a wrong offset.

By doing all time conversions and zone calculations in a central location, we can control the timezone database and guarantee correct conversions.
For more information about the tz-database/zoneinfo on UNIX systems, see: https://www.iana.org/time-zones

Theoretically, it was also possible to store future dates in a safe way without being affected by future, yet unknown, changes in timezone conversion rules.

### Setup
This section assumes there is a go, make, Docker and git installation available on the system.
A Makefile is located at the project root which should be used to test, build and run the program.

### Build

```bash
make build # builds all services
```

### Run
Serve the frontend, open a browser and navigate to `localhost` after all services have started.
```
make build
make run # starts all services
```

### Tests
There are several targets available to run tests.

```bash
make test # runs tests
make test-cover # creates a coverage profile
make test-race # tests service for race conditions
```

### Lint
There is a lint target which runs [golangci-lint](https://github.com/golangci/golangci-lint) in a docker container.

```bash
make lint
```

## Architecture
The program consists of three services running in a Docker container, exposed by a gateway which is reachable from the outside world.
For details on the setup the [docker-compose](https://github.com/fbngrm/time-tracker/blob/master/docker-compose.yaml) file.

### Database
A postgres database is used to store the time records for user sessions.
The schema is initialized on startup by the script in the `initdb` directory.
Four records are populated on start-up as examples.

### Backend
A golang backend service provides the API to store and fetch time records.
Builts of the backend service will be placed in the `/bin` directory of the time-tracker service.
Binaries use the latest git commit hash or tag as a version.

#### Private API Endpoints

`POST /record`

**Payload**

```json
{
	"user_id": 42,
	"name": "eat pølser",
	"start_time": 1579962216,
	"start_loc": "Europe/Copenhagen",
	"stop_time": 1579965816,
	"stop_loc": "Europe/Copenhagen",
	"duration": 3600
}
```

**Response**

```json
{
	"record_id": 9,
	"user_id": 42,
	"name": "eat pølser",
	"start_time": "25 Jan 2020 15:23:36",
	"start_loc": "Europe/Copenhagen",
	"stop_time": "25 Jan 2020 16:23:36",
	"stop_loc": "Europe/Copenhagen",
	"duration": "01:00:00"
}
```

**Role**

Persist a time record of a user session in the database.

**Behaviour**

The provided timestamps and timezones are used to get the start and stop times in the provided locations.
The time record is stored in the datastore which returns the record's ID.
A JSON representation of the record with the generated ID and formatted times and duration is returned.

---

`GET /records?user_id=42&tz=Europe/Berlin&ts=1579688104&period=week`

**Query parameters**

- `user_id: [0-9]+` - the user ID, currently hardcoded to 42
- `tz: [A-Za-z]+/[A-Za-z]+` - the user's time zone name according to the IANA zoneinfo definition
- `ts: [0-9]+` - timestamp as number of seconds since UNIX epoch
- `period: day|week|month` - the time period of requested records

**Response**
```json
[{
	"record_id": 5,
	"user_id": 42,
	"name": "hello world",
	"start_time": "27 Jan 2020 00:13:24",
	"start_loc": "Europe/Berlin",
	"stop_time": "27 Jan 2020 00:13:40",
	"stop_loc": "Europe/Berlin",
	"duration": "00:00:12"
}, ...]
```

**Role**

Fetch a list of records for a certain user for the current day, week or month.

**Behaviour**

The provided timestamp and timezone are used to get the time in the users location.
The start of the first day for the provided period in the provided location is calculated as the start date.
A list of JSON representations of all time records with a stop date past the start date is returned.

---

### Frontend
A react/redux frontend provides the user interface to interact with the backend API.
Builds of the frontend are created during the build of the gateway service and are served as static files.

### Gatweay
A nginx gateway is responsible for serving the frontend as well as providing public API enpoints.
Requests to the public endpoints are forwarded to the backend service.

#### Public API Endpoints

`GET /`

Serves the frontend.

`POST /time-tracker/record`

Forwarded to the private `/record` endpoint.

`GET /time-tracker/records`

Forwarded to the private `/records` endpoint.

## Dependency management
The program makes use of the standard library wherever possible.
For handling dependencies, go modules are used.
This requires to have a go version > 1.11 installed and setting `GO111MODULE=1`.
If the go version is >= 1.13, modules are enabled by default.

## Issues
The timer displayed in the frontend is not accurate when stopped and continued multiple times.
The time tracked and send to the backend is accurate though.
This just affects user experience/visualization and is not a data accuracy issue.
Still, this needed to be fixed in a future iteration.
