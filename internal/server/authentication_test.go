package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/env"
	"open-fermentations/internal/model"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testHandler(t *testing.T, secret string, statusCode int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		if _, err := w.Write([]byte(secret)); err != nil {
			t.Fatalf("could not write response in success case. Secret: %v, StatusCode: %v", secret, statusCode)
		}
	})
}

func TestLoginHandler(t *testing.T) {
	t.Run("with service login responding with a user",
		testCase(func(t *testing.T, c *testContext) {
			c.s.env = &env.Env{Jwt: env.JwtEnv{Key: "something"}, Cookie: env.CookieEnv{
				Secure:   false,
				Key:      "some-key",
				Duration: time.Second,
			}}
			c.mSvc.EXPECT().Login(mock.Anything, mock.Anything).
				Return(&model.User{Username: "admin"}, nil)
			server := httptest.NewServer(http.HandlerFunc(c.s.loginHandler))
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

func TestLogoutHandler(t *testing.T) {
	t.Run("sets cookie with blank value", testCase(func(t *testing.T, c *testContext) {
		c.s.env = &env.Env{Cookie: env.CookieEnv{Secure: false}}
		server := httptest.NewServer(http.HandlerFunc(c.s.logoutHandler))
		defer server.Close()

		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err.Error())
		}

		assert.EqualValues(t, http.StatusOK, resp.StatusCode)
		cookie := resp.Header.Get("Set-Cookie")
		assert.NotEmpty(t, cookie)

		parsedCookie := parseCookie(cookie)
		assert.EqualValues(t, "", parsedCookie["auth_token"])
	}))
}

func Test_authenticationMiddleware(t *testing.T) {
	e := &env.Env{
		Cookie: env.CookieEnv{
			Key: "some-key",
		},
		Jwt: env.JwtEnv{
			Key: "some-secure-key",
		}}
	t.Run("success case, should reach handler func", testCase(func(t *testing.T, c *testContext) {
		c.s.env = e
		successText := strconv.Itoa(time.Now().Nanosecond())
		server := httptest.NewServer(c.s.authenticationMiddleware(http.HandlerFunc(
			testHandler(t, successText, http.StatusOK))))
		defer server.Close()

		token, err := generateJwt([]byte(c.s.env.Jwt.Key), 30*time.Minute, &model.User{})
		if err != nil {
			t.Fatalf("could not create jwt: Err %v", err.Error())
		}
		cookie := generateCookie(c.s.env.Cookie.Key, token, time.Minute*30, false)

		req, _ := http.NewRequest("GET", server.URL, nil)
		req.Header.Set("Cookie", cookie.Name+"="+cookie.Value)

		client := server.Client()
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err.Error())
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("error reading body. Err: %v", err.Error())
		}

		assert.EqualValues(t, http.StatusOK, resp.StatusCode)
		assert.EqualValues(t, successText, strings.Trim(string(b), "\n"))
	}))

	t.Run("valid token signed by other key, should return 401 status code and not reach handler func",
		testCase(func(t *testing.T, c *testContext) {
			c.s.env = e
			successText := strconv.Itoa(time.Now().Nanosecond())

			server := httptest.NewServer(c.s.authenticationMiddleware(http.HandlerFunc(
				testHandler(t, successText, http.StatusOK))))
			defer server.Close()

			token, err := generateJwt([]byte("some-other-key"), 30*time.Minute, &model.User{})
			if err != nil {
				t.Fatalf("could not create jwt: Err %v", err.Error())
			}
			cookie := generateCookie(c.s.env.Cookie.Key, token, time.Minute*30, false)

			req, _ := http.NewRequest("GET", server.URL, nil)
			req.Header.Set("Cookie", cookie.Name+"="+cookie.Value)

			client := server.Client()
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("error making request to server. Err: %v", err.Error())
			}
			defer resp.Body.Close()

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("could not read body. Err: %v", err.Error())
			}

			assert.EqualValues(t, http.StatusUnauthorized, resp.StatusCode)
			assert.EqualValues(t, "unauthorized", strings.Trim(string(b), "\n"))
		}))

	t.Run("no cookie, should return 401 and not reach handler func", testCase(func(t *testing.T, c *testContext) {
		c.s.env = e
		successText := strconv.Itoa(time.Now().Nanosecond())
		server := httptest.NewServer(c.s.authenticationMiddleware(http.HandlerFunc(
			testHandler(t, successText, http.StatusOK))))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, nil)

		client := server.Client()
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err.Error())
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("could not read body. Err: %v", err.Error())
		}

		assert.EqualValues(t, http.StatusUnauthorized, resp.StatusCode)
		assert.EqualValues(t, "unauthorized", strings.Trim(string(b), "\n"))
	}))
}
