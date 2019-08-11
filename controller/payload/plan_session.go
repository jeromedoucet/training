package payload

import (
	"time"

	"github.com/google/uuid"
	"github.com/jeromedoucet/training/model"
)

type PlanSession struct {
	PlanId      string
	Day         time.Time `json:"day"`
	Description string    `json:"description"`
}

func (p PlanSession) ToModel() (*model.PlanSession, error) {
	var planId uuid.UUID
	var err error

	if err != nil {
		return nil, err
	}

	planId, err = uuid.Parse(p.PlanId)
	return &model.PlanSession{PlanId: planId, Day: p.Day, Description: p.Description}, nil
}
