create table if not exists users
(
    user_id  uuid         not null
        primary key,
    login    varchar(128) not null
        constraint login
            unique,
    password varchar(128) not null
);



