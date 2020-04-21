package tokens

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

const (
	issuer = "detectordag"
)

var (
	ErrUnexpectedSigningMethod = errors.New("Cannot parse JWT claims")
	ErrCannotParseClaims       = errors.New("Cannot parse JWT claims")
	ErrInvalidJWT              = errors.New("JWT is invalid")
)

type Tokens interface {
	Create(accountID string) (string, error)
	Validate(token string) (string, error)
}

type tokens struct {
	secret   string
	duration time.Duration
}

type CustomAuthClaims struct {
	AccountId string `json:"accountId"`
	jwt.StandardClaims
}

func New(secret string, duration time.Duration) Tokens {
	return &tokens{
		secret:   secret,
		duration: duration,
	}
}

func (t *tokens) Create(accountID string) (string, error) {
	// Create the Claims
	claims := CustomAuthClaims{
		accountID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(t.duration).Unix(),
			Issuer:    issuer,
		},
	}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.secret))
}

// Validate checks that the provided token is valid
func (t *tokens) Validate(tokenString string) (string, error) {
	// Parse takes the token string and a function for looking up the key.
	token, err := jwt.ParseWithClaims(tokenString, &CustomAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, ErrUnexpectedSigningMethod
		}
		// Return our secret
		return []byte(t.secret), nil
	})
	if err != nil {
		return "", err
	}
	// Check the token contents
	claims, ok := token.Claims.(*CustomAuthClaims)
	if !ok {
		return "", ErrCannotParseClaims
	}
	if !token.Valid {
		return "", ErrInvalidJWT
	}
	return claims.AccountId, nil
}
