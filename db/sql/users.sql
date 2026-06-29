-- name: FindUserByUserToken :one
select u.* from users u
join user_tokens ut on u.id = ut.user_id
where ut.inserted_at + cast(@validity as bigint) > unixepoch()
and ut.token = @token and ut.context = @context;

-- name: InsertUserToken :one
insert into user_tokens (id, user_id, token, context)
values (?, ?, ?, ?) returning *;
