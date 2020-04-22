package swagger

//go:generate mockgen -destination mock_db.go -package swagger -mock_names Client=MockDBClient github.com/briggysmalls/detectordag/shared/database Client
//go:generate mockgen -destination mock_shadow.go -package swagger -mock_names Client=MockShadowClient github.com/briggysmalls/detectordag/shared/shadow Client
//go:generate mockgen -destination mock_email.go -package swagger -mock_names Client=MockEmailClient github.com/briggysmalls/detectordag/shared/email Client
//go:generate mockgen -destination mock_tokens.go -package swagger -mock_names Client=MockTokens github.com/briggysmalls/detectordag/api/swagger/tokens Tokens

import (
	"bytes"
	"github.com/briggysmalls/detectordag/api/swagger/server"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var testDuration time.Duration

func init() {
	var err error
	testDuration, err = time.ParseDuration("1ms")
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func createRealRouter(t *testing.T) (*MockDBClient, *MockShadowClient, *MockEmailClient, *MockTokens, *mux.Router) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock database
	db := NewMockDBClient(ctrl)
	// Create mock shadow
	shadow := NewMockShadowClient(ctrl)
	// Create mock email
	email := NewMockEmailClient(ctrl)
	// Create mock tokens
	tokens := NewMockTokens(ctrl)
	// Create real server
	s := server.New(server.Params{
		Db:     db,
		Shadow: shadow,
		Email:  email,
		Tokens: tokens,
	})
	// Create the new router
	return db, shadow, email, tokens, NewRouter(db, s, tokens)
}

func runHandler(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	// Create a recorder for the response
	rr := httptest.NewRecorder()
	// Ask the router to handle the request
	router.ServeHTTP(rr, req)
	return rr
}

func createRequest(t *testing.T, method, route string, body []byte) *http.Request {
	req, err := http.NewRequest(method, route, bytes.NewReader(body))
	assert.NoError(t, err)
	return req
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}
