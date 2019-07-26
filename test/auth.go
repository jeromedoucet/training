package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jeromedoucet/training/configuration"
)

func CheckAuthCookie(cookies []*http.Cookie, conf *configuration.GlobalConf, t *testing.T) *jwt.Token {
	cookie := GetCookieByName(cookies, "auth")

	if cookie == nil {
		t.Fatal("Expect auth cookie, but no such cookie found")
	}

	if !cookie.HttpOnly {
		t.Fatal("Expect the token to be httpOnly, but it is not")
	}

	if cookie.Expires.IsZero() {
		t.Fatal("Expected Expire to be set")
	}

	token, err := ExtractToken(cookie.Value, conf.JwtSecret)

	if err != nil {
		t.Fatalf("Expect no error when validating token, but get %s", err)
	}

	if !token.Valid {
		t.Fatalf("Expect token %s to be valid", cookie.Value)
	}

	return token
}

func GetCookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func ExtractToken(tokenString string, jwtSecret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
	return token, err
}

func CreateToken(secret string, exp time.Time, id uuid.UUID, t *testing.T) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": exp.Unix(),
		"sub": id.String(),
	})
	res, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Got error %s when creating token for the test", err.Error())
		return ""
	}
	return res
}
