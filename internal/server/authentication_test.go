package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/env"
	"open-fermentations/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginHandler(t *testing.T) {
	t.Run("with service login responding with a user",
		testCase(func(t *testing.T, c *testContext) {
			c.s.env = &env.Env{Jwt: env.JwtEnv{Key: "something"}}
			c.mSvc.EXPECT().Login(mock.Anything, mock.Anything).
				Return(&model.User{Username: "admin"}, nil)
			server := httptest.NewServer(http.HandlerFunc(c.s.LoginHandler))
			defer server.Close()

			usr := dto.LoginBody{
				Username: "admin",
				Password: "admin",
			}
			json, _ := json.Marshal(usr)

			resp, err := http.Post(server.URL, ContentTypeJSON, bytes.NewReader(json))
			if err != nil {
				t.Fatalf("error making reuqest to server. Err: %v", err.Error())
			}
			defer resp.Body.Close()

			assert.EqualValues(t, http.StatusOK, resp.StatusCode)
			cookie := resp.Header.Get("Set-Cookie")
			assert.NotEmpty(t, cookie)
		}))
}
