package server

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/logging"
	"open-fermentations/internal/service"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var b dto.LoginBody
	if err := readBody(r, &b); err != nil {
		slog.Error("reading body", logging.Err(err))
		http.Error(w, FailedReadingBody, http.StatusBadRequest)
		return
	}

	user, err := s.svc.Login(b.Username, b.Password)
	if err != nil {
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
		http.Error(w, FailedMarshalling, http.StatusInternalServerError)
		return
	}

	token, err := generateJwt([]byte(s.env.Jwt.Key), s.env.Jwt.Expiry, user)
	if err != nil {
		panic(err)
	}

	cookie := generateCookie(s.env.Cookie.Key, token, s.env.Cookie.Duration, s.env.Cookie.Secure)
	http.SetCookie(w, cookie)

	setContentTypeJson(w)
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
			http.Error(w, Unauthorised, http.StatusUnauthorized)
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

				http.Error(w, Unauthorised, http.StatusUnauthorized)
				return
			}
			if rawId, ok := claims["id"].(string); ok {
				id, err := uuid.Parse(rawId)
				if err != nil {
					slog.Error("could not parse user id from claims", logging.Err(err), slog.String("id", rawId))
					http.Error(w, Unauthorised, http.StatusUnauthorized)
					return
				}
				ctx = context.WithValue(ctx, ContextUserIdKey, id)
			} else {
				slog.Error("could not get 'id' claim from token")
				http.Error(w, Unauthorised, http.StatusUnauthorized)
				return
			}
		} else {
			slog.Error("could not get claims from token")
			http.Error(w, Unauthorised, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
