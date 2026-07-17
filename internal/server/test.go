package server

import (
	mockdatabase "open-fermentations/internal/testing/mocks/database"
	mocksqlc "open-fermentations/internal/testing/mocks/database/sqlc"
	mockservice "open-fermentations/internal/testing/mocks/service"
	"testing"
)

type testContext struct {
	s     *Server
	mSvc  *mockservice.MockService
	mDb   *mockdatabase.MockService
	mSqlc *mocksqlc.MockQuerier
}

func setupContext(t *testing.T) *testContext {
	mockSqlc := mocksqlc.NewMockQuerier(t)
	mockDb := mockdatabase.NewMockService(t)
	mockSvc := mockservice.NewMockService(t)

	return &testContext{
		s:     &Server{db: mockDb, svc: mockSvc},
		mDb:   mockDb,
		mSvc:  mockSvc,
		mSqlc: mockSqlc,
	}
}

func testCase(test func(t *testing.T, c *testContext)) func(*testing.T) {
	return func(t *testing.T) {
		c := setupContext(t)
		test(t, c)
	}
}
