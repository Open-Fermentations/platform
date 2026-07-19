package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/logging"
	"open-fermentations/internal/model"
	"open-fermentations/internal/route"
	"open-fermentations/internal/service"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func generateJwt(key []byte, exp time.Duration, u *model.User) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "open-fermentations", // TODO get this from env
		"sub": u.Username,
		"id":  u.ID.String(),
		"exp": time.Now().Add(exp).Unix(),
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

func generateCookie(key, token string, dur time.Duration, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     key,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(dur),
		HttpOnly: true,
		Secure:   secure,
	}
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body in loginHandler",
			logging.Err(err),
		)
		http.Error(w, BadBodyRead, http.StatusBadRequest)
		return
	}

	var b dto.LoginBody

	if err := json.Unmarshal(rawBody, &b); err != nil {
		slog.Error("unmarshalling login body", logging.Err(err))
		http.Error(w, FailedToMarshall, http.StatusBadRequest)
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
		slog.Error("marshalling user dto", logging.Err(err))
		http.Error(w, "Failed to marshal user dto", http.StatusInternalServerError)
		return
	}

	token, err := generateJwt([]byte(s.env.Jwt.Key), s.env.Jwt.Expiry, user)
	if err != nil {
		panic(err)
	}

	cookie := generateCookie(s.env.Cookie.Key, token, s.env.Cookie.Duration, s.env.Cookie.Secure)
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", route.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResp); err != nil {
		slog.Error("writing response", logging.Err(err))
		return
	}
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     s.env.Cookie.Key,
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
		ctx := r.Context()
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

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			err := validateClaims(claims)
			if err != nil {
				var jwtExpiredErr ErrJWTExpired
				if errors.As(err, &jwtExpiredErr) {
					slog.Error("jwt expired", jwtExpiredErr.Slog()...)
				} else {
					slog.Error("validating claims", logging.Err(err))
				}

				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if rawId, ok := claims["id"].(string); ok {
				id, err := uuid.Parse(rawId)
				if err != nil {
					slog.Error("could not parse user id from claims", logging.Err(err), slog.String("id", rawId))
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
				ctx = context.WithValue(ctx, ContextUserIdKey, id)
			} else {
				slog.Error("could not get 'id' claim from token")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		} else {
			slog.Error("could not get claims from token")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateClaims(claims jwt.Claims) error {
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return err
	}

	n := time.Now()
	if exp.Compare(n) <= 0 {
		return ErrJWTExpired{Exp: exp.Time, Now: n}
	}

	return nil
}
