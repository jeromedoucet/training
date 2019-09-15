package plan

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jeromedoucet/route"
	"github.com/jeromedoucet/training/configuration"
	"github.com/jeromedoucet/training/controller/payload"
	"github.com/jeromedoucet/training/controller/response"
	"github.com/jeromedoucet/training/dao"
	"github.com/jeromedoucet/training/model"
)

func CreatePlanHandlerFunc(c *configuration.GlobalConf, conn *dao.Conn) func(context.Context, http.ResponseWriter, *http.Request) {
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
			response.RenderError(http.StatusBadRequest, err.Error(), w)
			return
		}

		plan, dbErr = conn.PlanDAO.Insert(ctx, plan)

		if dbErr != nil {
			response.RenderError(http.StatusInternalServerError, dbErr.Message, w)
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

func CreatePlanSessionHandlerFunc(c *configuration.GlobalConf, conn *dao.Conn) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createPlanSession(conn, ctx, w, r)
		} else if r.Method == http.MethodGet {
			listPlanSession(conn, ctx, w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func createPlanSession(conn *dao.Conn, ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	var dbErr *dao.DbError
	var payloadPlanSession *payload.PlanSession
	var planSession *model.PlanSession
	var body []byte

	d := json.NewDecoder(r.Body)
	err = d.Decode(&payloadPlanSession)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path := route.SplitPath(r.URL.Path)
	planId := path[len(path)-2]
	payloadPlanSession.PlanId = planId

	planSession, err = payloadPlanSession.ToModel()

	if err != nil {
		response.RenderError(http.StatusBadRequest, err.Error(), w)
		return
	}

	planSession, dbErr = conn.PlanSessionDAO.Insert(ctx, planSession)

	if dbErr != nil {
		var status int
		if dbErr.Type == dao.NOT_FOUND {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}

		response.RenderError(status, dbErr.Message, w)
		return
	}

	body, err = json.Marshal(response.FromPlanSessionModel(planSession))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func listPlanSession(conn *dao.Conn, ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	var planId string
	var rawFrom []byte
	var rawTo []byte
	from := time.Time{}
	to := time.Time{}
	var dbErr *dao.DbError
	var sessions []*model.PlanSession
	var res []*response.PlanSession
	var body []byte

	path := route.SplitPath(r.URL.Path)
	planId = path[len(path)-2]

	q := r.URL.Query()

	rawFrom, err = base64.StdEncoding.DecodeString(q.Get("from"))

	if err != nil {
		response.RenderError(http.StatusBadRequest, err.Error(), w)
		return
	}

	rawTo, err = base64.StdEncoding.DecodeString(q.Get("to"))

	if err != nil {
		response.RenderError(http.StatusBadRequest, err.Error(), w)
		return
	}

	err = from.UnmarshalText(rawFrom)

	if err != nil {
		response.RenderError(http.StatusBadRequest, err.Error(), w)
		return
	}

	err = to.UnmarshalText(rawTo)

	if err != nil {
		response.RenderError(http.StatusBadRequest, err.Error(), w)
		return
	}

	sessions, dbErr = conn.PlanSessionDAO.List(ctx, planId, from, to)

	if dbErr != nil {
		response.RenderError(http.StatusInternalServerError, dbErr.Message, w)
		return
	}

	res = make([]*response.PlanSession, len(sessions), len(sessions))

	for i, session := range sessions {
		res[i] = response.FromPlanSessionModel(session)
	}

	body, err = json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

}
