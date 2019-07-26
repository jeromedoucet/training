package security

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func GetAuthCookie(req *http.Request) *http.Cookie {
	c, err := req.Cookie("auth")
	if err != nil {
		return nil
	}
	return c
}

// SetAuthCookie will create a jwt token and write it in a http only cookie.
func SetAuthCookie(w http.ResponseWriter, secret string, exp time.Time, sub uuid.UUID) error {
	t, err := createToken(secret, exp, sub)
	if err != nil {
		return err
	}
	c := &http.Cookie{
		Name:     "auth",
		Value:    t,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  exp, // keep it a long time. The "real" expiration time is handle by jwt
	}

	http.SetCookie(w, c)
	return nil
}

func UnvalidatedAuthCookie(w http.ResponseWriter, r *http.Request) {
	cookie := GetAuthCookie(r)
	if cookie != nil {
		cookie.Expires = time.Now().Add(-1 * time.Minute)
		cookie.MaxAge = -1
		cookie.Path = "/"
		http.SetCookie(w, cookie)
	}
}

func createToken(secret string, exp time.Time, sub uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": exp.Unix(),
		"sub": sub.String(),
	})
	res, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return res, nil
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
