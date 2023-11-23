create table "project"
(
    id ulid NOT NULL DEFAULT gen_ulid() PRIMARY KEY,
    name varchar(128),
    creator_id ulid references "user" (id)
);

create table "user_to_project"
(
   id ulid not null default gen_ulid() primary key,
   project_id ulid references "project" (id),
   user_id ulid references "user" (id),
   permission smallint
);

create index project_user_idx
    on user_to_project (project_id, user_id)