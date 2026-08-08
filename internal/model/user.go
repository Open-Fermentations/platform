package model

import (
	"open-fermentations/internal/database/sqlc"
	"slices"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

func (u *User) FromModel(m *sqlc.User) *User {
	u.ID = m.ID
	u.Username = m.Username
	u.Password = m.Password
	u.Created = m.Created
	u.Modified = m.Modified

	return u
}

type AuthenticatedUser struct {
	User
	Permissions []Permission
	Roles       []Role
}

func (a *AuthenticatedUser) FromUserQuery(q []sqlc.GetUserByUsernameWithPasswordAndRolesAndPermissionsRow) (*AuthenticatedUser, error) {
	if len(q) == 0 {
		return nil, ErrNoElements
	}

	a.User = *new(User).FromModel(&q[0].User)

	for i, urp := range q {
		if urp.Role.ID == uuid.Nil && slices.ContainsFunc(a.Permissions, func(p Permission) bool {
			return urp.Permission.ID == p.ID
		}) == false {
			a.Permissions = append(a.Permissions, *new(Permission).FromModel(&q[i].Permission))
			continue
		}

		var role *Role
		roleIndex := -1
		for ri, r := range a.Roles {
			if urp.Role.ID == r.ID {
				role = &r
				roleIndex = ri
				break
			}
		}
		if role == nil {
			role = &Role{ID: urp.Role.ID, Name: urp.Role.Name}
			if urp.Permission.ID != uuid.Nil {
				role.Permissions = append(role.Permissions, Permission{ID: urp.Permission.ID, Name: urp.Permission.Name})
			}
			a.Roles = append(a.Roles, *role)
		} else {
			if urp.Permission.ID != uuid.Nil && slices.ContainsFunc(role.Permissions, func(p Permission) bool {
				return urp.Permission.ID == p.ID
			}) == false {
				a.Roles[roleIndex].Permissions = append(role.Permissions, Permission{ID: urp.Permission.ID, Name: urp.Permission.Name})
			}
		}
	}

	return a, nil
}
