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
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateJwt(key []byte, u *model.User) (string, error) {
	// TODO: add expiry time in token as well
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "open-fermentations", // TODO get this from env
		"sub": u.Username,
		"id":  u.ID.String(),
	})

	return t.SignedString(key)
}

func parseCookie(cookie string) map[string]string {
	cookieParts := strings.Split(cookie, ";")
	parsedCookie := map[string]string{}
	for _, part := range cookieParts {
		keyVal := strings.Split(part, "=")
		if len(keyVal) == 2 {
			parsedCookie[strings.Trim(keyVal[0], " ")] = keyVal[1]
		} else {
			parsedCookie[part] = ""
		}
	}

	return parsedCookie
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
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

	cookie := &http.Cookie{
		Name:     s.env.Cookie.Key,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(s.env.Cookie.Duration),
		HttpOnly: true,
		Secure:   s.env.Cookie.Secure,
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResp); err != nil {
		slog.Error("writing response", slog.String("error", err.Error()))
		return
	}
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   s.env.Cookie.Secure,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("")); err != nil {
		slog.Error("writing response", slog.String("error", err.Error()))
		return
	}
}

func (s *Server) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := r.Header.Get("cookie")

		parsedCookie := parseCookie(cookie)

		tokenString := parsedCookie[s.env.Cookie.Key]

		token, err := jwt.Parse(tokenString, func(*jwt.Token) (any, error) {
			return []byte(s.env.Jwt.Key), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			slog.Error("parsing jwt token from cookie",
				logging.Err(err),
			)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// TODO: validate expiry time of token

		_ = token

		next.ServeHTTP(w, r)
	})
}
