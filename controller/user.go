package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jeromedoucet/training/configuration"
	"github.com/jeromedoucet/training/controller/payload"
	"github.com/jeromedoucet/training/controller/response"
	"github.com/jeromedoucet/training/dao"
	"github.com/jeromedoucet/training/model"
)

func createUserHandlerFunc(c *configuration.GlobalConf, conn *dao.Conn) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

		var err error
		var dbErr *dao.DbError
		var payloadUser *payload.User
		var user *model.User
		var token string
		var body []byte

		d := json.NewDecoder(r.Body)
		err = d.Decode(&payloadUser)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // todo use a specific struct for error
			return
		}

		user, dbErr = conn.UserDAO.Insert(ctx, payloadUser.ToModel())

		if dbErr != nil {
			res := response.Error{Message: dbErr.Message}
			body, _ = json.Marshal(res)
			if dbErr.Type == dao.CONFLICT {
				w.WriteHeader(http.StatusConflict)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write(body)
			return
		}

		body, err = json.Marshal(response.FromUserModel(user))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // todo handle many case
			return
		}

		token, err = createToken(c.JwtSecret, time.Now().Add(time.Minute*10)) // todo make expiration time configurable

		if err != nil {
			log.Println("Error when encoding the token: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("authorization", fmt.Sprintf("Bearer %s", token))
		w.Write(body)
	}
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
