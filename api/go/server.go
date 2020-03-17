package swagger

import (
	"fmt"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	JwtSecret   string `split_words:"true"`
	JwtDuration string `split_words:"true`
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
