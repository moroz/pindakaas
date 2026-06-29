-- name: ListTunnels :many
select * from tunnels order by id;

-- name: GetTunnelByUsername :one
select * from tunnels where username = ?;
