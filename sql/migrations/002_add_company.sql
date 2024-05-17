-- +goose up

alter table users add primary key (id);

create table if not exists companies (
  id uuid not null default uuid_generate_v4()
  , name text
  , inn text
  , email text
  , user_id uuid references users(id) on delete cascade
);

-- +goose down

drop table companies;