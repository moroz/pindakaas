-- +goose Up
create table users (
  id uuid not null primary key,
  email text not null unique,
  user_role user_role not null default 'Regular',
  inserted_at timestamp not null default unixepoch(),
  updated_at timestamp not null default unixepoch()
);

-- +goose Down
drop table users;
