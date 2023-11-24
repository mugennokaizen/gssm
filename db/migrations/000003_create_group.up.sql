create table "secret_group"
(
    id ulid NOT NULL DEFAULT gen_ulid() PRIMARY KEY,
    name varchar(128),
    prefix varchar(32),
    project_id ulid references "project" (id)
);
