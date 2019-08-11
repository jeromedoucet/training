package plan_test

import (
	"github.com/google/uuid"
	"github.com/jeromedoucet/training/model"
	"github.com/jeromedoucet/training/test"
)

var userId uuid.UUID
var user *model.User
var plan *model.Plan

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
}
