-- name: CreateDevice :one
insert into "device" ("name", mac_address, user_id)
values ($1, $2, $3)
returning id, "name", mac_address, user_id, created, modified;
