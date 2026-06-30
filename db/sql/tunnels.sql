-- name: InsertTunnel :one
insert into tunnels (id, subdomain, username, password_hash, user_id)
values (?, ?, ?, ?, ?)
returning *;

-- name: ListTunnels :many
select * from tunnels order by id;

-- name: GetTunnelByUsername :one
select * from tunnels where username = ?;

-- name: GetTunnelForUser :one
select * from tunnels where id = @tunnel_id and user_id = ?;

-- name: ListTunnelsForUser :many
select * from tunnels where user_id = ? order by id;

-- name: DeleteTunnelForUser :exec
delete from tunnels where id = @tunnel_id and user_id = ?;
