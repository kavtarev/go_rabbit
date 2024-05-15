-- +goose up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists users2 (
  id uuid not null default uuid_generate_v4()
  , name text
  , surname text
  , email text
  , password text
);

-- +goose down

drop table users2;