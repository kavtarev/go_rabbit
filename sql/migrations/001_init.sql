-- +goose up

create extension if not exists "uuid-ossp";

create table if not exists users (
  id uuid not null default uuid_generate_v4()
  , name text
  , surname text
  , email text
  , password text
);

-- +goose down

drop table users;