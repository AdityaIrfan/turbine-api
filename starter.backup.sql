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
    total_bolts integer NOT NULL DEFAULT 0,
    current_torque double precision NOT NULL DEFAULT 0,
    max_torque double precision NOT NULL DEFAULT 0,
    title character varying(150) NOT NULL DEFAULT ''::character varying,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    created_by character varying(50) NOT NULL DEFAULT ''::character varying,
    deleted_at timestamp without time zone NULL,
    deleted_by character varying(50) DEFAULT ''::character varying
  );

ALTER TABLE
  public.turbines
ADD
  CONSTRAINT turbines_pkey PRIMARY KEY (id);

insert into "public"."towers" ("created_at", "deleted_at", "id", "name", "unit_number", "updated_at") values ('2024-06-12 21:51:22.970089', NULL, '01J06EX7EJHFXC47Y6AE8SSPTC', 'PLTA Sutami', '0H823A1', '2024-07-19 20:43:58.224267'), ('2024-06-12 22:03:55.305042', NULL, '01J06FM651MRQZC20PKZTQVE31', 'PLTA Sutami', '2', NULL), ('2024-06-16 05:18:32.057792', '2024-06-16 05:19:25.702557', '01J0EZP4QSX439KA7M3QZA54YF', 'PLTA Sutami', '3', NULL), ('2024-07-03 16:17:51.560531', NULL, '01J1VY5M13RVJAY7142H0BWQ9P', 'PLTA Sengguruh', '1', NULL), ('2024-07-03 16:17:51.865632', NULL, '01J1VY5M13RVJAY7142JBC8KCK', 'PLTA Sengguruh', '2', NULL), ('2024-07-03 16:17:52.046493', NULL, '01J1VY5M13RVJAY7142KS4H0HK', 'PLTA Wlingi', '1', NULL), ('2024-07-03 16:17:52.267142', NULL, '01J1VY5M13RVJAY7142PC61T6G', 'PLTA Wlingi', '2', NULL), ('2024-07-03 16:17:52.429293', NULL, '01J1VY5M13RVJAY7142R8DDZR6', 'PLTA Lodoyo', '1', NULL), ('2024-07-03 16:17:52.615735', NULL, '01J1VY5M13RVJAY7142RA1NA36', 'PLTA Tulung Agung', '1', NULL), ('2024-07-03 16:17:52.83501', NULL, '01J1VY5M13RVJAY7142S4MZJHQ', 'PLTA Tulung Agung', '2', NULL), ('2024-07-03 16:17:53.003122', NULL, '01J1VY5M13RVJAY7142SJK4XDK', 'PLTA Wonorejo', '1', NULL), ('2024-07-03 16:17:53.169802', NULL, '01J1VY5M13RVJAY7142T5SVTR1', 'PLTA Siman', '1', NULL), ('2024-07-03 16:17:53.361351', NULL, '01J1VY5M13RVJAY7142XBSFDQD', 'PLTA Siman', '2', NULL), ('2024-07-03 16:17:53.57721', NULL, '01J1VY5M13RVJAY7142Y3NK5QJ', 'PLTA Siman', '3', NULL), ('2024-07-03 16:17:53.751218', NULL, '01J1VY5M13RVJAY7142Z1EDJBW', 'PLTA Selorejo', '1', NULL), ('2024-07-03 16:17:53.922042', NULL, '01J1VY5M13RVJAY714301BY7EW', 'PLTA Mendalan', '1', NULL), ('2024-07-03 16:17:54.07845', NULL, '01J1VY5M13RVJAY71433T8NPMC', 'PLTA Mendalan', '2', NULL), ('2024-07-03 16:17:54.265985', NULL, '01J1VY5M13RVJAY71434720HST', 'PLTA Mendalan', '3', NULL), ('2024-07-03 16:17:54.490414', NULL, '01J1VY5M13RVJAY7143664QVSW', 'PLTA Mendalan', '4', NULL), ('2024-07-03 16:17:54.763913', NULL, '01J1VY5M13RVJAY71438KZ4N1Y', 'PLTA Golang', '1', NULL), ('2024-07-03 16:17:54.980758', NULL, '01J1VY5M13RVJAY714392N06RR', 'PLTA Golang', '2', NULL), ('2024-07-03 16:17:55.111366', NULL, '01J1VY5M13RVJAY71439QSFNHT', 'PLTA Golang', '3', NULL), ('2024-07-03 16:17:55.459045', NULL, '01J1VY5M13RVJAY7143C5EHQ0S', 'PLTA Giringan', '1', NULL), ('2024-07-03 16:17:56.062568', NULL, '01J1VY5M13RVJAY7143CMMMH3S', 'PLTA Giringan', '2', NULL), ('2024-07-03 16:17:56.397792', NULL, '01J1VY5M13RVJAY7143G5J49SC', 'PLTA Giringan', '3', NULL), ('2024-07-03 16:17:56.661394', NULL, '01J1VY5M13RVJAY7143KB3NS46', 'PLTA Ngebel', '1', NULL), ('2024-07-03 16:17:56.879869', NULL, '01J1VY5M13RVJAY7143Q9RZW2W', 'PLTA Cirata', '1', NULL), ('2024-07-03 16:17:57.233865', NULL, '01J1VY5M13RVJAY7143RPCMPMC', 'PLTA Cirata', '2', NULL), ('2024-07-03 16:17:57.360942', NULL, '01J1VY5M13RVJAY7143S5ZMG8C', 'PLTA Cirata', '3', NULL), ('2024-07-03 16:17:57.638864', NULL, '01J1VY5M13RVJAY7143SP8EVCP', 'PLTA Cirata', '4', NULL), ('2024-07-03 16:17:57.842257', NULL, '01J1VY5M13RVJAY7143T95F140', 'PLTA Cirata', '5', NULL), ('2024-07-03 16:17:57.998811', NULL, '01J1VY5M13RVJAY7143WAMC2RS', 'PLTA Cirata', '6', NULL), ('2024-07-03 16:17:58.241484', NULL, '01J1VY5M13RVJAY7143XPQ8CDY', 'PLTA Cirata', '7', NULL), ('2024-07-03 16:17:58.433067', NULL, '01J1VY5M13RVJAY7143XZ37RG3', 'PLTA Cirata', '8', NULL);
insert into "public"."divisions" ("created_at", "deleted_at", "id", "name", "updated_at") values ('2024-06-12 21:32:06', NULL, '01J06DRG2A5DEQY8CZH7GDE2N2', 'Engineer', NULL), ('2024-07-04 00:27:36.774591', NULL, '01J1WT6CJ6X460H4XNPR0MCHZF', 'HR', NULL), ('2024-07-04 00:27:53.25379', NULL, '01J1WT6WN5T6HWJY2SBCHX6PVK', 'Devops', NULL);
insert into "public"."users" ("created_at", "deleted_at", "division_id", "email", "id", "name", "password_hash", "password_salt", "role", "status", "updated_at", "username") values ('2024-06-12 14:35:01.929863', NULL, '01J06DRG2A5DEQY8CZH7GDE2N2', 'aditya@gmail.com', '01J06DSG8ZTVZTBW5NZ87XRVMX', 'aditya fullname', '8f82307499f0b213ea10761fb005965581a1dfe0ca6970fbf7e5994800b5f02a5a3a18ff611fd061eda1c4863674bb3403f9e9faf10196b2993de3bd1f5ebb9a66489290de763274c5600c79cb074240c75156f5de743b84', '688a376bf02e3828a67ce51d82e093c4e5c9e6aa20db18e68ece141d382a1c29e307fb4dd99c71e7681f5941ed61cc1e65c9aa5ad5b31557ffa19f447113b54a4ab1fca59d9209dafb0ba1a6aa1a6471249eb35baea36d86688effd1b359b67e93527413', 1, 1, '2024-07-30 16:41:51.670738', 'aditya');
insert into "public"."configs" ("data", "status", "type") values ('{"Long":-8.160927807736735,"Lat":112.44418490930914,"CoverageArea":5,"CoverageAreaType":"kilometer"}', false, 1);