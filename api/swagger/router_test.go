package swagger

//go:generate mockgen -destination mock_server.go -package swagger -mock_names Client=MockServer github.com/briggysmalls/detectordag/api/swagger/server Server

import (
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type expectFunc func(*MockServer)

func TestValidRoutes(t *testing.T) {
	// Create requests to check
	tps := []struct {
		method     string
		route      string
		expectFunc expectFunc
	}{
		{method: http.MethodPost, route: "/v1/auth", expectFunc: func(s *MockServer) {
			s.EXPECT().Auth(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodGet, route: "/v1/accounts/33b782d3-a2c8-40be-8aef-db5b44119bd5", expectFunc: func(s *MockServer) {
			s.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodPatch, route: "/v1/accounts/cfe7d5ed-826e-4e31-bb46-d62aa1cb58a7", expectFunc: func(s *MockServer) {
			s.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodGet, route: "/v1/accounts/f88948e6-5f93-4f11-8d58-15d48075069d/devices", expectFunc: func(s *MockServer) {
			s.EXPECT().GetDevices(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodPatch, route: "/v1/devices/c0e94a1b-a835-4cc2-9574-642bea13805a", expectFunc: func(s *MockServer) {
			s.EXPECT().UpdateDevice(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
	}
	// Run the test iterations
	for _, params := range tps {
		// Create a request
		r, err := http.NewRequest(params.method, params.route, nil)
		assert.NoError(t, err)
		// Run the request
		w := runTest(t, r, params.expectFunc)
		// Ensure we get a 200
		assert.Equal(t, http.StatusOK, w.StatusCode)
		// Assert Content-Type is set
		assert.Equal(t, getHeaderValue(t, w.Header, "Content-Type"), "application/json; charset=UTF-8")
		// Assert Access-Control-Allow-Origin is set
		assert.Equal(t, getHeaderValue(t, w.Header, "Access-Control-Allow-Origin"), "*")
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
	server, router := createTestRouter(t)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	w := httptest.NewRecorder()
	// Configure the expectations
	if expect != nil {
		expect(server)
	}
	// Get the router to handle the request
	router.ServeHTTP(w, r)
	// Return the response for checking
	return w.Result()
}

func setStatusOk(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func createTestRouter(t *testing.T) (*MockServer, *mux.Router) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock server
	s := NewMockServer(ctrl)
	return s, NewRouter(s)
}

func getHeaderValue(t *testing.T, header http.Header, key string) string {
	assert.Contains(t, header, key)
	assert.Len(t, header[key], 1)
	return header[key][0]
}
