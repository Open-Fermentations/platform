package dto

import (
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
