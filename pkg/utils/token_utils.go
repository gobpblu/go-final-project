package utils

import (
	"go-final-project/pkg/constants"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(password string) (string, error) {
	secret := []byte(password)

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	return jwtToken.SignedString(secret)
}

func IsTokenValid(tokenString string) bool {
	jwtToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// секретный ключ для всех токенов одинаковый, поэтому просто возвращаем его
		password := os.Getenv(constants.PasswordKey)
		return []byte(password), nil
	})

	if err != nil || !jwtToken.Valid {
		return false
	}

	return true
}
