package server

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/logging"
	"open-fermentations/internal/model"
	"open-fermentations/internal/service"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body", slog.String("error", err.Error()))
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var b dto.LoginBody

	if err := json.Unmarshal(rawBody, &b); err != nil {
		slog.Error("unmarshalling login body", logging.Err(err))
		http.Error(w, "Failed to unmarshal login body", http.StatusBadRequest)
		return
	}

	user, err := s.svc.Login(b.Username, b.Password)
	if err != nil {

		if errors.Is(err, service.ErrInvalidCredentials{}) {
			slog.Error("We got our error we were looking for")
		}

		var invalidCredentialError service.ErrInvalidCredentials
		if errors.As(err, &invalidCredentialError) {
			slog.Error("login", invalidCredentialError.SlogErr(err)...)
		} else {
			slog.Error("authenticating user", logging.Err(err))
		}
		http.Error(w, "Failed to authenticate user", http.StatusUnauthorized)
		return
	}

	userDto := new(dto.UserDTO).FromModel(user)
	jsonResp, err := json.Marshal(userDto)
	if err != nil {
		slog.Error("unmarshalling user dto", []any{logging.Err(err), userDto.Slog()}...)
		http.Error(w, "Failed to marshal user dto", http.StatusInternalServerError)
		return
	}

	token, err := generateJwt([]byte(s.env.Jwt.Key), user)
	if err != nil {
		panic(err)
	}

	expiry := time.Now().Add(30 * time.Minute)

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Expires:  expiry,
		HttpOnly: true,
		Secure:   s.env.CookieSecure,
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResp); err != nil {
		slog.Error("writing response", slog.String("error", err.Error()))
		return
	}
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   s.env.CookieSecure,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("")); err != nil {
		slog.Error("writing response", slog.String("error", err.Error()))
		return
	}
}

func generateJwt(key []byte, u *model.User) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "open-fermentations", // TODO get this from env
		"sub": u.Username,
		"id":  u.ID.String(),
	})

	return t.SignedString(key)
}
