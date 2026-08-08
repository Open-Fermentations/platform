-- name: GetRolesWithPermissions :many
select sqlc.embed(r), sqlc.embed(p)
from "role" r
left join "permission" p on r.id = p.role_id;