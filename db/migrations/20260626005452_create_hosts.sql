-- +goose Up
create table hosts (
  id integer not null primary key,
  subdomain text not null unique,
  username text not null unique,
  password_hash text not null,
  inserted_at bigint not null default (unixepoch()),
  updated_at bigint not null default (unixepoch())
);

-- +goose Down
drop table hosts;
