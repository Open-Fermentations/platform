package service

import (
	"fmt"
	"open-fermentations/internal/model"
)

// Login implements [Service].
func (s service) Login(username string, password string) (*model.User, error) {
	u, err := s.db.Queries().GetUserByUsernameWithPassword(s.ctx, username)
	if err != nil {
		return nil, fmt.Errorf("fetching user from db: %w", err)
	}

	user := new(model.User).FromUsernameWithPasswordRow(&u)

	return user, nil
}
