package dao

import (
	"context"
	"database/sql"
	"log"

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
		INSERT INTO "plan_session" ("id", "plan_id", "day", "description", "comments") VALUES ($1, $2, $3, $4, $5)
	`); err != nil {
		log.Fatalf("An error is returned during inserPlanSession statement initialization %s", err.Error())
	}

	return p
}

// Insert a new plan
func (p *PlanSessionDAO) Insert(ctx context.Context, planSession *model.PlanSession) (*model.PlanSession, *DbError) {
	planSession.Id = uuid.New()
	_, err := p.insertPlanSession.ExecContext(ctx, planSession.Id.String(), planSession.PlanId.String(), planSession.Day, planSession.Description, planSession.Comments)
	if err != nil {
		return nil, &DbError{Message: err.Error(), Type: UNKNOWN}
	}

	return planSession, nil
}
