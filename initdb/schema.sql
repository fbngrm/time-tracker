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
  stop_time_loc varchar(50) NOT NULL
);

INSERT INTO users(id) VALUES(42);

INSERT INTO
  time_records(
    user_id,
    name,
    start_time,
    start_time_loc,
    stop_time,
    stop_time_loc
  )
VALUES
  (
    42,
    'foo',
    '2020-01-01 00:00:00+01', /* 1577833200 */
    'Europe/Berlin',
    '2020-01-01 01:00:00+01', /* 1577836800 */
    'Europe/Berlin'
  ),
  (
    42,
    'foo',
    '2020-01-20 00:00:00+01', /* 1579474800 */
    'Europe/Berlin',
    '2020-01-20 01:00:00+01', /* 1579478400 */
    'Europe/Berlin'
  )
