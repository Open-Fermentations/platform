package model

import (
	"open-fermentations/internal/database/sqlc"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_FromGetRolesWithPermissionsToRoles(t *testing.T) {
	t.Run("with only a role, should create the role", func(t *testing.T) {
		r := sqlc.Role{
			Name: "role1",
		}
		rps := []sqlc.GetRolesWithPermissionsRow{
			{Role: r},
		}

		actual := FromGetRolesWithPermissionsToRoles(rps)
		assert.EqualValues(t, []Role{{Name: r.Name}}, actual)
	})

	t.Run("with roles and permissions, should construct the roles and permissions", func(t *testing.T) {

		p1 := sqlc.Permission{
			ID:   uuid.New(),
			Name: "permission1",
		}
		p2 := sqlc.Permission{
			ID:   uuid.New(),
			Name: "permission2",
		}
		r1 := sqlc.Role{
			ID:   uuid.New(),
			Name: "role1",
		}
		r2 := sqlc.Role{
			ID:   uuid.New(),
			Name: "role2",
		}
		r3 := sqlc.Role{
			ID:   uuid.New(),
			Name: "role3",
		}

		rps := []sqlc.GetRolesWithPermissionsRow{
			{Role: r1, Permission: p1},
			{Role: r1, Permission: p2},
			{Role: r2, Permission: p1},
			{Role: r3},
		}

		actual := FromGetRolesWithPermissionsToRoles(rps)

		expected := []Role{{
			ID:          r1.ID,
			Name:        r1.Name,
			Permissions: []Permission{{ID: p1.ID, Name: p1.Name}, {ID: p2.ID, Name: p2.Name}},
		}, {
			ID:          r2.ID,
			Name:        r2.Name,
			Permissions: []Permission{{ID: p1.ID, Name: p1.Name}},
		}, {
			ID:          r3.ID,
			Name:        r3.Name,
			Permissions: []Permission{},
		}}

		assert.EqualValues(t, expected, actual)
	})
}
