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

		userDto := new(dto.UserDTO).FromModel(user)
		slog.Info("UserDTO", userDto.Slog())
		jsonResp, err := json.Marshal(userDto)
		if err != nil {
			slog.Error("unmarshalling user dto", slog.String("error", err.Error()), userDto.Slog())
			http.Error(w, "Failed to marshal user dto", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(jsonResp); err != nil {
			slog.Error("writing response", slog.String("error", err.Error()))
			return
		}
	}
}
