CREATE TABLE users (
    id INT PRIMARY KEY
);

CREATE TABLE time_records (
  id BIGSERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) NOT NULL,
  name varchar(256),
  start_time TIMESTAMP NOT NULL,
  start_time_tz varchar(50) NOT NULL,
  stop_time TIMESTAMP NOT NULL,
  stop_time_tz varchar(50) NOT NULL
);

INSERT INTO users(id) VALUES
(42);

INSERT INTO time_records(user_id,name,start_time,start_time_tz,stop_time,stop_time_tz) VALUES
(42,'foo','2020-01-01 00:00:00','Europe/Berlin','2020-01-01 02:00:00','Europe/Berlin'),
(42,'bar','2020-01-19 12:00:00','Europe/Berlin','2020-01-19 16:00:00','Europe/Berlin');


