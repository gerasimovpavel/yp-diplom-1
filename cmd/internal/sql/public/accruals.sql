create table if not exists accruals
(
    "order" varchar(128) not null,
    summa   double precision
);

create index if not exists "order"
    on accruals ("order" varchar_ops);

