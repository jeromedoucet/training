package plan_test

import (
	"github.com/google/uuid"
	"github.com/jeromedoucet/training/model"
	"github.com/jeromedoucet/training/test"
)

var userId uuid.UUID

func insertDataSet() {
	userId = uuid.New()
	test.InsertUser(&model.User{Id: userId, Login: "jerdct", Password: "titi_123456_tata"}, db)
}
