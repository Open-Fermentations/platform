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
		if slices.ContainsFunc(a.Permissions, func(p Permission) bool {
			return urp.Permission.ID == p.ID
		}) == false {
			a.Permissions = append(a.Permissions, *new(Permission).FromModel(&q[i].Permission))
		}

		if slices.ContainsFunc(a.Roles, func(r Role) bool {
			return urp.Role.ID == r.ID
		}) == false {
			a.Roles = append(a.Roles, *new(Role).FromModel(&q[i].Role))
		}
	}

	return a, nil
}
