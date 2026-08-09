package dto

import (
	"log/slog"
	"open-fermentations/internal/logging"
	"open-fermentations/internal/model"
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

func (u *UserDTO) FromModel(m *model.User) *UserDTO {
	u.ID = m.ID
	u.Username = m.Username
	u.Created = m.Created
	u.Modified = m.Modified

	return u
}

func (u UserDTO) Slog() []any {
	return []any{slog.Group(
		"user",
		slog.String("id", u.ID.String()),
		slog.String("username", u.Username),
		slog.Time("created", u.Created),
		slog.Time("modified", u.Modified),
	)}
}

var _ logging.Slog = UserDTO{}

type AuthenticatedUserDTO struct {
	User        *UserDTO        `json:"user"`
	Roles       []RoleDTO       `json:"roles"`
	Permissions []PermissionDTO `json:"permissions"`
}

func (a *AuthenticatedUserDTO) FromModel(m *model.AuthenticatedUser) *AuthenticatedUserDTO {
	a.User = new(UserDTO).FromModel(&m.User)

	a.Roles = []RoleDTO{}
	for _, r := range m.Roles {
		a.Roles = append(a.Roles, *new(RoleDTO).FromModel(&r))
	}

	a.Permissions = []PermissionDTO{}
	for _, p := range m.Permissions {
		a.Permissions = append(a.Permissions, *new(PermissionDTO).FromModel(&p))
	}

	return a
}
