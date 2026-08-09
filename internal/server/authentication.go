package server

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"open-fermentations/internal/constants"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/logging"
	"open-fermentations/internal/model"
	"open-fermentations/internal/service"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	JWTIdKey          string = "id"
	JWTPermissionsKey string = "permissions"
	JWTRolesKey       string = "roles"
)

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var b dto.LoginBody
	if err := readBody(r, &b); err != nil {
		slog.Error("reading body", logging.Err(err))
		http.Error(w, FailedReadingBody, http.StatusBadRequest)
		return
	}

	user, err := s.svc.Login(r.Context(), b.Username, b.Password)
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

	userDto := new(dto.AuthenticatedUserDTO).FromModel(user)
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

			ctx, err = s.mapClaimsToContext(ctx, claims)
			if err != nil {
				slog.Error("mapping claims to context", logging.Err(err))
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

func addUserIdToContext(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {
	if rawId, ok := claims[JWTIdKey].(string); ok {
		id, err := uuid.Parse(rawId)
		if err != nil {
			return nil, err
		}

		return context.WithValue(ctx, constants.ContextUserIdKey, id), nil
	} else {
		return nil, ErrJWTClaim{Key: JWTIdKey}
	}
}

func (s *Server) addRolesToContext(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {
	switch v := claims[JWTRolesKey].(type) {
	case []any:
		roles := make([]model.Role, len(v))
		for i, item := range v {
			if strVal, ok := item.(string); ok {
				var role model.Role
				for _, roleItem := range s.roles {
					if roleItem.Name == strVal {
						role = roleItem
					}
				}
				roles[i] = role
			}
		}
		return context.WithValue(ctx, constants.ContextRolesKey, roles), nil
	default:
		return nil, ErrJWTClaim{Key: JWTRolesKey}
	}
}

func (s *Server) addPermissionsToContext(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {
	switch v := claims[JWTPermissionsKey].(type) {
	case []string:
		return context.WithValue(ctx, constants.ContextPermissionsKey, v), nil
	case []interface{}:
		strPermissions := make([]string, len(v))
		for i, p := range v {
			if strVal, ok := p.(string); ok {
				strPermissions[i] = strVal
			}
		}
		return context.WithValue(ctx, constants.ContextPermissionsKey, strPermissions), nil
	default:
		return nil, ErrJWTClaim{Key: JWTPermissionsKey}
	}
}

func (s *Server) mapClaimsToContext(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {
	var err error
	ctx, err = addUserIdToContext(ctx, claims)
	if err != nil {
		return nil, err
	}
	ctx, err = s.addRolesToContext(ctx, claims)
	if err != nil {
		return nil, err
	}
	ctx, err = s.addPermissionsToContext(ctx, claims)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}
