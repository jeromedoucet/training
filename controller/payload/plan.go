package payload

import (
	"errors"

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
	var err error
	var creatorId uuid.UUID
	var traineeId uuid.UUID

	creatorId, err = uuid.Parse(plan.CreatorId)

	if err != nil {
		return nil, err
	}

	traineeId, err = uuid.Parse(plan.TraineeId)

	if err != nil {
		return nil, err
	}

	if plan.Name == "" {
		return nil, errors.New("Expect non empty name")
	}
	return &model.Plan{CreatorId: creatorId, TraineeId: traineeId, Name: plan.Name}, nil
}
