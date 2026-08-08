package model

import (
	"open-fermentations/internal/database/sqlc"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID
	Name        string
	Permissions []Permission
}

func (r *Role) FromModel(m *sqlc.Role) *Role {
	r.ID = m.ID
	r.Name = m.Name
	return r
}

func (r *Role) WithPermissions(ps []Permission) *Role {
	r.Permissions = ps
	return r
}

func FromGetRolesWithPermissionsToRoles(rps []sqlc.GetRolesWithPermissionsRow) []Role {
	roles := []Role{}
	for _, rp := range rps {
		var role *Role
		roleIndex := -1
		for i, r := range roles {
			if r.ID == rp.Role.ID {
				role = &r
				roleIndex = i
				break
			}
		}

		if role == nil {
			role = &Role{
				ID:   rp.Role.ID,
				Name: rp.Role.Name,
			}

			if rp.Permission.ID != uuid.Nil {
				role.Permissions = []Permission{Permission{ID: rp.Permission.ID, Name: rp.Permission.Name}}
			}

			roles = append(roles, *role)
		} else {
			if rp.Permission.ID != uuid.Nil {
				roles[roleIndex].Permissions = []Permission{Permission{ID: rp.Permission.ID, Name: rp.Permission.Name}}
			}
		}
	}

	return roles
}
