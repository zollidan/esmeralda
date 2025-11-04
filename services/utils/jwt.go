package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zollidan/esmeralda/config"
)

// TODO: Refactor to accept parameters in loop instead of fixed ones
func GenerateJWT(cfg *config.Config, params ...string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": params[0],
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
