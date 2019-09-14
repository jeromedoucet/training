package plan_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/model"
	"github.com/jeromedoucet/training/test"
)

var userId uuid.UUID
var user *model.User
var plan *model.Plan
var firstSession *model.PlanSession

func insertDataSet() {
	// user
	userId = uuid.New()
	user = &model.User{Id: userId, Login: "jerdct", Password: "titi_123456_tata"}
	test.InsertUser(user, db)

	// plan
	plan = &model.Plan{
		Id:        uuid.New(),
		Name:      "some training plan for test",
		CreatorId: userId,
		TraineeId: userId,
	}
	test.InsertPlan(plan, db)

	// plan session
	firstSession = &model.PlanSession{
		Id:          uuid.New(),
		PlanId:      plan.Id,
		From:        time.Date(2019, time.September, 30, 12, 0, 0, 0, time.UTC),
		To:          time.Date(2019, time.September, 30, 13, 0, 0, 0, time.UTC),
		Description: "Long run in the mountains.",
		Comments:    "Nice training !",
	}

}
