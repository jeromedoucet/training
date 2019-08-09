package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jeromedoucet/training/configuration"
	"github.com/jeromedoucet/training/controller/payload"
	"github.com/jeromedoucet/training/controller/response"
	"github.com/jeromedoucet/training/dao"
	"github.com/jeromedoucet/training/model"
)

func createPlanHandlerFunc(c *configuration.GlobalConf, conn *dao.Conn) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var err error
		var dbErr *dao.DbError
		var payloadPlan *payload.Plan
		var plan *model.Plan
		var body []byte

		d := json.NewDecoder(r.Body)
		err = d.Decode(&payloadPlan)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		plan, err = payloadPlan.ToModel()

		if err != nil {
			renderError(http.StatusBadRequest, err.Error(), w)
			return
		}

		plan, dbErr = conn.PlanDAO.Insert(ctx, plan)

		if dbErr != nil {
			renderError(http.StatusInternalServerError, dbErr.Message, w)
			return
		}

		body, err = json.Marshal(response.FromPlanModel(plan))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}
