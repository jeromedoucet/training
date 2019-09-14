package test

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jeromedoucet/training/configuration"
	"github.com/jeromedoucet/training/model"
)

func Conf() *configuration.GlobalConf {
	return &configuration.GlobalConf{DbName: "training_test", User: "postgres", JwtExpiration: 120 * time.Minute}
}

// Clean empty all table in the training model
func CleanDB(db *sql.DB) {
	var err error
	_, err = db.Exec(`TRUNCATE "plan_session", "plan", "user"`)
	if err != nil {
		log.Fatalf("An error is returned during db clean up %s", err.Error())
	}
}

func InsertUser(user *model.User, db *sql.DB) {
	var res sql.Result
	var err error
	res, err = db.Exec(`
				INSERT INTO "user" ("id", "login", "first_name", "last_name", "email", "password") 
				VALUES ($1, $2, $3, $4, $5, crypt($6, gen_salt('bf')))`, user.Id.String(), user.Login, user.FirstName, user.LastName, user.LastName, user.Password)

	if err != nil {
		log.Fatalf("An error is returned during user insertion %s", err.Error())
	}

	affectedRows, _ := res.RowsAffected()

	fmt.Println(fmt.Sprintf("%d user(s) injected", affectedRows))
}

func InsertPlan(plan *model.Plan, db *sql.DB) {
	var res sql.Result
	var err error
	res, err = db.Exec(`
				INSERT INTO "training"."plan" ("id", "creator_id", "trainee_id", "name") VALUES ($1, $2, $3, $4)`, plan.Id.String(), plan.CreatorId.String(), plan.TraineeId.String(), plan.Name)

	if err != nil {
		log.Fatalf("An error is returned during plan insertion %s", err.Error())
	}

	affectedRows, _ := res.RowsAffected()

	fmt.Println(fmt.Sprintf("%d plan(s) injected", affectedRows))
}

func InsertPlanSession(planSession *model.PlanSession, db *sql.DB) {
	var res sql.Result
	var err error
	res, err = db.Exec(`
				INSERT INTO "training"."plan_session" ("id", "plan_id", "from", "to", "description", "comments") VALUES ($1, $2, $3, $4, $5, $6)`, planSession.Id.String(), planSession.PlanId.String(), planSession.From, planSession.To, planSession.Description, planSession.Comments)

	if err != nil {
		log.Fatalf("An error is returned during planSession insertion %s", err.Error())
	}

	affectedRows, _ := res.RowsAffected()

	fmt.Println(fmt.Sprintf("%d planSession(s) injected", affectedRows))
}
