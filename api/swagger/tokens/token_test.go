package tokens

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	timeLeeway, err := time.ParseDuration("10ns")
	assert.NoError(t, err)
	// Create some test inputs
	testParams := []struct {
		secret    string
		duration  string
		accountID string
		error     error
	}{
		{secret: "mysecret", duration: "2h", accountID: "35581BF4-32C8-4908-8377-2E6A021D3D2B"},
		{secret: "anothersecret", duration: "1m", accountID: "22222222-32C8-4908-8377-2E6A021D3D2B"},
	}
	for _, params := range testParams {
		// Create a tokens
		duration, err := time.ParseDuration("10s")
		assert.NoError(t, err)
		tokens := New(params.secret, duration)
		// Create a token
		ss, err := tokens.Create(params.accountID)
		// Check the error
		assert.Equal(t, params.error, err)
		// Parse the token contents
		token, err := jwt.ParseWithClaims(ss, &CustomAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret
			return []byte(params.secret), nil
		})
		// Check the claims
		if claims, ok := token.Claims.(*CustomAuthClaims); ok && token.Valid {
			// Confirm the accounts match
			if claims.AccountId != params.accountID {
				t.Fatalf("Token did not save correct account ID")
			}
			// Confirm the expiry time
			dur, err := time.ParseDuration(params.duration)
			if err != nil {
				t.Fatalf(err.Error())
			}
			exp := time.Unix(claims.ExpiresAt, 0)
			if math.Abs(float64(time.Until(exp)-dur)) < float64(timeLeeway) {
				t.Fatalf("Token did not save correct duration")
			}
		} else {
			t.Fatalf(err.Error())
		}
	}
}

func TestCheckValid(t *testing.T) {
	// Create a dummy expiry duration
	duration, err := time.ParseDuration("10s")
	assert.NoError(t, err)
	// Create some test inputs
	testParams := []struct {
		secret    string
		token     string
		now       time.Time
		accountID string
		error     error
	}{
		{
			// Valid token
			secret:    "mysecret",
			accountID: "35581BF4-32C8-4908-8377-2E6A021D3D2B",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiIzNTU4MUJGNC0zMkM4LTQ5MDgtODM3Ny0yRTZBMDIxRDNEMkIiLCJleHAiOjE1ODQ3OTk0NzQsImlzcyI6ImRldGVjdG9yZGFnIn0.qqMDypPk5BT1dz_8KT6S9eNLABWcYIfnaRr_BroisKo",
			now:       createTime(t, "2020/03/21 12:06:00"),
			error:     nil,
		},
		{
			// Expired token
			secret:    "mysecret",
			accountID: "35581BF4-32C8-4908-8377-2E6A021D3D2B",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiIzNTU4MUJGNC0zMkM4LTQ5MDgtODM3Ny0yRTZBMDIxRDNEMkIiLCJleHAiOjE1ODQ3OTk0NzQsImlzcyI6ImRldGVjdG9yZGFnIn0.qqMDypPk5BT1dz_8KT6S9eNLABWcYIfnaRr_BroisKo",
			now:       createTime(t, "2020/03/22 12:06:00"),
			error:     ErrBadToken,
		},
		{
			// Missing token
			secret:    "mysecret",
			accountID: "35581BF4-32C8-4908-8377-2E6A021D3D2B",
			token:     "",
			now:       createTime(t, "2020/03/22 12:06:00"),
			error:     ErrBadToken,
		},
		{
			// Token with RS256 signing algorithm
			secret:    "mysecret",
			accountID: "35581BF4-32C8-4908-8377-2E6A021D3D2B",
			token:     "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.POstGetfAytaZS82wHcjoTyoqhMyxXiWdR7Nn7A29DNSl0EiXLdwJ6xC6AfgZWF1bOsS_TuYI3OG85AmiExREkrS6tDfTQ2B3WXlrr-wp5AokiRbz3_oB4OxG-W9KcEEbDRcZc0nH3L7LzYptiy1PtAylQGxHTWZXtGz4ht0bAecBgmpdgXMguEIcoqPJ1n3pIWk_dUZegpqx0Lka21H6XxUTxiy8OcaarA8zdnPUnV6AmNP3ecFawIFYdvJB_cm-GvpCSbr8G8y_Mllj8f4x9nBH8pQux89_6gUY618iYv7tuPWBFfEbLxtF2pZS6YC1aSfLQxeNe8djT9YjpvRZA",
			now:       createTime(t, "2020/03/22 12:06:00"),
			error:     ErrBadToken,
		},
	}
	for _, params := range testParams {
		// Create a tokens
		tokens := New(params.secret, duration)
		// Check if the token authorises the supplied account
		at(params.now, func() {
			accountID, err := tokens.Validate(params.token)
			assert.Equal(t, params.error, err)
			if err == nil {
				// We weren't expecting an error
				assert.Equal(t, params.accountID, accountID)
				return
			}
		})
	}
}

// Override time value for tests.  Restore default value after.
func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}
