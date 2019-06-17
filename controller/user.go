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
	"github.com/jeromedoucet/training/controller/response"
	"github.com/jeromedoucet/training/controller/security"
	"github.com/jeromedoucet/training/dao"
	"github.com/jeromedoucet/training/model"
)

func createUserHandlerFunc(c *configuration.GlobalConf, conn *dao.Conn) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

		var err error
		var dbErr *dao.DbError
		var payloadUser *payload.User
		var user *model.User
		var body []byte

		d := json.NewDecoder(r.Body)
		err = d.Decode(&payloadUser)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !payloadUser.CreationPayloadValid() {
			renderError(http.StatusBadRequest, "Missing some mandatory fields", w)
			return
		}

		user, dbErr = conn.UserDAO.Insert(ctx, payloadUser.ToModel())

		if dbErr != nil {
			var status int
			if dbErr.Type == dao.CONFLICT {
				status = http.StatusConflict
			} else {
				status = http.StatusInternalServerError
			}
			renderError(status, dbErr.Message, w)
			return
		}

		body, err = json.Marshal(response.FromUserModel(user))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
		w.Write(body)
	}
}
