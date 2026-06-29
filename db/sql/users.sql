-- name: FindUserByUserToken :one
select u.* from users u
join user_tokens ut on u.id = ut.user_id
where ut.inserted_at + cast(@validity as bigint) > unixepoch()
and ut.token = @token and ut.context = @context;

-- name: InsertUserToken :one
insert into user_tokens (id, user_id, token, context)
values (?, ?, ?, ?) returning *;

-- name: GetUserByEmail :one
select * from users where email = ?;

-- name: InsertUser :one
insert into users (id, email, given_name, family_name, avatar)
values (?, ?, ?, ?, ?) returning *;

-- name: UpsertUser :one
insert into users (id, email, given_name, family_name, avatar)
values (?, ?, ?, ?, ?)
on conflict (email)
do update set
  given_name = excluded.given_name,
  family_name = excluded.family_name,
  avatar = excluded.avatar,
  updated_at = (unixepoch())
returning *;

