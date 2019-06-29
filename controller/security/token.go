package security

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// SetAuthCookie will create a jwt token and write it in a http only cookie.
func SetAuthCookie(w http.ResponseWriter, secret string, exp time.Time) error {
	t, err := createToken(secret, exp)
	if err != nil {
		return err
	}
	c := &http.Cookie{
		Name:     "auth",
		Value:    t,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, c)
	return nil
}

func createToken(secret string, exp time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": exp.Unix(),
	})
	res, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return res, nil
}
