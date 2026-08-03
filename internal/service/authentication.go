package service

import (
	"errors"
	"open-fermentations/internal/model"
	"slices"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Login implements [Service].
func (s service) Login(username string, password string) (*model.AuthenticatedUser, error) {
	userRolePermissions, err := s.db.Querier().GetUserByUsernameWithPasswordAndRolesAndPermissions(s.ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(userRolePermissions) == 0 {
		return nil, pgx.ErrNoRows
	}

	authedUser := model.AuthenticatedUser{
		User:        *new(model.User).FromModel(&userRolePermissions[0].User),
		Permissions: []model.Permission{},
		Roles:       []model.Role{},
	}
	for i, urp := range userRolePermissions {

		if slices.ContainsFunc(authedUser.Permissions, func(p model.Permission) bool {
			return urp.Permission.ID == p.ID
		}) == false {
			authedUser.Permissions = append(authedUser.Permissions, *new(model.Permission).FromModel(&userRolePermissions[i].Permission))
		}

		if slices.ContainsFunc(authedUser.Roles, func(r model.Role) bool {
			return urp.Role.ID == r.ID
		}) == false {
			authedUser.Roles = append(authedUser.Roles, *new(model.Role).FromModel(&userRolePermissions[i].Role))
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(authedUser.Password), []byte(password))
	if err != nil {
		return nil, errors.Join(ErrInvalidCredentials{ID: authedUser.ID, Username: username}, err)
	}
	return &authedUser, nil
}
