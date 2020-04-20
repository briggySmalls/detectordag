package server

//go:generate mockgen -destination mock_db.go -package server -mock_names Client=MockDBClient github.com/briggysmalls/detectordag/shared/database Client
//go:generate mockgen -destination mock_shadow.go -package server -mock_names Client=MockShadowClient github.com/briggysmalls/detectordag/shared/shadow Client

import (
	"bytes"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createMocks(t *testing.T) (*MockDBClient, *MockShadowClient) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock database client
	mockDb := NewMockDBClient(ctrl)
	// Create mock shadow client
	mockShadow := NewMockShadowClient(ctrl)
	return mockDb, mockShadow
}

func runHandler(h func(http.ResponseWriter, *http.Request), req *http.Request) *httptest.ResponseRecorder {
	// Run the handler using test code
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h)
	handler.ServeHTTP(rr, req)
	return rr
}

func createRequest(t *testing.T, method, route string, body []byte) *http.Request {
	req, err := http.NewRequest(method, route, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func assertStatus(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expected)
	}
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	if err != nil {
		t.Fatalf(err.Error())
	}
	return tme
}

// Override time value for tests.  Restore default value after.
func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}
