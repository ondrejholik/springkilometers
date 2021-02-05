-- Table: public.trips

-- DROP TABLE public.trips;

CREATE TABLE public.trips
(
    id serial not null,
    name text COLLATE pg_catalog."default" NOT NULL,
    content text COLLATE pg_catalog."default",
    km numeric,
    withbike boolean,
    created_on date DEFAULT CURRENT_DATE,
    updated_on date DEFAULT CURRENT_DATE,
    modified_on date DEFAULT CURRENT_DATE,
    deleted_on date DEFAULT CURRENT_DATE,
    author text COLLATE pg_catalog."default" NOT NULL,
    tiny text COLLATE pg_catalog."default" NOT NULL,
    medium text COLLATE pg_catalog."default" NOT NULL,
    small text COLLATE pg_catalog."default" NOT NULL,
    large text COLLATE pg_catalog."default" NOT NULL,
    year integer,
    month integer,
    day integer,
    hour integer,
    minute integer,
    "timestamp" integer,
    gpx text COLLATE pg_catalog."default",
    CONSTRAINT trips_pkey PRIMARY KEY (id)
);

TABLESPACE pg_default;

ALTER TABLE public.trips
    OWNER to postgres;


-- Table: public.users

-- DROP TABLE public.users;

CREATE TABLE public.users
(
    id serial NOT NULL ,
    username text COLLATE pg_catalog."default" NOT NULL,
    password text COLLATE pg_catalog."default" NOT NULL,
    salt text COLLATE pg_catalog."default" NOT NULL,
    created_on date NOT NULL DEFAULT CURRENT_DATE,
    updated_on date NOT NULL DEFAULT CURRENT_DATE,
    modified_on date NOT NULL DEFAULT CURRENT_DATE,
    deleted_on date NOT NULL DEFAULT CURRENT_DATE,
    CONSTRAINT users_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.users
    OWNER to postgres;

-- Table: public.user_trip

-- DROP TABLE public.user_trip;

CREATE TABLE public.user_trip
(
    user_id serial NOT NULL,
    trip_id integer NOT NULL,
    CONSTRAINT trip_user_pkey PRIMARY KEY (trip_id, user_id),
    CONSTRAINT user_trip_trip_id_fkey FOREIGN KEY (trip_id)
        REFERENCES public.trips (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT user_trip_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE public.user_trip
    OWNER to postgres;

