-- name: CreateBatch :one
insert into "batch" ("name", user_id)
values ($1, $2)
returning id, "name", user_id, created, modified;

-- name: DeleteBatch :exec
delete from "batch" where id = @id and user_id = @userId;

-- name: SearchBatches :many
select *, count(id) over() as total from "batch"
where user_id = @userId
    AND "name" LIKE @name::text
order by created desc
limit @limitVal::integer
offset @offsetVal::integer;

-- name: GetBatchById :one
select * from "batch" where id = @id and user_id = @userId;

-- name: UpdateBatch :one
update "batch" set "name" = @name, modified = current_timestamp where id = @id and user_id = @userId
returning *;