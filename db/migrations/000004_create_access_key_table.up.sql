create table "access_key"
(
    id ulid NOT NULL DEFAULT gen_ulid() PRIMARY KEY,
    project_id ulid references "project" (id),
    user_id ulid references "user" (id),
    key varchar(24),
    mask varchar(40),
    signature bytea,
    expires timestamp
);
