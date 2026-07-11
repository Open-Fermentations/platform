package server

import (
	"context"
	"encoding/json"
	"net/http"
	"open-fermentations/internal/dto"
)

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	user, err := s.db.Queries().GetUserByUsername(context.Background(), "admin")
	if err != nil {
		panic(err)
	}

	jsonResp, err := json.Marshal(new(dto.UserDTO).FromModel(&user))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		panic(err)
	}
}
