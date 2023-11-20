CREATE EXTENSION ulid;
create table "user"
(
    id ulid NOT NULL DEFAULT gen_ulid() PRIMARY KEY,
    email         varchar(128),
    password_hash bytea,
    salt          bytea
);