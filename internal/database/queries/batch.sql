-- name: CreateBatch :one
insert into "batch" ("name", user_id)
values ($1, $2)
returning id, "name", user_id, created, modified;

-- name: DeleteBatch :exec
delete from "batch" where id = $1;

-- name: GetBatches :many
select *, count(id) over() as total from "batch"
where "name" like $1::text
order by created
limit $2
offset $3;

-- name: GetBatchById :one
select * from "batch" where id = $1;

-- name: UpdateBatch :one
update "batch" set "name" = $2 where id = $1
returning id, "name", "user_id", created, modified;