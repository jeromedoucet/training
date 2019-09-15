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
var firstSession, secondSession, thirdSession *model.PlanSession

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
	test.InsertPlanSession(firstSession, db)

	secondSession = &model.PlanSession{
		Id:          uuid.New(),
		PlanId:      plan.Id,
		From:        time.Date(2019, time.October, 1, 12, 0, 0, 0, time.UTC),
		To:          time.Date(2019, time.October, 1, 13, 0, 0, 0, time.UTC),
		Description: "VMA interval training (1' X 1')",
	}
	test.InsertPlanSession(secondSession, db)

	thirdSession = &model.PlanSession{
		Id:          uuid.New(),
		PlanId:      plan.Id,
		From:        time.Date(2019, time.October, 5, 12, 0, 0, 0, time.UTC),
		To:          time.Date(2019, time.October, 5, 13, 0, 0, 0, time.UTC),
		Description: "VMA interval training (2' X 2')",
	}
	test.InsertPlanSession(thirdSession, db)

}
