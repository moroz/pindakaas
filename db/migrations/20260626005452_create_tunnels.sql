-- +goose Up
create table tunnels (
  id uuid not null primary key,
  subdomain text not null unique,
  username text not null unique,
  password_encrypted blob not null,
  inserted_at timestamp not null default (unixepoch()),
  updated_at timestamp not null default (unixepoch())
);

-- +goose Down
drop table hosts;
