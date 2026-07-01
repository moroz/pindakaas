-- name: InsertTunnel :one
insert into tunnels (id, subdomain, username, password_encrypted, user_id)
values (?, ?, ?, ?, ?)
returning *;

-- name: GetTunnelByUsername :one
select * from tunnels where username = ?;

-- name: GetTunnelForUser :one
select * from tunnels where id = @tunnel_id and user_id = ?;

-- name: ListTunnelsForUser :many
select * from tunnels where user_id = ? order by id desc;

-- name: DeleteTunnelForUser :exec
delete from tunnels where id = @tunnel_id and user_id = ?;
