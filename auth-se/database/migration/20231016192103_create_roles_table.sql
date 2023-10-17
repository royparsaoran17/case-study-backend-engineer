-- +goose Up
create table if not exists roles
(
    id         uuid
    constraint roles_pk
    primary key,
    name   varchar   not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp
);

create unique index roles_name_uindex
    on roles (name);

-- +goose Down
drop table if exists roles;
