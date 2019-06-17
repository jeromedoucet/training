package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/jeromedoucet/training/configuration"
	"github.com/jeromedoucet/training/controller/payload"
	"github.com/jeromedoucet/training/controller/security"
	"github.com/jeromedoucet/training/dao"
)

func authenticationHandlerFunc(c *configuration.GlobalConf, conn *dao.Conn) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var err error
		var payloadUser *payload.User
		var isAuthenticated bool

		var dbErr *dao.DbError

		d := json.NewDecoder(r.Body)
		err = d.Decode(&payloadUser)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !payloadUser.AuthenticationPayloadValid() {
			renderError(http.StatusBadRequest, "Missing some mandatory fields", w)
			return
		}

		isAuthenticated, dbErr = conn.UserDAO.CheckPassword(ctx, payloadUser.ToModel())

		if dbErr != nil {
			renderError(http.StatusInternalServerError, dbErr.Message, w)
			return
		}

		if !isAuthenticated {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = security.SetAuthCookie(w, c.JwtSecret, time.Now().Add(time.Minute*10))

		if err != nil {
			log.Println("Error when encoding the token: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		w.Header().Set("Content-Type", "application/json")

	}
}
