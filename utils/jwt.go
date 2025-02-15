package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key untuk JWT
var jwtSecret = []byte("your-secret-key")

// GenerateJWT untuk membuat token JWT dengan UserID dan Role
func GenerateJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
