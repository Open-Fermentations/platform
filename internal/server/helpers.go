package server

import (
	"encoding/json"
	"io"
	"net/http"
	"open-fermentations/internal/model"
	"open-fermentations/internal/route"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func getUserId(r *http.Request) uuid.UUID {
	id := r.Context().Value(ContextUserIdKey).(uuid.UUID)

	return id
}

func readBody(r *http.Request, d interface{}) error {
	rawBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	return json.Unmarshal(rawBody, d)
}

func setContentTypeJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", route.ContentTypeJSON)
}

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
