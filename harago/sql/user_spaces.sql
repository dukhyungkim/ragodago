create table user_spaces
(
    id         bigserial
        primary key,
    name       varchar(16) not null,
    email      varchar(64) not null
        unique,
    space      varchar(32) not null
        unique,
    created_at timestamp with time zone
);

alter table user_spaces
    owner to harago;
