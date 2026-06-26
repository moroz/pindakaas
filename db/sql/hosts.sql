-- name: ListHosts :many
select * from hosts order by id;

-- name: GetHostByUsername :one
select * from hosts where username = ?;
