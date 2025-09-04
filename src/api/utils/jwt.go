package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	ID          int64    `json:"id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	FirstLogin  bool     `json:"first_login"`
	jwt.RegisteredClaims
}

// GenerateJWT crea un token JWT para el usuario.
func GenerateJWT(id int64, role string, permissions []string, firstLogin bool) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		ID:          id,
		Role:        role,
		Permissions: permissions,
		FirstLogin:  firstLogin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}
