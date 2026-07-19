-- name: CreateDevice :one
insert into "device" ("name", mac_address, user_id)
values ($1, $2, $3)
returning id, "name", mac_address, user_id, created, modified;

-- name: SearchDevices :many
select *, count(id) over() as total from "device"
where "name" like @name::text
order by created 
limit @limitVal::integer 
offset @offsetVal::integer;