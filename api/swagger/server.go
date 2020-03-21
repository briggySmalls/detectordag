package swagger

import (
	"errors"
	"fmt"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/kelseyhightower/envconfig"
	"net/http"
	"strings"
	"time"
)

const (
	AuthenticationHeaderPrefix = "Bearer "
)

var (
	ErrNoAuthHeader           = errors.New("Authorization header not set")
	ErrMalformattedAuthHeader = errors.New("Authorization header badly formed")
)

type Config struct {
	JwtSecret      string `split_words:"true"`
	JwtDuration    string `split_words:"true`
	ShadowEndpoint string `split_words:"true"`
}

type server struct {
	db     database.Client
	config Config
}

type CustomAuthClaims struct {
	AccountId string `json:"accountId"`
	jwt.StandardClaims
}

func (c *Config) ParseDuration() (time.Duration, error) {
	return time.ParseDuration(c.JwtDuration)
}

func NewConfig() (*Config, error) {
	// Load config
	var c Config
	var err error
	err = envconfig.Process("detectordag", &c)
	if err != nil {
		return nil, err
	}
	// Ensure duration is valid
	dur, err := c.ParseDuration()
	if err != nil {
		return nil, err
	}
	if dur.Seconds() < 1 {
		return nil, fmt.Errorf("JWT expiry duration insufficient: %f", dur.Seconds())
	}
	return &c, nil
}

func (s *server) createToken(accountId string) (string, error) {
	// Get the duration tokens should be alive for
	dur, err := s.config.ParseDuration()
	if err != nil {
		return "", err
	}
	// Create the Claims
	claims := CustomAuthClaims{
		accountId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(dur).Unix(),
			Issuer:    issuer,
		},
	}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JwtSecret))
}

// checkAuthorized checks that the token authorises access to the specified account
func (s *server) checkAuthorized(tokenString, accountId string) error {
	// Parse takes the token string and a function for looking up the key.
	token, err := jwt.ParseWithClaims(tokenString, &CustomAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return our secret
		return []byte(s.config.JwtSecret), nil
	})
	if err != nil {
		return err
	}
	// Check the token contents
	if claims, ok := token.Claims.(*CustomAuthClaims); ok && token.Valid {
		// Confirm the account IDs match
		if claims.AccountId != accountId {
			return fmt.Errorf("Not authorized to access account: %s", accountId)
		}
		// The token authorises access
		return nil
	}
	return err
}

func (s *server) getToken(header *http.Header) (string, error) {
	// Check the auth header is set
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}
	// Ensure we've been given a JWT how we expect
	if !strings.HasPrefix(authHeader, AuthenticationHeaderPrefix) {
		return "", ErrMalformattedAuthHeader
	}
	// Return the token
	return strings.TrimPrefix(authHeader, AuthenticationHeaderPrefix), nil
}
