package payload

import (
	"github.com/google/uuid"
	"github.com/jeromedoucet/training/model"
)

type Plan struct {
	Name      string `json:"name"`
	CreatorId string `json:"creator_id"`
	TraineeId string `json:"trainee_id"`
}

// ToModel create aplan model
func (plan *Plan) ToModel() (*model.Plan, error) {
	creatorId, _ := uuid.Parse(plan.CreatorId)
	traineeId, _ := uuid.Parse(plan.TraineeId)
	return &model.Plan{CreatorId: creatorId, TraineeId: traineeId, Name: plan.Name}, nil
}
