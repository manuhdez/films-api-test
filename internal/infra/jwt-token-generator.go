package infra

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const TokenDuration = time.Hour * 24

var (
	ErrParsingToken = errors.New("error parsing token")
	ErrInvalidToken = errors.New("invalid token")
)

type JWTGenerator struct {
	secretKey      string
	expirationDate time.Time
}

func NewJWTGenerator() JWTGenerator {
	secret := os.Getenv("JWT_SECRET_KEY")

	return JWTGenerator{
		secretKey:      secret,
		expirationDate: time.Now().Add(TokenDuration),
	}
}

func (g JWTGenerator) Generate(id string) (string, error) {
	claims := jwt.MapClaims{
		"exp":  g.expirationDate.Unix(),
		"user": id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(g.secretKey))
}

func (g JWTGenerator) Validate(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(g.secretKey), nil
	})

	if err != nil {
		return errors.Join(err, ErrParsingToken)
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}
