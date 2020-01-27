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
This is due to the possiblity that a user's timezone database is not up-to-date and thus delievers dates with a wrong UTC-offset.

E.g. if the European Union decides to abondon daylight saving time, the database on the user's system may not be updated in time while she starts or stops a session and reports a wrong offset.

By doing all time conversions and zone calculations in a central location, we can control the timezone database and guarantee correct conversions.
For more information about the tz-database/zoneinfo on UNIX systems, see: https://www.iana.org/time-zones

Theoretically, it was also possible to store future dates in a safe way without being affected by future, yet unknown, changes in timezone conversion rules.
Against the widespread believe that it is true to store dates in UTC format to be a bulletproof method, this is not true for storing future dates.

### Setup
This section assumes there is a go, make and git installation available on the system.
A Makefile is located at the project root which should be used to test, build and run the program.

### Build

```bash
make build # builds all services
```

### Run
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
The program consists of three services.
Services run in a Docker container, exposed by a gateway service which is reachable from the outside world.
For details on the setup the [docker-compose]() file.

### Database
A postgres database is used to store the time records for user sessions.
The schema is initialized on startup by the scipt in the `initdb` directory.
Four records are pre-populated on start-up as examples.

### Backend
A golang backend service provides the API to store and fetch time records.
Builds of the backend service will be placed in the `/bin` directory of the time-tracker service.
Binaries use the latest git commit hash or tag as a version.

#### API Endpoints

`GET /records?user_id=42&tz=Europe/Berlin&ts=1579688104&period=week

**Query parameters**
- user_id: [0-9]+ - the user ID, currenlty hardcoded to 42
- tz: [A-Za-z]+/[A-Za-z]+ - the user's time zone name according to the IANA zoneinfo definition
- ts: [0-9]+ - timestamp as number of seconds since UNIX epoch
- period: day|week|month - the time period of requested records

**Response**
```json
{
	"record_id": 5,
	"user_id": 42,
	"name": "hello world",
	"start_time": "27 Jan 2020 00:13:24",
	"start_loc": "Europe/Berlin",
	"stop_time": "27 Jan 2020 00:13:40",
	"stop_loc": "Europe/Berlin",
	"duration": "00:00:12"
}
```

**Role:**

During a typical day, thousands of drivers send their coordinates every 5 seconds to this endpoint.

**Behaviour**

Coordinates received on this endpoint are converted to [NSQ](https://github.com/nsqio/nsq) messages listened by the `Driver Location` service.

---

`GET /drivers/:id`

**Response**


### Frontend
Builds of the frontend are created during the build of the gateway service.


### Gatweay

## Dependency management
The program makes use of the standard library wherever possible.
For handling dependencies, go modules are used.
This requires to have a go version > 1.11 installed and setting `GO111MODULE=1`.
If the go version is >= 1.13, modules are enabled by default.
