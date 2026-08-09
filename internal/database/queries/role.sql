-- name: GetRolesWithPermissions :many
select sqlc.embed(r), sqlc.embed(p)
from "role" r
left join "role_permission" rp on r.id = rp.role_id
left join "permission" p on p.id = rp.permission_id;