package model

import (
	"open-fermentations/internal/database/sqlc"
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
