package test

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func TokenValid(tokenString string, jwtSecret string) (res bool, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}
