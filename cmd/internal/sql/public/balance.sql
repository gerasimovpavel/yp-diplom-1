create table if not exists balance
(
    user_id  uuid
        constraint user_id
            unique,
    accrual  numeric(19, 2) default 0,
    withdraw numeric(19, 2) default 0,
    current  numeric(19, 2) default 0
);


