create schema collection;

create table collection.events
(
    id          uuid                    not null,
    title       varchar(100)            not null,
    description varchar,
    start_time  timestamp               not null,
    end_time    timestamp               not null,
    created_at  timestamp default now() not null
);
