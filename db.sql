create table configs
(
    type integer not null unique,
    data json    not null,
    status boolean NOT NULL DEFAULT false
);

create table divisions
(
    id         varchar(50)             not null primary key,
    name       varchar(50)             not null,
    created_at timestamp default now() not null,
    updated_at timestamp,
    deleted_at timestamp
);

create table towers
(
    id          varchar                 not null primary key,
    name        varchar(100)            not null,
    unit_number varchar(20)             not null,
    created_at  timestamp default now() not null,
    updated_at  timestamp,
    deleted_at  timestamp
);

create table users
(
    id            varchar(50)              not null,
    name          varchar(150)             not null,
    username      varchar(50)              not null unique,
    email         varchar(100)             not null unique,
    division_id   varchar(50)              not null,
    role          integer                  not null,
    status        integer                  not null,
    password_hash varchar,
    password_salt varchar,
    created_at    timestamp  default now() not null,
    updated_at    timestamp default null,
    deleted_at    timestamp  default null,

    primary key (id),
    constraint fk_division
        foreign key (division_id)
            references divisions (id)
);

comment on column users.role is 'super admin = 1, admin = 2, user = 3';

comment on column users.status is 'inactive = 0, active = 1, blocked by admin = 2';

CREATE TABLE
  public.turbines (
    id character varying(50) NOT NULL,
    tower_id character varying(50) NOT NULL,
    gen_bearing_to_coupling double precision NOT NULL,
    coupling_to_turbine double precision NOT NULL,
    data json NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NULL,
    deleted_at timestamp without time zone NULL,
    created_by character varying(50) NOT NULL DEFAULT ''::character varying,
    total_bolts integer NOT NULL DEFAULT 0,
    current_torque double precision NOT NULL DEFAULT 0,
    max_torque double precision NOT NULL DEFAULT 0,
    title character varying(150) NOT NULL DEFAULT ''::character varying
  );

ALTER TABLE
  public.turbines
ADD
  CONSTRAINT turbines_pkey PRIMARY KEY (id)