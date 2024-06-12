create table configs
(
    type integer not null unique,
    data json    not null
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

create table turbines
(
    id                      varchar(50)              not null,
    tower_id                varchar(50)              not null,
    gen_bearing_to_coupling double precision         not null,
    coupling_to_turbine     double precision         not null,
    data                    json                     not null,
    created_at              timestamp  default now() not null,
    updated_at              timestamp default null,
    deleted_at              timestamp  default null,
    created_by              varchar(50) default ''   not null,

    primary key (id),
    constraint fk_tower
        foreign key (tower_id)
            references towers (id),
    add constraint fk_created_by
        foreign key (created_by)
            references users (id)
);
