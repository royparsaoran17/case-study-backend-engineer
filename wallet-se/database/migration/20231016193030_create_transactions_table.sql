-- +goose Up
create table if not exists transactions
(
    id              uuid
        constraint transactions_pk
        primary key,
    type            varchar   not null,
    amount           float   not null,
    reference_id    uuid   not null,
    status           varchar   not null,
    transaction_at   timestamp   not null,
    transaction_by   uuid   not null,
    wallet_id        uuid   not null
        constraint wallet_transaction_uid_fk
        references wallet (id)
        on delete cascade,
    created_at  timestamp not null,
    updated_at  timestamp not null,
    deleted_at  timestamp
);

create unique index transactions_reference_id_uindex
    on transactions (reference_id);

-- +goose Down
drop table if exists transactions;
