-- name: CreateDevice :one
insert into "device" ("name", mac_address, user_id)
values ($1, $2, $3)
returning id, "name", mac_address, user_id, created, modified;

-- name: SearchDevices :many
select *, count(*) over() as total
from "device"
where "name" like @search::text
order by
  case when sqlc.arg('orderCol')::text = 'name' and sqlc.arg('asc')::bool then "device"."name" end asc nulls last,
  case when sqlc.arg('orderCol')::text = 'name' and sqlc.arg('asc')::bool != true then "device"."name" end desc nulls last,
  case when sqlc.arg('orderCol')::text = 'created' and sqlc.arg('asc')::bool then "device"."created" end asc nulls last,
  case when sqlc.arg('orderCol')::text = 'created' and sqlc.arg('asc')::bool != true then "device"."created" end desc nulls last,
  case when sqlc.arg('orderCol')::text = 'modified' and sqlc.arg('asc')::bool then "device"."modified" end asc nulls last,
  case when sqlc.arg('orderCol')::text = 'modified' and sqlc.arg('asc')::bool != true then "device"."modified" end desc nulls last,
  case when sqlc.arg('orderCol')::text = 'id' and sqlc.arg('asc')::bool then "device"."id" end asc nulls last,
  case when sqlc.arg('orderCol')::text = 'id' and sqlc.arg('asc')::bool != true then "device"."id" end desc nulls last,
  case when sqlc.arg('orderCol')::text = 'user_id' and sqlc.arg('asc')::bool then "device"."user_id" end asc nulls last,
  case when sqlc.arg('orderCol')::text = 'user_id' and sqlc.arg('asc')::bool != true then "device"."user_id" end desc nulls last,
  case when sqlc.arg('orderCol')::text = 'mac_address' and sqlc.arg('asc')::bool then "device"."mac_address" end asc nulls last,
  case when sqlc.arg('orderCol')::text = 'mac_address' and sqlc.arg('asc')::bool != true then "device"."mac_address" end desc nulls last,
  "device"."created" asc nulls last
limit @limitVal
offset @offsetVal;