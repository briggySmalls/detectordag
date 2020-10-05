package app

<<<<<<< HEAD:api/swagger/utils_test.go
//go:generate go run github.com/golang/mock/mockgen -destination mock_db.go -package swagger -mock_names Client=MockDBClient -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/shared/database Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_shadow.go -package swagger -mock_names Client=MockShadowClient -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/shared/shadow Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_email.go -package swagger -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/shared/email Verifier
//go:generate go run github.com/golang/mock/mockgen -destination mock_iot.go -package swagger -mock_names Client=MockIoTClient -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/shared/iot Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_tokens.go -package swagger -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/api/swagger/tokens Tokens
=======
//go:generate go run github.com/golang/mock/mockgen -destination mock_db.go -package app -mock_names Client=MockDBClient -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/shared/database Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_shadow.go -package app -mock_names Client=MockShadowClient -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/shared/shadow Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_email.go -package app -mock_names Client=MockEmailClient -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/shared/email Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_iot.go -package app -mock_names Client=MockIoTClient -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/shared/iot Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_tokens.go -package app -mock_names Client=MockTokens -self_package github.com/briggysmalls/detectordag/api github.com/briggysmalls/detectordag/api/app/tokens Tokens
>>>>>>> Move modules out of 'swagger' into 'app':api/app/test_utils.go

import (
	"bytes"
	"github.com/briggysmalls/detectordag/api/app/server"
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

func createRealRouter(t *testing.T) (*MockDBClient, *MockShadowClient, *MockVerifier, *MockIoTClient, *MockTokens, *mux.Router) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock database
	db := NewMockDBClient(ctrl)
	// Create mock shadow
	shadow := NewMockShadowClient(ctrl)
	// Create mock email
	email := NewMockVerifier(ctrl)
	// Create mock iot
	iot := NewMockIoTClient(ctrl)
	// Create mock tokens
	tokens := NewMockTokens(ctrl)
	// Create real server
	s := server.New(db, shadow, email, iot, tokens)
	// Create the new router
	return db, shadow, email, iot, tokens, NewRouter(iot, s, tokens)
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
