package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jeromedoucet/training/configuration"
	"github.com/jeromedoucet/training/controller/payload"
	"github.com/jeromedoucet/training/controller/response"
	"github.com/jeromedoucet/training/controller/security"
	"github.com/jeromedoucet/training/dao"
	"github.com/jeromedoucet/training/model"
)

func checkAuthHandlerFunc(c *configuration.GlobalConf, conn *dao.Conn) func(http.ResponseWriter, *http.Request) bool {
	return func(w http.ResponseWriter, r *http.Request) bool {
		cookie := security.GetAuthCookie(r)
		if cookie != nil {
			token, err := security.ExtractToken(cookie.Value, c.JwtSecret)
			if err != nil || !token.Valid {
				response.RenderError(http.StatusUnauthorized, "invalid token", w)
				return false
			} else {
				return true
			}
		} else {
			response.RenderError(http.StatusUnauthorized, "no token", w)
			return false
		}
	}
}

func authenticationHandlerFunc(c *configuration.GlobalConf, conn *dao.Conn) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var err error
		var payloadUser *payload.User
		var user *model.User
		var body []byte

		var dbErr *dao.DbError

		d := json.NewDecoder(r.Body)
		err = d.Decode(&payloadUser)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !payloadUser.AuthenticationPayloadValid() {
			response.RenderError(http.StatusBadRequest, "Missing some mandatory fields", w)
			return
		}

		user, dbErr = conn.UserDAO.CheckPassword(ctx, payloadUser.ToModel())

		if dbErr != nil {
			if dbErr.Type == dao.NOT_FOUND {
				response.RenderError(http.StatusUnauthorized, dbErr.Message, w)
			} else {
				response.RenderError(http.StatusInternalServerError, dbErr.Message, w)
			}
			return
		}

		body, err = json.Marshal(response.FromUserModel(user))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = security.SetAuthCookie(w, c.JwtSecret, time.Now().Add(c.JwtExpiration), user.Id)

		if err != nil {
			log.Println("Error when encoding the token: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}

func logoutHandlerFunc(c *configuration.GlobalConf, conn *dao.Conn) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		security.UnvalidatedAuthCookie(w, r)
	}
}

func getSessionHandlerFunc(c *configuration.GlobalConf, conn *dao.Conn) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var userId string
		var user *model.User
		var dbErr *dao.DbError
		var body []byte
		cookie := security.GetAuthCookie(r)
		if cookie != nil {
			token, err := security.ExtractToken(cookie.Value, c.JwtSecret)
			if err != nil || !token.Valid {
				security.UnvalidatedAuthCookie(w, r)
				response.RenderError(http.StatusUnauthorized, "invalid token", w)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userId = claims["sub"].(string)

			user, dbErr = conn.UserDAO.GetUser(ctx, userId)

			if dbErr != nil {
				log.Println("Error when fetching the user: ", dbErr.Message)
				response.RenderError(http.StatusInternalServerError, dbErr.Message, w)
				return
			}

			err = security.SetAuthCookie(w, c.JwtSecret, time.Now().Add(c.JwtExpiration), user.Id)

			if err != nil {
				log.Println("Error when encoding the token: ", err)
				response.RenderError(http.StatusInternalServerError, err.Error(), w)
				return
			}

			body, err = json.Marshal(response.FromUserModel(user))

			if err != nil {
				log.Println("Error when marshaling the response: ", err)
				response.RenderError(http.StatusInternalServerError, err.Error(), w)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(body)

		} else {
			response.RenderError(http.StatusUnauthorized, "not authenticated", w)
		}
	}
}
