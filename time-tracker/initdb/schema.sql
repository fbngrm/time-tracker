CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    first_name varchar(32) NOT NULL,
    last_name  varchar(32) NOT NULL
);

CREATE TABLE time_record (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) NOT NULL,
  name varchar(256) NOT NULL,
  start_time TIMESTAMP WITH TIME ZONE,
  end_time TIMESTAMP WITH TIME ZONE
);
