-- name: CreateUser :exec
insert into "user" (username, password)
values ($1, $2);

-- name: GetUserByUsername :one
select id, username, created, modified
from "user"
where id = $1;