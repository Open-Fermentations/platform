package service

import (
	"errors"
	"open-fermentations/internal/model"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Login implements [Service].
func (s service) Login(username string, password string) (*model.User, error) {
	u, err := s.db.Querier().GetUserByUsernameWithPassword(s.ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, errors.Join(ErrInvalidCredentials{ID: u.ID, Username: username}, err)
	}

	user := new(model.User).FromUsernameWithPasswordRow(&u)

	return user, nil
}
