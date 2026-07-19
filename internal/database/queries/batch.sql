-- name: CreateBatch :one
insert into "batch" ("name", user_id)
values ($1, $2)
returning id, "name", user_id, created, modified;

-- name: DeleteBatch :exec
delete from "batch" where id = $1;