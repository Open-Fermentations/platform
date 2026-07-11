package server

import (
	"context"
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

		// Fetch user by the username provided in the body
		user, err := s.db.Queries().GetUserByUsernameWithPassword(context.Background(), b.Username)
		if err != nil {
			slog.Error("failed to find user", slog.String("error", err.Error()))
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// NOTE: You must compare the password here with a hashed version stored in the database.
		// For now, we just check if the retrieved user model has a dummy password match for demonstration.
		if user.Password != b.Password { // Replace with actual password comparison logic!
			slog.Info("login failed: incorrect password", slog.String("user", b.Username))
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
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
}
