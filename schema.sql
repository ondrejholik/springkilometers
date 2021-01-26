CREATE TABLE users (
  user_id serial PRIMARY KEY,  -- implicit primary key constraint
  username text NOT NULL,
  pass text NOT NULL,
  salt text NOT NULL
);

CREATE TABLE trips (
  trip_id  serial PRIMARY KEY,
  name     text NOT NULL,
  content  text ,
  km       decimal,
  with_bike bool,
  created  date NOT NULL DEFAULT CURRENT_DATE
  
);

CREATE TABLE trip_user (
  trip_id    int REFERENCES trips (trip_id) ON UPDATE CASCADE ON DELETE CASCADE
, user_id int REFERENCES users (user_id) ON UPDATE CASCADE
, CONSTRAINT trip_user_pkey PRIMARY KEY (trip_id, user_id)  -- explicit pk
);


