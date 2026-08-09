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
	roles := []*Role{}
	for _, rp := range rps {
		var role *Role
		for _, r := range roles {
			if r.ID == rp.Role.ID {
				role = r
				break
			}
		}

		if role == nil {
			role = &Role{ID: rp.Role.ID, Name: rp.Role.Name, Permissions: []Permission{}}
			roles = append(roles, role)
		}

		if rp.Permission.ID != uuid.Nil {
			var permission *Permission
			for _, p := range role.Permissions {
				if p.ID == rp.Permission.ID {
					permission = &p
				}
			}
			if permission == nil {
				role.Permissions = append(role.Permissions, Permission{ID: rp.Permission.ID, Name: rp.Permission.Name})
			}
		}
	}

	result := make([]Role, len(roles))
	for i, r := range roles {
		result[i] = *r
	}

	return result
}
