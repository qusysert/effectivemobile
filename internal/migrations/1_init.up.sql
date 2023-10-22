create table if not exists "user"
(
    id         serial
        constraint user_pk
            primary key,
    name       text not null,
    surname    text not null,
    patronymic text,
    age int,
    gender text,
    nation text
);

