-- name: CreateBatch :one
insert into "batch" ("name", user_id)
values ($1, $2)
returning id, "name", user_id, created, modified;

-- name: DeleteBatch :exec
delete from "batch" where id = $1;

-- name: SearchBatches :many
select *, count(id) over() as total from "batch"
where "name" like @name::text
order by created
limit @limitVal::integer
offset @offsetVal::integer;

-- name: GetBatchById :one
select * from "batch" where id = $1;

-- name: UpdateBatch :one
update "batch" set "name" = $2, modifie = current_timestamp where id = $1
returning *;