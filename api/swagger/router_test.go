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
