-- +goose Up
create table if not exists customers
(
    id          uuid
        constraint customers_pk
        primary key,
    name        varchar   not null,
    phone       varchar   not null,
    password    varchar   not null,
    role_id     uuid   not null
        constraint customer_role_uid_fk
        references roles (id)
        on delete cascade,
    created_at  timestamp not null,
    updated_at  timestamp not null,
    deleted_at  timestamp
);

create unique index customers_phone_uindex
    on customers (phone);

-- +goose Down
drop table if exists customers;
