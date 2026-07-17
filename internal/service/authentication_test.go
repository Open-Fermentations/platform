package service

import (
	"errors"
	"open-fermentations/internal/database/sqlc"
	"open-fermentations/internal/env"
	mockdatabase "open-fermentations/internal/testing/mocks/database"
	mocksqlc "open-fermentations/internal/testing/mocks/database/sqlc"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	ErrMock = errors.New("some mock error")
)

type testContext struct {
	svc      Service
	mockDb   *mockdatabase.MockService
	mockSqlc *mocksqlc.MockQuerier
}

func setupContext(t *testing.T) *testContext {
	m := mockdatabase.NewMockService(t)
	mSqlc := mocksqlc.NewMockQuerier(t)

	m.EXPECT().Querier().Return(mSqlc)

	c := testContext{
		mockDb:   m,
		mockSqlc: mSqlc,
		svc:      New(t.Context(), &env.Env{}, m),
	}

	return &c
}

func (c *testContext) afterEach(t *testing.T) {
	// Set the service instance to nil in order to recreate it
	serviceInstance = nil
}

func testCase(test func(t *testing.T, c *testContext)) func(*testing.T) {
	return func(t *testing.T) {
		c := setupContext(t)
		// before each goes here
		defer c.afterEach(t)
		test(t, c)
	}
}

func Test_Login(t *testing.T) {
	mockUser := sqlc.GetUserByUsernameWithPasswordRow{
		Username: "MockUsername",
		// Password is 'admin' hashed
		Password: "$2a$10$5nmh/cOu.dzk05V7lfBqQua9FO6nG.aQTGTJQFB26DGMSMwp5FWxu",
	}

	t.Run("with db.GetUserByUsernameWithPassword returning an error",
		testCase(func(t *testing.T, c *testContext) {
			u := "non-existent-username"
			c.mockSqlc.EXPECT().
				GetUserByUsernameWithPassword(mock.Anything, u).
				Once().
				Return(sqlc.GetUserByUsernameWithPasswordRow{}, ErrMock)

			usr, err := c.svc.Login(u, "")
			assert.Nil(t, usr)
			assert.ErrorIs(t, err, ErrMock)
		}))

	t.Run("with an incorrect password, should return ErrInvalidCredentialsErr",
		testCase(func(t *testing.T, c *testContext) {
			c.mockSqlc.EXPECT().
				GetUserByUsernameWithPassword(mock.Anything, mockUser.Username).
				Once().
				Return(mockUser, nil)

			usr, err := c.svc.Login(mockUser.Username, "some incorrect password")
			assert.Nil(t, usr)
			var invalidCredentialsErr ErrInvalidCredentials
			assert.ErrorAs(t, err, &invalidCredentialsErr)
			assert.EqualValues(t, invalidCredentialsErr.Username, mockUser.Username)
		}))

	t.Run("with correct password, should return model.User",
		testCase(func(t *testing.T, c *testContext) {
			c.mockSqlc.EXPECT().
				GetUserByUsernameWithPassword(mock.Anything, mockUser.Username).
				Once().
				Return(mockUser, nil)

			usr, err := c.svc.Login(mockUser.Username, "admin")
			assert.Nil(t, err)
			assert.EqualValues(t, mockUser.Username, usr.Username)
		}))
}
