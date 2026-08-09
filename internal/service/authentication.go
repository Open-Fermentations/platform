package service

import (
	"context"
	"errors"
	"open-fermentations/internal/model"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Login implements [Service].
func (s service) Login(ctx context.Context, username string, password string) (*model.AuthenticatedUser, error) {
	userRolePermissions, err := s.db.Querier().GetUserByUsernameWithPasswordAndRolesAndPermissions(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(userRolePermissions) == 0 {
		return nil, pgx.ErrNoRows
	}

	authedUser, err := new(model.AuthenticatedUser).FromUserQuery(userRolePermissions)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(authedUser.Password), []byte(password))
	if err != nil {
		return nil, errors.Join(ErrInvalidCredentials{ID: authedUser.ID, Username: username}, err)
	}
	return authedUser, nil
}
