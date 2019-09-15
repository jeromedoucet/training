package dao

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/model"
)

// PlanDAO allow accessing plan table
type PlanSessionDAO struct {
	db                *sql.DB
	insertPlanSession *sql.Stmt
	listPlanSession   *sql.Stmt
}

func createPlanSessionDAO(db *sql.DB) *PlanSessionDAO {
	var err error
	p := &PlanSessionDAO{db: db}

	if p.insertPlanSession, err = p.db.Prepare(`
		INSERT INTO "plan_session" ("id", "plan_id", "from", "to", "description", "comments") VALUES ($1, $2, $3, $4, $5, $6)
	`); err != nil {
		log.Fatalf("An error is returned during inserPlanSession statement initialization %s", err.Error())
	}

	if p.listPlanSession, err = p.db.Prepare(`
			SELECT
				"id",
				"plan_id",
				"from",
				"to",
				"description",
				"comments"
			FROM "plan_session" as "s"
			WHERE "s"."plan_id" = $1 AND "s"."from" >= $2 AND "s"."to" <= $3
	`); err != nil {
		log.Fatalf("An error is returned during inserPlanSession statement initialization %s", err.Error())
	}

	return p
}

// Insert a new plan
func (p *PlanSessionDAO) Insert(ctx context.Context, planSession *model.PlanSession) (*model.PlanSession, *DbError) {
	planSession.Id = uuid.New()
	_, err := p.insertPlanSession.ExecContext(ctx, planSession.Id.String(), planSession.PlanId.String(), planSession.From, planSession.To, planSession.Description, planSession.Comments)
	if err != nil {
		if strings.Contains(err.Error(), "plan_session_plan_id_fkey") {
			return nil, &DbError{Message: err.Error(), Type: NOT_FOUND}
		} else {
			return nil, &DbError{Message: err.Error(), Type: UNKNOWN}
		}
	}

	return planSession, nil
}

func (p *PlanSessionDAO) List(ctx context.Context, planId string, from time.Time, to time.Time) ([]*model.PlanSession, *DbError) {
	rows, err := p.listPlanSession.QueryContext(ctx, planId, from, to)

	if err != nil {
		return nil, &DbError{Message: err.Error(), Type: UNKNOWN}
	}

	res := []*model.PlanSession{}

	for rows.Next() {
		var strId string
		var parsedId uuid.UUID
		var strPlanId string
		var parsedPlanId uuid.UUID
		session := &model.PlanSession{}

		err = rows.Scan(
			&strId,
			&strPlanId,
			&session.From,
			&session.To,
			&session.Description,
			&session.Comments,
		)

		if err != nil {
			return nil, &DbError{Message: err.Error(), Type: UNKNOWN}
		}

		parsedId, err = uuid.Parse(strId)

		if err != nil {
			return nil, &DbError{Message: err.Error(), Type: UNKNOWN}
		}

		parsedPlanId, err = uuid.Parse(strPlanId)

		if err != nil {
			return nil, &DbError{Message: err.Error(), Type: UNKNOWN}
		}

		session.Id = parsedId
		session.PlanId = parsedPlanId

		res = append(res, session)
	}

	return res, nil

}
