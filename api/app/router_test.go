package app

//go:generate go run github.com/golang/mock/mockgen -destination mock_server.go -package app -mock_names Client=MockServer -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/api/app/server Server
//go:generate go run github.com/golang/mock/mockgen -destination mock_iot.go -package app -mock_names Client=MockIoTClient -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/shared/iot Client

import (
	"fmt"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiIzNTU4MUJGNC0zMkM4LTQ5MDgtODM3Ny0yRTZBMDIxRDNEMkIiLCJleHAiOjkyMjMzNzIwMzY4NTQ3NzU4MDcsImlzcyI6ImRldGVjdG9yZGFnIn0.CzyaCEIXlq1E0F89HR2Z9wbUn5gBDyQKTOCxTsX6iiQ"

type expectFunc func(*MockServer, *MockIoTClient, *MockTokens)

func TestValidRoutes(t *testing.T) {
	// Create requests to check
	tps := []struct {
		method     string
		route      string
		expectFunc expectFunc
	}{
		{method: http.MethodPost, route: "/v1/auth", expectFunc: func(s *MockServer, _ *MockIoTClient, _ *MockTokens) {
			// Expect the handler to be called
			s.EXPECT().Auth(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodGet, route: "/v1/accounts/33b782d3-a2c8-40be-8aef-db5b44119bd5", expectFunc: func(s *MockServer, _ *MockIoTClient, tokens *MockTokens) {
			// Expect the auth middleware to validate the token
			expectAuth(tokens, "33b782d3-a2c8-40be-8aef-db5b44119bd5")
			// Expect the handler to be called
			s.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodPatch, route: "/v1/accounts/cfe7d5ed-826e-4e31-bb46-d62aa1cb58a7", expectFunc: func(s *MockServer, _ *MockIoTClient, tokens *MockTokens) {
			// Expect the auth middleware to validate the token
			expectAuth(tokens, "cfe7d5ed-826e-4e31-bb46-d62aa1cb58a7")
			// Expect the handler to be called
			s.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodGet, route: "/v1/accounts/f88948e6-5f93-4f11-8d58-15d48075069d/devices", expectFunc: func(s *MockServer, _ *MockIoTClient, tokens *MockTokens) {
			// Expect the auth middleware to validate the token
			expectAuth(tokens, "f88948e6-5f93-4f11-8d58-15d48075069d")
			// Expect the handler to be called
			s.EXPECT().GetDevices(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodPatch, route: "/v1/devices/c0e94a1b-a835-4cc2-9574-642bea13805a", expectFunc: func(s *MockServer, i *MockIoTClient, tokens *MockTokens) {
			// Expect the auth middleware to get the device from database
			accountID := "f88948e6-5f93-4f11-8d58-15d48075069d"
			i.EXPECT().GetThing(gomock.Eq("c0e94a1b-a835-4cc2-9574-642bea13805a")).Return(&iot.Device{AccountId: accountID}, nil)
			// Expect the auth middleware to validate the token
			expectAuth(tokens, accountID)
			// Expect the handler to be called
			s.EXPECT().UpdateDevice(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
	}
	// Run the test iterations
	for _, params := range tps {
		log.Printf("Testing route (%s) '%s'", params.method, params.route)
		// Create a request
		r, err := http.NewRequest(params.method, params.route, nil)
		assert.NoError(t, err)
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testToken))
		// Run the request
		w := runTest(t, r, params.expectFunc)
		// Ensure we get a 200
		assert.Equal(t, http.StatusOK, w.StatusCode)
		// Assert Content-Type is set
		assert.Equal(t, "application/json; charset=UTF-8", getHeaderValue(t, w.Header, "Content-Type"))
		// Assert Access-Control-Allow-Origin is set
		assert.Equal(t, "*", getHeaderValue(t, w.Header, "Access-Control-Allow-Origin"))
	}
}

func TestOptionsRoutes(t *testing.T) {
	// Create requests to check
	tps := []struct {
		route   string
		methods []string
	}{
		{route: "/v1/auth"},
		{route: "/v1/accounts/33b782d3-a2c8-40be-8aef-db5b44119bd5"},
		{route: "/v1/accounts/f88948e6-5f93-4f11-8d58-15d48075069d/devices"},
		{route: "/v1/devices/c0e94a1b-a835-4cc2-9574-642bea13805a"},
	}
	// Run the test iterations
	for _, params := range tps {
		log.Printf("Testing route '%s'", params.route)
		// Create a request
		r, err := http.NewRequest(http.MethodOptions, params.route, nil)
		assert.NoError(t, err)
		// Run the request
		w := runTest(t, r, nil)
		// Ensure we get a 200
		assert.Equal(t, http.StatusOK, w.StatusCode)
		// Ensure we have the expected allowed headers
		allowedHeaders := strings.Split(getHeaderValue(t, w.Header, "Access-Control-Allow-Headers"), ",")
		assert.Contains(t, allowedHeaders, "Content-Type")
		assert.Contains(t, allowedHeaders, "Authorization")
		assert.Len(t, allowedHeaders, 2)
		// Ensure we have the expected methods
		allowedMethods := strings.Split(getHeaderValue(t, w.Header, "Access-Control-Allow-Headers"), ",")
		for _, method := range params.methods {
			assert.Contains(t, allowedMethods, method)
		}
		assert.Len(t, allowedMethods, len(allowedMethods))
	}
}

func runTest(t *testing.T, r *http.Request, expect expectFunc) *http.Response {
	// Create router
	i, tokens, server, router := createStubbedRouter(t)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	w := httptest.NewRecorder()
	// Configure the expectations
	if expect != nil {
		expect(server, i, tokens)
	}
	// Get the router to handle the request
	router.ServeHTTP(w, r)
	// Return the response for checking
	return w.Result()
}

func setStatusOk(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func expectAuth(tokens *MockTokens, accountID string) {
	// Expect the auth middleware to validate the token
	tokens.EXPECT().Validate(gomock.Eq(testToken)).Return(accountID, nil)
}

func createStubbedRouter(t *testing.T) (*MockIoTClient, *MockTokens, *MockServer, *mux.Router) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock database
	i := NewMockIoTClient(ctrl)
	// Create mock shadow
	s := NewMockServer(ctrl)
	// Create mock tokens
	tokens := NewMockTokens(ctrl)
	// Create the new router
	return i, tokens, s, NewRouter(i, s, tokens)
}

func getHeaderValue(t *testing.T, header http.Header, key string) string {
	assert.Contains(t, header, key)
	assert.Len(t, header[key], 1)
	return header[key][0]
}
