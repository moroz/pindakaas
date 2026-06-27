-- +goose Up
-- TODO: Rename `hosts` to `tunnels`
create table hosts (
  id uuid not null primary key,
  subdomain text not null unique,
  username text not null unique,
  password_hash text not null,
  inserted_at timestamp not null default (unixepoch()),
  updated_at timestamp not null default (unixepoch())
);

-- +goose Down
drop table hosts;
