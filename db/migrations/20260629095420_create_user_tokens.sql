-- +goose Up
create table user_tokens (
  id uuid not null primary key,
  user_id uuid not null references users (id) on delete cascade,
  token blob not null unique,
  context user_token_context not null,
  inserted_at timestamp not null default (unixepoch())
);

-- +goose Down
drop table user_tokens;
