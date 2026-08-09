-- name: CreateBatch :one
insert into "batch" ("name", user_id)
values ($1, $2)
returning id, "name", user_id, created, modified;

-- name: DeleteBatch :exec
delete from "batch" where id = @id and user_id = @user_id;

-- name: SearchBatches :many
select *, count(id) over() as total from "batch"
where user_id = @user_id
    AND "name" LIKE @name::text
order by created desc
limit @limitVal::integer
offset @offsetVal::integer;

-- name: GetBatchById :one
select * from "batch" where id = @id and user_id = @user_id;

-- name: UpdateBatch :one
update "batch" set "name" = @name, modified = current_timestamp where id = @id and user_id = @user_id
returning *;