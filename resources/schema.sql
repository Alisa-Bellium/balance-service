CREATE TABLE "users"
(
    id      bigint not null primary key,
    balance bigint not null check ( balance >= 0)
);

CREATE TABLE "transactions"
(
    id                      serial    not null
        constraint transactions_pk
            primary key,
    user_id                 bigint    not null
        constraint transactions_users_fk0
            references users,
    created_at              timestamp not null,
    amount                  bigint    not null,
    service_id              bigint,
    order_id                bigint,
    is_reserve_account      boolean   not null,
    canceled_transaction_id bigint
        constraint transactions_cancelled_fk0
            references transactions
);

CREATE TABLE "reports"
(
    id                  serial    not null primary key,
    month               int       not null,
    year                int       not null,
    created_at          timestamp not null,
    file_path           TEXT      null,
    last_transaction_id bigint
        constraint reports_transactions_null_fk
            references transactions
            on update cascade on delete cascade,
    constraint reports_unique_date_transaction_id
        unique (year, month, last_transaction_id)
);





