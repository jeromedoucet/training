package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/jeromedoucet/training/configuration"
)

func CheckAuthCookie(cookies []*http.Cookie, conf *configuration.GlobalConf, t *testing.T) {
	var err error

	cookie := GetCookieByName(cookies, "auth")

	if cookie == nil {
		t.Fatal("Expect auth cookie, but no such cookie found")
	}

	if !cookie.HttpOnly {
		t.Fatal("Expect the token to be httpOnly, but it is not")
	}

	var tokenValid bool
	tokenValid, err = TokenValid(cookie.Value, conf.JwtSecret)

	if err != nil {
		t.Fatalf("Expect no error when validating token, but get %s", err)
	}

	if !tokenValid {
		t.Fatalf("Expect token %s to be valid", cookie.Value)
	}
}

func GetCookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

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
