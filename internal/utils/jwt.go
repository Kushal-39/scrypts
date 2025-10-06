package utils

import (
	"net/http"
	"scrypts/internal/auth"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func GetUsernameFromJWT(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", http.ErrNoCookie // Use as generic error
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return auth.JwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}
	username, ok := claims["username"].(string)
	if !ok {
		return "", err
	}
	return username, nil
}
