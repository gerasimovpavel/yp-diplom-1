create table if not exists orders
(
    number      varchar(20)
        constraint number
            unique,
    status      varchar(50),
    accrual     numeric(19, 2) default 0,
    uploaded_at timestamp with time zone,
    user_id     uuid,
    order_id    uuid
);


