package swagger

//go:generate mockgen -destination mock_server.go -package swagger -mock_names Client=MockServer github.com/briggysmalls/detectordag/api/swagger/server Server

import (
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ExpectFunc func(*MockServer)

func TestValidRoutes(t *testing.T) {
	// Create requests to check
	testParams := []struct {
		method     string
		route      string
		expectFunc ExpectFunc
	}{
		{method: http.MethodPost, route: "/v1/auth", expectFunc: func(s *MockServer) {
			s.EXPECT().Auth(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodGet, route: "/v1/accounts/33b782d3-a2c8-40be-8aef-db5b44119bd5", expectFunc: func(s *MockServer) {
			s.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodGet, route: "/v1/accounts/f88948e6-5f93-4f11-8d58-15d48075069d/devices", expectFunc: func(s *MockServer) {
			s.EXPECT().GetDevices(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodPatch, route: "/v1/accounts/cfe7d5ed-826e-4e31-bb46-d62aa1cb58a7", expectFunc: func(s *MockServer) {
			s.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
		{method: http.MethodPatch, route: "/v1/devices/c0e94a1b-a835-4cc2-9574-642bea13805a", expectFunc: func(s *MockServer) {
			s.EXPECT().UpdateDevice(gomock.Any(), gomock.Any()).Do(setStatusOk)
		}},
	}

	for _, params := range testParams {
		// Create router
		server, router := createTestRouter(t)
		// Create a request
		r, err := http.NewRequest(params.method, params.route, nil)
		if err != nil {
			t.Fatal(err)
		}
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		w := httptest.NewRecorder()
		// Configure the expectations
		params.expectFunc(server)
		// Get the router to handle the request
		router.ServeHTTP(w, r)
		// Ensure we get a 200
		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v", status)
		}
	}
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
