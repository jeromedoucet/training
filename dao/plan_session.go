package dao

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/model"
)

// PlanDAO allow accessing plan table
type PlanSessionDAO struct {
	db                *sql.DB
	insertPlanSession *sql.Stmt
}

func createPlanSessionDAO(db *sql.DB) *PlanSessionDAO {
	var err error
	p := &PlanSessionDAO{db: db}

	if p.insertPlanSession, err = p.db.Prepare(`
		INSERT INTO "plan_session" ("id", "plan_id", "from", "to", "description", "comments") VALUES ($1, $2, $3, $4, $5, $6)
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
