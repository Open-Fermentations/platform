-- name: CreateUser :exec
insert into "user" (username, password)
values ($1, $2);

-- name: GetUserByUsername :one
select *
from "user"
where username = $1;

-- name: GetUserByUsernameWithPasswordAndRolesAndPermissions :many
select sqlc.embed(u), sqlc.embed(r), sqlc.embed(p) 
from "user" u 
left join "user_role" ur on ur.user_id = u.id 
left join "user_permission" up on up.user_id = u.id 
left join "role" r on ur.role_id = r.id 
left join "role_permission" rp on r.id = rp.role_id 
left join "permission" p on up.permission_id = p.id or p.id = rp.permission_id 
where username = @username;

-- name: GetUserById :one
select id, username, created, modified
from "user"
where id = $1;