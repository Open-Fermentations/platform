package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"open-fermentations/internal/dto"
)

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("failed to read request body", slog.String("error", err.Error()))
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		var b LoginBody

		if err := json.Unmarshal(rawBody, &b); err != nil {
			slog.Error("unmarshalling login body", slog.String("error", err.Error()))
			http.Error(w, "Failed to unmarshal login body", http.StatusBadRequest)
			return
		}

		user, err := s.svc.Login(b.Username, b.Password)
		if err != nil {
			slog.Error("authenticating user", slog.String("error", err.Error()))
			http.Error(w, "Failed to authenticate user", http.StatusUnauthorized)
			return
		}

		jsonResp, err := json.Marshal(new(dto.UserDTO).FromModel(user))
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(jsonResp); err != nil {
			panic(err)
		}
	}
}
