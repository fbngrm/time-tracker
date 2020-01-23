CREATE TABLE users (
    id INT PRIMARY KEY
);

CREATE TABLE time_records (
  id BIGSERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) NOT NULL,
  name varchar(256),
  start_time TIMESTAMP WITH TIME ZONE NOT NULL,
  start_time_loc varchar(50) NOT NULL,
  stop_time TIMESTAMP WITH TIME ZONE NOT NULL,
  stop_time_loc varchar(50) NOT NULL,
  duration BIGINT NOT NULL
);

INSERT INTO users(id) VALUES(42);

INSERT INTO
  time_records(
    user_id,
    name,
    start_time,
    start_time_loc,
    stop_time,
    stop_time_loc,
    duration
  )
VALUES
  (
    42,
    'foo',
    '2020-01-01 00:00:00+01',
    'Europe/Berlin',
    '2020-01-01 01:00:00+01',
    'Europe/Berlin',
    3600
  ),
  (
    42,
    'bar',
    '2020-01-10 00:00:00+01',
    'Europe/Berlin',
    '2020-01-10 00:00:00+00',
    'Europe/London',
    3600
  ),
  (
    42,
    'baz',
    '2020-01-20 00:00:00+01',
    'Europe/Berlin',
    '2020-01-20 01:00:00+01',
    'Europe/Berlin',
    3600
  ),
  (
    42,
    'foobar',
    '2020-01-21 00:00:00+01',
    'Europe/Berlin',
    '2020-01-21 20:00:00+09',
    'Asia/Tokyo',
    43200
  )
