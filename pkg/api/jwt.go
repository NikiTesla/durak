package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	Token    *jwt.Token
	jwt.StandardClaims
}

func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func parseClaims(r *http.Request) (*Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, ErrAuthorisationHeaderMissing
	}

	tokenStrings := strings.Split(authHeader, "Bearer ")
	if len(tokenStrings) != 2 {
		return nil, ErrInvalidToken
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStrings[1], claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	claims.Token = token

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
