package dao

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/model"
)

// PlanDAO allow accessing plan table
type PlanDAO struct {
	db         *sql.DB
	insertPlan *sql.Stmt
}

func createPlanDAO(db *sql.DB) *PlanDAO {
	var err error
	p := &PlanDAO{db: db}

	if p.insertPlan, err = p.db.Prepare(`
		INSERT INTO "plan" ("id", "creator_id", "trainee_id", "name") VALUES ($1, $2, $3, $4)
	`); err != nil {
		log.Fatalf("An error is returned during inserPlan statement initialization %s", err.Error())
	}

	return p
}

// Insert a new plan
func (p *PlanDAO) Insert(ctx context.Context, plan *model.Plan) (*model.Plan, *DbError) {
	plan.Id = uuid.New()
	_, err := p.insertPlan.ExecContext(ctx, plan.Id.String(), plan.CreatorId.String(), plan.TraineeId.String(), plan.Name)
	if err != nil {
		return nil, &DbError{Message: err.Error(), Type: UNKNOWN}
	}

	return plan, nil
}
