package swagger

import (
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	JwtSecret string `split_words:"true"`
}

type server struct {
	db     database.Client
	config Config
}

type CustomAuthClaims struct {
	AccountId string `json:"accountId"`
	jwt.StandardClaims
}

func NewConfig() (*Config, error) {
	// Load config
	var c Config
	err := envconfig.Process("detectordag", &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *server) createToken(accountId string) (string, error) {
	// Create the Claims
	claims := CustomAuthClaims{
		accountId,
		jwt.StandardClaims{
			ExpiresAt: expiryDuration,
			Issuer:    issuer,
		},
	}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JwtSecret))
}
