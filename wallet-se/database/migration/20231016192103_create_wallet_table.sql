-- +goose Up
create table if not exists wallet
(
    id         uuid
    constraint wallet_pk
    primary key,
    balance     float     not null,
    status      varchar not null,
    enabled_at  timestamp   not null,
    owned_by    uuid      not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp
);

create unique index wallet_owned_uindex
    on wallet (owned_by);

-- +goose Down
drop table if exists wallet;
