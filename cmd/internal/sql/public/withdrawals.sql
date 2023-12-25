create table if not exists withdrawals
(
    "order"      varchar(128),
    summa        numeric(19, 2),
    processed_at timestamp with time zone,
    user_id      uuid
);


