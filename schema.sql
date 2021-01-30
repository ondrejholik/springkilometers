CREATE TABLE users (
  id serial PRIMARY KEY,  -- implicit primary key constraint
  username text NOT NULL,
  pass text NOT NULL,
  salt text NOT NULL,
  created_on  date NOT NULL DEFAULT CURRENT_DATE,
  updated_on date NOT NULL DEFAULT CURRENT_DATE,
  modified_on date NOT NULL DEFAULT CURRENT_DATE,
  deleted_on date NOT NULL DEFAULT CURRENT_DATE
);

CREATE TABLE trips (
  id  serial PRIMARY KEY,
  name     text NOT NULL,
  content  text ,
  km       decimal,
  withbike bool,
  created_on  date NOT NULL DEFAULT CURRENT_DATE,
  updated_on date NOT NULL DEFAULT CURRENT_DATE,
  modified_on date NOT NULL DEFAULT CURRENT_DATE,
  deleted_on date NOT NULL DEFAULT CURRENT_DATE
  
);

CREATE TABLE trip_user (
  trip_id    int REFERENCES trips (id) ON UPDATE CASCADE ON DELETE CASCADE
, user_id int REFERENCES users (id) ON UPDATE CASCADE
, CONSTRAINT trip_user_pkey PRIMARY KEY (trip_id, user_id)  -- explicit pk
);


